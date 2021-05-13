package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/yt-spam/internal/checks/blchecks"
	"github.com/darmiel/yt-spam/internal/checks/copycat"
	"github.com/darmiel/yt-spam/internal/checks/fmtspam"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	VideoNotFound = errors.New("video not found")
)

func (cmd *Command) CheckCommand(c *cli.Context) error {
	videoID := c.String("video-id")
	force := c.Bool("cache-yes")
	return cmd.c(videoID, force)
}

func (cmd *Command) c(videoID string, forceUseCache bool) error {
	cached, err := cmd.a(videoID, forceUseCache)
	if err != nil {
		return err
	}
	checker := ytspam.NewCommentChecker(cached.WrapArray())

	checker.Check(blchecks.NewChannelBlacklistCheck(checker.ResChannels),
		blchecks.NewCommentBlacklistCheck(checker.ResComments),
		blchecks.NewNameBlacklistCheck(checker.ResChannels),
		copycat.NewCommentCopyCatCheck(checker.ResComments, 28),
		fmtspam.NewFormatSpamCheck(checker.ResComments, 7, 4))

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go checker.ReadChannels(os.Stdout, wg)
	go checker.ReadComments(os.Stdout, wg)
	go checker.ReadErrors(os.Stdout, wg)
	wg.Wait()

	return nil
}

// TODO: give me a descriptive name, please :(
func (cmd *Command) a(videoID string, forceUseCache bool) (*ytspam.CachedComments, error) {
	data := make(map[string]*ytspam.CachedComments)
	if videoID == "" {
		var p survey.Prompt

		// ask if from cache
		p = &survey.Confirm{Message: "No video ID specified. Do you want to select a video from the cache?"}
		var fc bool
		if err := survey.AskOne(p, &fc); err != nil {
			return nil, err
		}
		if !fc {
			return nil, VideoNotFound
		}

		// read all cached files
		glob, err := filepath.Glob(path.Join("data", "cache", "*.json"))
		if err != nil {
			return nil, err
		}

		log.Println(":: Reading", len(glob), "files...")

		var keys []string
		for _, pth := range glob {
			// read cache #n
			c, err := ytspam.CacheFromPath(pth)
			if err != nil {
				log.Println("Failed at file:", pth)
				return nil, err
			}
			display := fmt.Sprintf("[%s]: %s", c.LastUpdate.Format("02.01.2006 15:04:05"), c.VideoTitle)

			data[c.VideoTitle] = c
			data[c.VideoID] = c
			data[display] = c

			keys = append(keys, display)
		}

		p = &survey.Select{Message: "Select Cached Video", Options: keys}
		var title string
		if err := survey.AskOne(p, &title, survey.WithPageSize(25)); err != nil {
			return nil, err
		}
		if title == "" {
			return nil, VideoNotFound
		}

		// get id
		c, found := data[title]
		if !found {
			return nil, VideoNotFound
		}

		videoID = c.VideoID
		log.Println("Using Video ID:", videoID)
	}

	fromCache := true

	pth := path.Join("data", "cache", videoID+".json")
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		fromCache = false
	} else {
		if !forceUseCache {
			q := &survey.Confirm{Message: "Load comments from cache?"}
			if err := survey.AskOne(q, &fromCache); err != nil {
				return nil, err
			}
		}
	}

	// load from cache
	if fromCache {
		if c, ok := data[videoID]; ok {
			return c, nil
		}

		// read from file
		return ytspam.CacheFromVideoID(videoID)
	}

	// load from youtube
	service := cmd.Service
	call := service.Videos.List([]string{"snippet", "statistics"}).Id(videoID)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) <= 0 {
		return nil, VideoNotFound
	}
	video := resp.Items[0]

	p := termenv.ColorProfile()
	fmt.Println(termenv.String("YT-SPAM").Foreground(p.Color("0")).Background(p.Color("#E88388")),
		"Reading comments from",
		termenv.String(video.Snippet.Title).Foreground(p.Color("#A8CC8C")),
		"(",
		termenv.String(video.Snippet.PublishedAt).Foreground(p.Color("#D290E4")),
		")")

	reader := ytspam.NewCommentReader(service, video, &ytspam.CommentReaderConfig{
		DisplayBar: true,
		Silent:     false,
	})

	if err := reader.Read(); err != nil {
		return nil, err
	}

	var comments []*youtube.Comment
	for _, v := range reader.GetComments() {
		comments = append(comments, v)
	}

	cmt := &ytspam.CachedComments{
		VideoID:    video.Id,
		VideoTitle: video.Snippet.Title,
		LastUpdate: time.Now(),
		Comments:   comments,
	}

	// save to file
	if err := cmt.Save(); err != nil {
		return nil, err
	}

	return cmt, nil
}

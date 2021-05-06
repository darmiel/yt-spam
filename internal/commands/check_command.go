package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	blacklist_checks "github.com/darmiel/yt-spam/internal/checks/blacklist-checks"
	"github.com/darmiel/yt-spam/internal/checks/copycat"
	fmt_spam "github.com/darmiel/yt-spam/internal/checks/fmt-spam"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"path"
	"path/filepath"
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
	checker := ytspam.NewCommentChecker(cached.Wrap())
	if err := checker.Check(
		&blacklist_checks.NameBlacklistCheck{},
		&blacklist_checks.CommentBlacklistCheck{},
		&fmt_spam.FormatSpamCheck{},
		&copycat.CommentCopyCatCheck{}); err != nil {
		return err
	}
	/*
		fmt.Println()
		log.Println("Found:")
			for id, violations := range checker.Violations() {
				log.Println("*", "https://www.youtube.com/channel/"+id, ":")
				ratings := make(map[string]checks.Rating)
				for _, vl := range violations {
					r, ok := ratings[vl.Check.Name()]
					if !ok {
						r = 0
					}
					r += vl.Rating
					ratings[vl.Check.Name()] = r
				}
				for cn, cr := range ratings {
					log.Println("  ├", cn, "::", cr, "(", cr.IsViolation(), ")")
				}
			}
	*/
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
		glob, err := filepath.Glob(path.Join("data", "cache", "**"))
		if err != nil {
			return nil, err
		}

		var keys []string
		for _, pth := range glob {
			// read cache #n
			c, err := ytspam.CacheFromPath(pth)
			if err != nil {
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
		if err := survey.AskOne(p, &title); err != nil {
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

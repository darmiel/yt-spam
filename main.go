package main

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks"
	nameblacklist "github.com/darmiel/yt-spam/internal/checks/name-blacklist"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/muesli/termenv"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const videoId string = "ddBBWrBa6QU"

func main() {
	ctx := context.Background()

	// read api key from "api-key.txt"
	b, err := ioutil.ReadFile("api-key.txt")
	if err != nil {
		log.Fatalln("Error reading api-key.txt:", err)
		return
	}

	apiKey := string(b)

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}

	call := service.Videos.List([]string{"snippet", "statistics"}).Id(videoId)
	resp, err := call.Do()
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}
	if len(resp.Items) <= 0 {
		log.Fatalln("No video found")
		return
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

	// TODO: CACHE ONLY FOR TESTING PURPOSES.
	cp := path.Join("data", "cache", videoId+".json")
	if _, err := os.Stat(cp); os.IsNotExist(err) {
		log.Println("READ :: From YouTube")
		// not cache
		if err := reader.Read(); err != nil {
			log.Fatalln("WARN ::", err)
			return
		}
		b, err := reader.ToJSON()
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Println("READ :: Saving To Cache")
		if err = ioutil.WriteFile(cp, b, 0777); err != nil {
			log.Fatalln(err)
			return
		}
	} else {
		log.Println("READ :: From Cache")
		// from cache
		b, err := ioutil.ReadFile(cp)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if err = reader.FromJSON(b); err != nil {
			log.Fatalln(err)
			return
		}
	}

	checker := ytspam.NewCommentChecker(reader.GetComments())
	if err := checker.Check(
		// &copycat.CommentCopyCatCheck{},
		&nameblacklist.NameBlacklistCheck{}); err != nil {
		log.Fatalln("WARN ::", err)
		return
	}

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
}

package main

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks/copycat"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/muesli/termenv"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

const videoId string = "SacqnEO770E"

func main() {
	ctx := context.Background()

	// read api key from "api-key.txt"
	b, err := ioutil.ReadFile("api-key.txt")
	if err != nil {
		log.Fatalln("Error reading api-key.txt:", err)
		return
	}

	apiKey := string(b)
	log.Println("Read API-Key:", apiKey)

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}

	call := service.Videos.List([]string{"snippet"}).
		Id(videoId)
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

	reader := ytspam.NewCommentReader(service, videoId)
	if err := reader.Read(); err != nil {
		log.Println("WARN ::", err)
	}

	checker := ytspam.NewCommentChecker(reader.GetComments())
	if err := checker.Check(&copycat.CommentCopyCatCheck{}); err != nil {
		log.Println("WARN ::", err)
	}
}

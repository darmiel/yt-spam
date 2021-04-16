package main

import (
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/darmiel/yt-spam/internal/ytspam/checks"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

const videoId string = "lR1LdQS9KCo"

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

	reader := ytspam.NewCommentReader(service, videoId)
	if err := reader.Read(); err != nil {
		log.Println("WARN ::", err)
	}

	checker := ytspam.NewCommentChecker(reader.GetComments())
	if err := checker.Check(&checks.CommentCopyCatCheck{}); err != nil {
		log.Println("WARN ::", err)
	}
}

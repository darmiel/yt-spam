// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List([]string{part})
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func readVideoComments(service *youtube.Service, videoID string) {
	call := service.CommentThreads.List([]string{"id", "replies", "snippet"}).
		// Filters
		VideoId(videoID).
		// Parameters
		MaxResults(100).
		Order("relevance")

	resp, err := call.Do()
	if err != nil {
		log.Fatalln("Error fetching comments:", err)
		return
	}
	log.Println("Comments:")
	for i, t := range resp.Items {
		b, err := json.Marshal(t)
		if err != nil {
			fmt.Println("### Skipping", i, "::", err)
			continue
		}
		// header
		fmt.Println()
		fmt.Println("((", i, ")) ::", t.Id, "[", t.Snippet.TopLevelComment.Snippet.PublishedAt, ")")
		fmt.Println(t.Snippet.TopLevelComment.Snippet.TextOriginal)
		fmt.Println("---")
		fmt.Println(string(b))
		fmt.Println("---")
		fmt.Println()
	}
}

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
	// channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
	readVideoComments(service, "O85OWBJ2ayo")
}

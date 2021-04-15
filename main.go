// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/darmiel/yt-spam/internal/ytspam/checks"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var c = []ytspam.Check{
	checks.CommentDownloadProfilePictureCheck{},
}
var a []*WrappedComment

type WrappedComment struct {
	AuthorName string
	Body       string
	Date       string
	Replies    []*WrappedReply
}
type WrappedReply struct {
	AuthorName string
	Body       string
	Date       string
}

func readVideoComments(service *youtube.Service, videoID, npt string) error {
	call := service.CommentThreads.List([]string{"id", "replies", "snippet"}).
		// Filters
		VideoId(videoID).
		// Parameters
		MaxResults(100).
		Order("relevance")

	if npt != "" {
		log.Println("Using Next Page Token")
		call = call.PageToken(npt)
	}

	resp, err := call.Do()
	if err != nil {
		return err
	}

	// fetch comments
	for i, t := range resp.Items {
		com := &WrappedComment{
			AuthorName: t.Snippet.TopLevelComment.Snippet.AuthorDisplayName,
			Body:       t.Snippet.TopLevelComment.Snippet.TextOriginal,
			Date:       t.Snippet.TopLevelComment.Snippet.PublishedAt,
		}

		fmt.Println()
		fmt.Println("((", i, ")) ::", t.Id, "[", t.Snippet.TopLevelComment.Snippet.PublishedAt, ")")
		fmt.Println()
		fmt.Println(t.Snippet.TopLevelComment.Snippet.TextOriginal)
		fmt.Println()

		// replies
		if t.Replies != nil {
			l := len(t.Replies.Comments)
			if l != 0 {
				for _, r := range t.Replies.Comments {
					com.Replies = append(com.Replies, &WrappedReply{
						AuthorName: r.Snippet.AuthorDisplayName,
						Body:       r.Snippet.TextOriginal,
						Date:       r.Snippet.PublishedAt,
					})
					fmt.Println(" ├", r.Snippet.AuthorDisplayName, ":", r.Snippet.TextOriginal)
				}
				if int64(l) != t.Snippet.TotalReplyCount {
					fmt.Println(" └ and", t.Snippet.TotalReplyCount-int64(l), "more...")
				}
			}
		}

		a = append(a, com)

		fmt.Println("(( Checking comment with comment checks ... ))")
		for i, ch := range c {
			cc, ok := ch.(ytspam.CommentCheck)
			if !ok {
				continue
			}
			fmt.Print(" ├ Check ", i, ": ")

			rating, err := cc.CheckComment(t.Snippet.TopLevelComment)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Rating (", rating, ")")
		}

		fmt.Println()
		fmt.Println()
		fmt.Println("---")
		fmt.Println()
	}

	if resp.NextPageToken != "" {
		return readVideoComments(service, videoID, resp.NextPageToken)
	}

	return nil
}

const videoId string = "O85OWBJ2ayo"

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

	p := path.Join("data", "comments")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0666); err != nil {
			log.Fatalln(err)
			return
		}
	}
	p = path.Join(p, videoId+".json")

	log.Println("error? :=", readVideoComments(service, videoId, ""))

	b, err = json.Marshal(a)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("error (write) :=", ioutil.WriteFile(p, b, 0644))
}

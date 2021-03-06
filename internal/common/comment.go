package common

import (
	"google.golang.org/api/youtube/v3"
	"strings"
	"time"
)

func CommentBefore(a, b *youtube.Comment) bool {
	t1, err := time.Parse(time.RFC3339, a.Snippet.PublishedAt)
	if err != nil {
		return false
	}
	t2, err := time.Parse(time.RFC3339, b.Snippet.PublishedAt)
	if err != nil {
		return false
	}
	return t1.Before(t2)
}

func TrimBody(c *youtube.Comment) (res string) {
	res = c.Snippet.TextOriginal

	// lower case
	res = strings.ToLower(res)

	// trim new line
	res = strings.TrimRight(res, "\n .,:-")
	res = strings.TrimLeft(res, "\n .,:-")

	return
}

func GetMatchingComments(orig *youtube.Comment, all []*youtube.Comment) (matches []*youtube.Comment) {
	bodyA := TrimBody(orig)
	for _, a := range all {
		bodyB := TrimBody(a)
		if bodyA == bodyB {
			matches = append(matches, a)
		}
		// TODO: Text Matching
	}
	return
}

func CommentFromHimHimself(a, b *youtube.Comment) bool {
	if a.Snippet == nil || b.Snippet == nil {
		return false
	}
	if a.Snippet.AuthorChannelId == nil || b.Snippet.AuthorChannelId == nil {
		return false
	}
	return a.Snippet.AuthorChannelId.Value == b.Snippet.AuthorChannelId.Value
}

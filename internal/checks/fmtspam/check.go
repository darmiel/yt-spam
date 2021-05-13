package fmtspam

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/common"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"strings"
)

type formatSpamCheck struct {
	channel chan *checks.CommentRatingNotify
	words   map[string]uint64
	minLen  int
	minOcc  uint64
}

func NewFormatSpamCheck(channel chan *checks.CommentRatingNotify, minLen int, minOcc uint64) *formatSpamCheck {
	return &formatSpamCheck{channel, make(map[string]uint64), minLen, minOcc}
}

func (c *formatSpamCheck) Name() string {
	return "Format-Spam"
}

func (c *formatSpamCheck) Prefix() termenv.Style {
	return common.CreatePrefix("🔁", "FMT-SPAM", "DBAB79")
}

// 0 -> comment
// 1 -> word
// 2 -> occurrences
// 3 -> rating
func (c *formatSpamCheck) SendViolation(i ...interface{}) {
	var (
		comment = i[0].(*youtube.Comment)
		word    = i[1].(string)
		occ     = i[2].(uint64)
		rating  = i[3].(checks.Rating)
	)
	c.channel <- &checks.CommentRatingNotify{
		Reason:  fmt.Sprintf("comment contained '%s' (%vx)", word, occ),
		Comment: comment,
		Rating:  rating,
		Check:   c,
	}
}

///

func (c *formatSpamCheck) CheckComments(comments []*youtube.Comment) {
	for _, comment := range comments {
		body := comment.Snippet.TextOriginal
		if compare.ContainsHomoglyphs(body) {
			body = compare.Normalize(body)
		}

		// extract formatted
		for _, w := range compare.ExtractFormattedWords(body) {
			w = strings.ToLower(w)

			// ignore short words
			if len(w) < c.minLen {
				continue
			}

			num := c.words[w] + 1
			c.words[w] = num
		}
	}

	// remove all extracted formatted-texts where occurrences < min
	for idx, num := range c.words {
		if num < c.minOcc {
			delete(c.words, idx)
		}
	}

	for _, comment := range comments {
		body := comment.Snippet.TextOriginal
		if compare.ContainsHomoglyphs(body) {
			body = compare.Normalize(body)
		}

		// extract formatted
		for _, w := range compare.ExtractFormattedWords(body) {
			w = strings.ToLower(w)

			// ignore short words
			if len(w) < c.minLen {
				continue
			}

			occ, ok := c.words[w]
			if !ok {
				continue
			}

			c.SendViolation(comment, w, occ, checks.Rating(100))
		}
	}
}

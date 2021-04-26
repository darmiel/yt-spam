package fmt_spam

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"google.golang.org/api/youtube/v3"
	"strings"
)

const (
	FmtSpamMinLen         = 5
	FmtSpamMinOccurrences = 3
)

type CommentBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
	words      map[string]uint64
}

func (c *CommentBlacklistCheck) Name() string {
	return "Format-Spam"
}

func (c *CommentBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	c.words = make(map[string]uint64)
	return nil
}

func (c *CommentBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *CommentBlacklistCheck) CheckComments(all map[string]*youtube.Comment) error {
	bar := pb.New(len(all) * 2)
	for _, comment := range all {
		bar.Increment()

		body := comment.Snippet.TextOriginal
		if compare.ContainsHomoglyphs(body) {
			body = compare.Normalize(body)
		}

		// extract formatted
		for _, w := range extractFormattedWords(body) {
			w = strings.ToLower(w)

			// ignore short words
			if len(w) < FmtSpamMinLen {
				continue
			}

			num := c.words[w] + 1
			c.words[w] = num
		}
	}

	// remove all extracted formatted-texts where occurrences < min
	for idx, num := range c.words {
		if num < FmtSpamMinOccurrences {
			delete(c.words, idx)
		}
	}

	for _, comment := range all {
		bar.Increment()

		body := comment.Snippet.TextOriginal
		if compare.ContainsHomoglyphs(body) {
			body = compare.Normalize(body)
		}

		// extract formatted
		for _, w := range extractFormattedWords(body) {
			w = strings.ToLower(w)

			// ignore short words
			if len(w) < FmtSpamMinLen {
				continue
			}

			occ, ok := c.words[w]
			if !ok {
				continue
			}

			printSpamMessage(comment, w, occ)
			c.violations[comment] = 100
		}
	}
	bar.Finish()
	return nil
}
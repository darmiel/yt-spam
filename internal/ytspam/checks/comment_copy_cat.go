package checks

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"time"
)

const (
	CopyCatMinCommentLength = 16
)

type CommentCopyCatCheck struct {
	violations map[*youtube.Comment]ytspam.Rating
	wrapped    []*WrappedCopyCatComment
}

func (c *CommentCopyCatCheck) Name() string {
	return "Copy Cat"
}
func (c *CommentCopyCatCheck) Clean() {
	c.violations = make(map[*youtube.Comment]ytspam.Rating)
	c.wrapped = make([]*WrappedCopyCatComment, 0)
}
func (c *CommentCopyCatCheck) Finalize() map[*youtube.Comment]ytspam.Rating {
	return c.violations
}

type WrappedCopyCatComment struct {
	ID         string
	AuthorID   string
	AuthorName string
	Body       string
	Time       time.Time
	Original   *youtube.Comment
}

func ccPrefix() termenv.Style {
	return termenv.String("🐈 COPY-CAT").Foreground(p().Color("0")).Background(p().Color("#D290E4"))
}

func printCCMessage(oc, ccc *youtube.Comment) {
	b := oc.Snippet.TextOriginal
	if len(b) > 35 {
		b = b[:32] + "..."
	}
	oa := oc.Snippet.AuthorDisplayName
	cc := ccc.Snippet.AuthorDisplayName
	fmt.Println(ccPrefix(),
		termenv.String(cc).Foreground(p().Color("#E88388")),
		"copied",
		termenv.String(oa).Foreground(p().Color("#A8CC8C")),
		"w/",
		termenv.String(b).Foreground(p().Color("#66C2CD")),
		"+", oc.Id, ", -", ccc.Id, "]")
}

func (c *CommentCopyCatCheck) CheckComment(comment *youtube.Comment) error {
	body := comment.Snippet.TextOriginal

	// check body length
	if len(body) < CopyCatMinCommentLength {
		return nil
	}

	authorID := comment.Snippet.AuthorChannelId.Value
	authorName := comment.Snippet.AuthorDisplayName
	t, err := time.Parse(time.RFC3339, comment.Snippet.PublishedAt)
	if err != nil {
		return err
	}

	// find copy catted comments
	for _, ccc := range c.wrapped {
		if ccc.AuthorID == authorID || ccc.ID == comment.Id {
			continue
		}
		if ccc.Body != body {
			continue
		}
		if t.After(ccc.Time) {
			c.violations[comment] = 10
			printCCMessage(ccc.Original, comment)
		} else {
			c.violations[ccc.Original] = 10
			printCCMessage(comment, ccc.Original)
		}
	}

	c.wrapped = append(c.wrapped, &WrappedCopyCatComment{
		ID:         comment.Id,
		AuthorID:   authorID,
		AuthorName: authorName,
		Body:       body,
		Time:       t,
		Original:   comment,
	})
	return nil
}

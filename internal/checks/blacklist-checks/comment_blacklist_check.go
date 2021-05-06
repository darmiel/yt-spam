package blacklist_checks

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/yt-spam/internal/blacklists"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

var p = termenv.ColorProfile()

type CommentBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
}

func (c *CommentBlacklistCheck) Prefix() termenv.Style {
	return termenv.String("✍️ BODY").Foreground(p.Color("0")).Background(p.Color("#71BEF2"))
}

func (c *CommentBlacklistCheck) Name() string {
	return "Comment-Blacklist"
}

func (c *CommentBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	return nil
}

func (c *CommentBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *CommentBlacklistCheck) CheckComments(all map[string]*youtube.Comment) error {
	bar := pb.New(len(all))
	for _, comment := range all {
		bar.Increment()
		body := comment.Snippet.TextOriginal
		bodyNorm := body
		if compare.ContainsHomoglyphs(body) {
			bodyNorm = compare.Normalize(body)
		}

		if cmp := blacklists.CommentBlacklist.AnyAnyMatch(body, bodyNorm); cmp != nil {
			old, ok := c.violations[comment]
			if !ok {
				old = 0
			}
			old += 100
			c.violations[comment] = old

			fmt.Println(c.Prefix(),
				termenv.String(comment.Snippet.AuthorDisplayName).Foreground(p.Color("#E88388")),
				"<<",
				termenv.String(cmp.String()).Foreground(p.Color("#D290E4")))
		}
	}
	bar.Finish()
	return nil
}

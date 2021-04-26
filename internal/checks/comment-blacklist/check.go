package name_blacklist

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"path"
)

type CommentBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
	blacklist  []compare.StringCompare
}

func (c *CommentBlacklistCheck) Name() string {
	return "Comment-Blacklist"
}

func (c *CommentBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	// read blacklist
	var err error
	pa := path.Join("data", "input", "comment-blacklist.txt")
	if c.blacklist, err = compare.FromFile(pa); err != nil {
		return err
	}
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

		for _, b := range c.blacklist {
			normCmp := false
			if body != bodyNorm {
				normCmp = b.Compare(bodyNorm)
			}
			if b.Compare(body) || normCmp {
				old, ok := c.violations[comment]
				if !ok {
					old = 0
				}
				old += 100
				c.violations[comment] = old

				fmt.Println(nbPrefix(),
					termenv.String(comment.Snippet.AuthorDisplayName).Foreground(p.Color("#E88388")),
					"<<",
					termenv.String(b.String()).Foreground(p.Color("#D290E4")))
			}
		}
	}
	bar.Finish()
	return nil
}

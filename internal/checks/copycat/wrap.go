package copycat

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

var p termenv.Profile

func init() {
	p = termenv.ColorProfile()
}

// implements CommentCheck
type CommentCopyCatCheck struct {
	violations map[*youtube.Comment]checks.Rating
}

func (c *CommentCopyCatCheck) Name() string {
	return "Copy Cat"
}

func (c *CommentCopyCatCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	return nil
}

func (c *CommentCopyCatCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *CommentCopyCatCheck) Prefix() termenv.Style {
	return termenv.String("🐈 COPY-CAT").Foreground(p.Color("0")).Background(p.Color("#D290E4"))
}

func (c *CommentCopyCatCheck) printCCMessage(oc, ccc *youtube.Comment) {
	b := trimBody(oc)
	if len(b) > 35 {
		b = b[:32] + "..."
	}
	oa := oc.Snippet.AuthorDisplayName
	cc := ccc.Snippet.AuthorDisplayName
	fmt.Println(c.Prefix(),
		termenv.String(cc).Foreground(p.Color("#E88388")),
		"copied",
		termenv.String(oa).Foreground(p.Color("#A8CC8C")),
		"w/",
		termenv.String(b).Foreground(p.Color("#66C2CD")),
		ccc.Id, "<<", oc.Id)
}

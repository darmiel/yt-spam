package copycat

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/common"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type CommentCopyCatCheck struct {
	channel chan *checks.CommentRatingNotify
	minLen  int
}

func (c *CommentCopyCatCheck) Name() string {
	return "Copy Cat"
}

func (c *CommentCopyCatCheck) Prefix() termenv.Style {
	return common.CreatePrefix("🐈", "COPY-CAT", "D290E4")
}

func (c *CommentCopyCatCheck) SendViolation(i ...interface{}) {
	var (
		original = i[0].(*youtube.Comment)
		copycat  = i[1].(*youtube.Comment)
	)
	b := common.TrimBody(original)
	if len(b) > 35 {
		b = b[:32] + "..."
	}
	oauth := original.Snippet.AuthorDisplayName
	cauth := copycat.Snippet.AuthorDisplayName
	c.channel <- &checks.CommentRatingNotify{
		Reason: fmt.Sprintf("'%s' copied '%s' with '%s' (orig: %s, copy: %s)",
			common.Colored(cauth, "#E88388"),
			common.Colored(oauth, "#A8CC8C"),
			common.Colored(b, "#66C2CD"),
			original.Id, copycat.Id),
		Comment: copycat,
		Rating:  2,
		Check:   c,
	}
}

///

func (c *CommentCopyCatCheck) CheckComments(comments []*youtube.Comment) {
	checked := make(map[string]bool)
	for _, cursor := range comments {
		// check length
		if len(common.TrimBody(cursor)) < c.minLen {
			continue
		}
		if _, checked := checked[cursor.Id]; checked {
			continue
		}
		earliest := cursor
		matches := common.GetMatchingComments(cursor, comments)
		for _, matching := range matches {
			checked[matching.Id] = true
			if common.CommentBefore(matching, earliest) {
				earliest = matching
			}
		}
		for _, matching := range matches {
			if matching != earliest {
				// check if author copied himself
				if common.CommentFromHimHimself(matching, earliest) {
					continue
				}
				c.SendViolation(earliest, matching)
			}
		}
	}
}

package blchecks

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/blacklists"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/common"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type CommentBlacklistCheck struct {
	channel chan *checks.CommentRatingNotify
}

func (c *CommentBlacklistCheck) Prefix() termenv.Style {
	return common.CreatePrefix("📝", "BODY-BL", "71BEF2")
}

func (c *CommentBlacklistCheck) Name() string {
	return "Body-Blacklist"
}

func (c *CommentBlacklistCheck) SendViolation(i ...interface{}) {
	var (
		comment   = i[0].(*youtube.Comment)
		blacklist = i[1].(compare.StringCompare)
	)
	b := common.TrimBody(comment)
	if len(b) > 35 {
		b = b[:32] + "..."
	}
	c.channel <- &checks.CommentRatingNotify{
		Reason: fmt.Sprintf("'%s' used '%s' in '%s'",
			common.Colored(comment.Snippet.AuthorDisplayName, "E88388"),
			common.Colored(blacklist.String(), "DBAB79"),
			common.Colored(b, "66C2CD")),
		Comment: comment,
		Rating:  25,
		Check:   c,
	}
}

func (c *CommentBlacklistCheck) CheckComment(comment *youtube.Comment) {
	body := comment.Snippet.TextOriginal
	bodyNorm := body
	if compare.ContainsHomoglyphs(body) {
		bodyNorm = compare.Normalize(body)
	}
	if cmp := blacklists.CommentBlacklist.AnyAnyMatch(body, bodyNorm); cmp != nil {
		c.SendViolation(comment, cmp)
	}
}

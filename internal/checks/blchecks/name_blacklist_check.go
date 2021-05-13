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

type nameBlacklistCheck struct {
	channel chan *checks.ChannelRatingNotify
	checked map[string]bool
}

func NewNameBlacklistCheck(channel chan *checks.ChannelRatingNotify) *nameBlacklistCheck {
	return &nameBlacklistCheck{channel, make(map[string]bool)}
}

func (c *nameBlacklistCheck) Prefix() termenv.Style {
	return common.CreatePrefix("👦", "NAME-BL", "DBAB79")
}

func (c *nameBlacklistCheck) Name() string {
	return "Name-Blacklist"
}

func (c *nameBlacklistCheck) SendViolation(i ...interface{}) {
	var (
		comment = i[0].(*youtube.Comment)
		cmp     = i[1].(compare.StringCompare)
	)
	c.channel <- &checks.ChannelRatingNotify{
		ChannelName: comment.Snippet.AuthorDisplayName,
		ChannelID:   comment.Snippet.AuthorChannelId.Value,
		Reason: fmt.Sprintf("'%s' contains '%s'",
			common.Colored(comment.Snippet.AuthorDisplayName, "E88388"),
			common.Colored(cmp.String(), "DBAB79")),
		Rating: 25,
		Check:  c,
	}
}

///

func (c *nameBlacklistCheck) CheckChannelByComment(comment *youtube.Comment) {
	authorID := comment.Snippet.AuthorChannelId.Value
	if _, checked := c.checked[authorID]; checked {
		return
	}
	c.checked[authorID] = true
	authorName := comment.Snippet.AuthorDisplayName
	authorNameNorm := authorName
	if compare.ContainsHomoglyphs(authorName) {
		authorNameNorm = compare.Normalize(authorName)
	}
	if cmp := blacklists.NameBlacklist.AnyAnyMatch(authorName, authorNameNorm); cmp != nil {
		c.SendViolation(comment, cmp)
	}
}

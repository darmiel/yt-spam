package blchecks

import (
	"github.com/darmiel/yt-spam/internal/blacklists"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/common"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type ChannelBlacklistCheck struct {
	channel chan *checks.ChannelRatingNotify
}

func (c *ChannelBlacklistCheck) Prefix() termenv.Style {
	return common.CreatePrefix("📺", "CHAN-BL", "DBAB79")
}

func (c *ChannelBlacklistCheck) Name() string {
	return "Channel-Blacklist"
}

func (c *ChannelBlacklistCheck) SendViolation(i ...interface{}) {
	var (
		name = i[0].(string)
		id   = i[1].(string)
	)
	c.channel <- &checks.ChannelRatingNotify{
		ChannelName: name,
		ChannelID:   id,
		Rating:      1000,
		Check:       c,
	}
}

///

func (c *ChannelBlacklistCheck) CheckChannelByComment(comment *youtube.Comment) {
	if i := blacklists.ChannelBlacklist.AnyMatch(comment.Snippet.AuthorChannelId.Value); i != nil {
		c.SendViolation(comment.Snippet.AuthorDisplayName,
			comment.Snippet.AuthorChannelId.Value)
	}
}

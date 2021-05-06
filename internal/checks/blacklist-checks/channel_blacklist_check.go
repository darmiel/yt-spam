package blacklist_checks

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/blacklists"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type ChannelBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
}

func (c *ChannelBlacklistCheck) Prefix() termenv.Style {
	return termenv.String("✍️ CHAN").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

func (c *ChannelBlacklistCheck) Name() string {
	return "Channel-Blacklist"
}

func (c *ChannelBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *ChannelBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	return nil
}

func (c *ChannelBlacklistCheck) CheckChannels(all []*youtube.Channel) error {
	for _, channel := range all {
		if i := blacklists.ChannelBlacklist.AnyMatch(channel.Id); i != nil {
			fmt.Println(c.Prefix(), "Found:", i.String(), "<->", channel.Id, "--", channel.Snippet.Title)
		}
	}
	return nil
}

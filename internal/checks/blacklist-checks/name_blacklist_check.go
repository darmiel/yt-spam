package blacklist_checks

import (
	"fmt"
	"github.com/darmiel/yt-spam/internal/blacklists"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type NameBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
}

func (c *NameBlacklistCheck) Prefix() termenv.Style {
	return termenv.String("✍️ NAME").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

func (c *NameBlacklistCheck) Name() string {
	return "Name-Blacklist"
}

func (c *NameBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)

	return nil
}

func (c *NameBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *NameBlacklistCheck) CheckComments(all map[string]*youtube.Comment) error {
	checked := make(map[string]bool)
	for _, comment := range all {
		authorID := comment.Snippet.AuthorChannelId.Value
		if _, checked := checked[authorID]; checked {
			continue
		}
		authorName := comment.Snippet.AuthorDisplayName
		authorNameNorm := authorName
		if compare.ContainsHomoglyphs(authorName) {
			authorNameNorm = compare.Normalize(authorName)
		}

		if cmp := blacklists.NameBlacklist.AnyAnyMatch(authorName, authorNameNorm); cmp != nil {
			old, ok := c.violations[comment]
			if !ok {
				old = 0
			}
			old += 100
			c.violations[comment] = old

			fmt.Println(c.Prefix(),
				termenv.String(authorName).Foreground(p.Color("#E88388")),
				"<<",
				termenv.String(cmp.String()).Foreground(p.Color("#D290E4")))
		}
	}
	return nil
}

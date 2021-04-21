package name_blacklist

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"log"
	"path"
	"strings"
)

type NameBlacklistCheck struct {
	violations map[*youtube.Comment]checks.Rating
	blacklist  []compare.StringCompare
}

func (c *NameBlacklistCheck) Name() string {
	return "Name-Blacklist"
}

func (c *NameBlacklistCheck) Clean() error {
	c.violations = make(map[*youtube.Comment]checks.Rating)
	// read blacklist
	var err error
	pa := path.Join("data", "input", "name-blacklist.txt")
	if c.blacklist, err = compare.FromFile(pa); err != nil {
		return err
	}
	return nil
}

func (c *NameBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *NameBlacklistCheck) CheckComments(all map[string]*youtube.Comment) error {
	bar := pb.New(len(all))
	checked := make(map[string]bool)
	for _, comment := range all {
		bar.Increment()
		authorID := comment.Snippet.AuthorChannelId.Value
		if _, checked := checked[authorID]; checked {
			continue
		}
		authorName := comment.Snippet.AuthorDisplayName
		authorNameNorm := authorName
		if compare.ContainsHomoglyphs(authorName) {
			log.Println("WARN ::", authorName, "has homoglyphs! Normalizing...")
			authorNameNorm = compare.Normalize(authorName)
			log.Println("  ::", authorNameNorm)
		}

		authorNameLc := strings.ToLower(authorName)

		for _, b := range c.blacklist {
			normCmp := false
			if authorName != authorNameNorm {
				normCmp = b.Compare(authorNameNorm)
			}
			if b.Compare(authorNameLc) || normCmp {
				old, ok := c.violations[comment]
				if !ok {
					old = 0
				}
				old += 100
				c.violations[comment] = old

				fmt.Println(nbPrefix(),
					termenv.String(authorName).Foreground(p.Color("#E88388")),
					"<<",
					termenv.String(b.String()).Foreground(p.Color("#D290E4")))
			}
		}
	}
	bar.Finish()
	return nil
}

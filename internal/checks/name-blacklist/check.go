package name_blacklist

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/yt-spam/internal/checks"
	"github.com/darmiel/yt-spam/internal/compare"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
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
	c.blacklist = make([]compare.StringCompare, 0)
	pa := path.Join("data", "input", "name-blacklist.txt")
	data, err := ioutil.ReadFile(pa)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		w := strings.TrimSpace(strings.ToLower(line))
		if len(w) <= 0 || strings.HasPrefix(w, "#") {
			continue
		}

		var cmp compare.StringCompare

		// regex
		if compare.HasPrefixSuffix(w, "/") {
			fmt.Println(nbPrefix(), "using", w, "as regex")
			cmp = compare.NewStringRegexCompare(w)
		} else {
			cmp = compare.NewStringLowerCompare(w)
		}

		if cmp == nil {
			fmt.Println(nbPrefix(), "skipping", w)
			continue
		}

		c.blacklist = append(c.blacklist, cmp)
		fmt.Println(nbPrefix(), "read word", termenv.String(w).Foreground(p.Color("#E88388")))
	}
	return nil
}

func (c *NameBlacklistCheck) Finalize() map[*youtube.Comment]checks.Rating {
	return c.violations
}

func (c *NameBlacklistCheck) CheckComments(all map[string]*youtube.Comment) error {
	fmt.Println()
	bar := pb.New(len(all))
	checked := make(map[string]bool)
	for _, comment := range all {
		bar.Increment()
		authorID := comment.Snippet.AuthorChannelId.Value
		if _, checked := checked[authorID]; checked {
			continue
		}
		authorName := comment.Snippet.AuthorDisplayName
		authorNameLc := strings.ToLower(authorName)
		for _, b := range c.blacklist {
			if b.Compare(authorNameLc) {
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
	fmt.Println()
	return nil
}

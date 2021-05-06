package fmt_spam

import (
	"fmt"
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
	"regexp"
	"strconv"
)

var pattern = regexp.MustCompile("(?i)(?:\\*(\\w+)\\*)|(?:\\_(\\w+)\\_)|(?:\\~(\\w+)\\~)")

func extractFormattedWords(body string) (res []string) {
	for _, s := range pattern.FindAllStringSubmatch(body, -1) {
		if len(s) != 2 {
			continue
		}
		str := s[1]
		res = append(res, str)
	}
	return
}

var p = termenv.ColorProfile()

func fsPrefix() termenv.Style {
	return termenv.String("🔁 FMT-SPAM").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

func printSpamMessage(comment *youtube.Comment, word string, occurrences uint64) {
	fmt.Println(fsPrefix(),
		termenv.String(comment.Snippet.AuthorDisplayName).Foreground(p.Color("#E88388")),
		"used word",
		termenv.String(word).Foreground(p.Color("#A8CC8C")),
		"w/",
		termenv.String(strconv.FormatUint(occurrences, 10)).Foreground(p.Color("#66C2CD")),
		"occurrences")
}

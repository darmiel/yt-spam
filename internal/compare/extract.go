package compare

import "regexp"

var pattern = regexp.MustCompile("(?i)(?:\\*(\\w+)\\*)|(?:\\_(\\w+)\\_)|(?:\\~(\\w+)\\~)")

func ExtractFormattedWords(body string) (res []string) {
	for _, s := range pattern.FindAllStringSubmatch(body, -1) {
		if len(s) != 2 {
			continue
		}
		str := s[1]
		res = append(res, str)
	}
	return
}

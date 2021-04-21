package compare

import "regexp"

type StringRegexCompare struct {
	needle *regexp.Regexp
	str    string
}

func NewStringRegexCompare(needle string) *StringRegexCompare {
	if HasPrefixSuffix(needle, "/") {
		needle = needle[1 : len(needle)-1]
	}
	r, err := regexp.Compile(needle)
	if err != nil {
		return nil
	}
	return &StringRegexCompare{
		needle: r,
		str:    needle,
	}
}

func (c *StringRegexCompare) Compare(str string) bool {
	return c.needle.MatchString(str)
}

func (c *StringRegexCompare) String() string {
	return c.str
}

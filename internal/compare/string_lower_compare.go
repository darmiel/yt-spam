package compare

import "strings"

func NewStringLowerCompare(needle string) *StringLowerCompare {
	orig := needle
	needle = strings.TrimSpace(strings.ToLower(needle))
	if HasPrefixSuffix(needle, `""`) {
		needle = needle[1 : len(needle)-1]
	}
	return &StringLowerCompare{
		needle: needle,
		str:    orig,
	}
}

type StringLowerCompare struct {
	needle string
	str    string
}

func (c *StringLowerCompare) Compare(str string) bool {
	return strings.Contains(strings.ToLower(str), c.needle)
}

func (c *StringLowerCompare) String() string {
	return c.str
}

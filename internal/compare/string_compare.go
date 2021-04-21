package compare

import "strings"

type StringCompare interface {
	Compare(str string) bool
	String() string
}

///

func HasPrefixSuffix(str, sub string) bool {
	return strings.HasPrefix(str, sub) && strings.HasSuffix(str, sub)
}

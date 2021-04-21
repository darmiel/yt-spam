package compare

import (
	"io/ioutil"
	"log"
	"strings"
)

func FromFile(file string) (res []StringCompare, err error) {
	log.Println("read from file", file)

	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		w := strings.TrimSpace(strings.ToLower(line))
		if len(w) <= 0 || strings.HasPrefix(w, "#") {
			continue
		}

		var cmp StringCompare

		// regex
		if HasPrefixSuffix(w, "/") {
			cmp = NewStringRegexCompare(w)
		} else {
			cmp = NewStringLowerCompare(w)
		}

		if cmp == nil {
			continue
		}

		res = append(res, cmp)
	}

	log.Println("Read", len(res), "string compares")

	return
}

func MustFromFile(file string) []StringCompare {
	res, err := FromFile(file)
	if err != nil {
		panic(err)
	}
	return res
}

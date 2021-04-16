package checks

import (
	"github.com/muesli/termenv"
)

var profile *termenv.Profile

func p() termenv.Profile {
	if profile == nil {
		t := termenv.ColorProfile()
		profile = &t
		return t
	}
	return *profile
}

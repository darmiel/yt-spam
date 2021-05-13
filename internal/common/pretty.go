package common

import (
	"fmt"
	"github.com/muesli/termenv"
)

var profile = termenv.ColorProfile()

func Profile() termenv.Profile {
	return profile
}

// TODO: Background
func CreatePrefix(chr, name, col string) termenv.Style {
	return Colored(fmt.Sprintf("%s %s", chr, name), col)
}

func Colored(str, col string) termenv.Style {
	return termenv.String(str).Foreground(Profile().Color(fmt.Sprintf("#%s", col)))
}

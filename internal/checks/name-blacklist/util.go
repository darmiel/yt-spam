package name_blacklist

import "github.com/muesli/termenv"

var p = termenv.ColorProfile()

func nbPrefix() termenv.Style {
	return termenv.String("✍️ NAME").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

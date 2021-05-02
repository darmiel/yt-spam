package cmt_blacklist

import "github.com/muesli/termenv"

var p = termenv.ColorProfile()

func nbPrefix() termenv.Style {
	return termenv.String("✍️ BODY").Foreground(p.Color("0")).Background(p.Color("#71BEF2"))
}

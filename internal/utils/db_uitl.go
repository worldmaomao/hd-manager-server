package utils

import "regexp"

var (
	SpecCharRegexpCompile, _ = regexp.Compile(`([\*\.\?\+\$\^\[\]\(\)\{\}\|\\\/]+)`)
)

func FilterSpecChar(text string) string {
	return SpecCharRegexpCompile.ReplaceAllString(text, `\${1}`)
}

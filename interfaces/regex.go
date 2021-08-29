package interfaces

import (
	"regexp"
	"strings"
)

var AsciiRegex = regexp.MustCompile("[[:^ascii:]]")

func PrepareStringToBeUsedForHtmlId(text string) string {
	text = AsciiRegex.ReplaceAllLiteralString(text, "")
	if len(text) > 30 {
		text = text[:30]
	}
	text = strings.Replace(strings.ToLower(text), " ", "_", -1)
	text = strings.Replace(strings.ToLower(text), ".", "_", -1)
	return text
}

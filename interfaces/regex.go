package interfaces

import "regexp"

var AsciiRegex = regexp.MustCompile("[[:^ascii:]]")


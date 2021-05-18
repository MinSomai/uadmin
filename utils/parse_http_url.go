package utils

import (
	"regexp"
)

var (
	UrlRemainderRegex, _ = regexp.Compile(`(?P<Path>[a-zA-Z0-9]+)`)
)

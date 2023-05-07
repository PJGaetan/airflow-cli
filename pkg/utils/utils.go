package utils

import (
	"regexp"
	"strings"
)

var lineRE = regexp.MustCompile(`(?m)^`)

func Indent(s, indent string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}

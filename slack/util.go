package slack

import (
	"regexp"
	"strings"
)

var (
	replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

	escapedUserReference = regexp.MustCompile("&lt;(@[A-Z0-9]+\\|[A-Za-z0-9\\._-]+)&gt;") // escaped user reference, i.e. <@U123456|username>
)

func EscapeMessage(s string) string {
	escaped := []byte(replacer.Replace(s))

	// restore escaped user references
	escaped = escapedUserReference.ReplaceAll(escaped, []byte("<${1}>"))

	return string(escaped)
}

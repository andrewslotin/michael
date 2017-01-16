package slack

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

	escapedLink          = regexp.MustCompile("&lt;(https?://\\S+?)&gt;")                 // escaped URL, i.e. <https://google.com?q=search+term>
	escapedUserReference = regexp.MustCompile("&lt;(@[A-Z0-9]+\\|[A-Za-z0-9\\._-]+)&gt;") // escaped user reference, i.e. <@U123456|username>
)

func EscapeMessage(s string) string {
	escaped := []byte(replacer.Replace(s))

	// restore escaped user references
	escaped = escapedUserReference.ReplaceAll(escaped, []byte("<${1}>"))

	// restore escaped URLs
	escaped = escapedLink.ReplaceAllFunc(escaped, func(m []byte) []byte {
		m[3] = '<'
		m[len(m)-4] = '>'
		m = m[3 : len(m)-3] // strip &lt at the beginning and gt; at the endj

		return bytes.Replace(m, []byte("&amp;"), []byte{'&'}, -1) // restore &amp; -> &
	})

	return string(escaped)
}

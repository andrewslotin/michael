package slack

import "strings"

var replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

func EscapeMessage(s string) string {
	return replacer.Replace(s)
}

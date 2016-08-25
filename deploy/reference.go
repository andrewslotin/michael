package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	referenceRegexes = []*regexp.Regexp{
		regexp.MustCompile("^(?P<repository>\\S+/\\S+)#(?P<number>\\d+)[^A-Za-z]?$"),                        // octocat/helloworld#12
		regexp.MustCompile("^https?://github.com/(?P<repository>\\S+/\\S+)/pull/(?P<number>\\d+)(?:\\?|$)"), // https://github.com/octocat/helloworld/pull/12
	}
)

type Reference struct {
	ID         string
	Repository string
}

func FindReferences(s string) []Reference {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	var refs []Reference
	for scanner.Scan() {
		word := scanner.Text()

		for _, re := range referenceRegexes {
			if ref, ok := extractReference(word, re); ok {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}

func extractReference(s string, re *regexp.Regexp) (ref Reference, ok bool) {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return ref, false
	}

	for i, name := range re.SubexpNames() {
		switch name {
		case "repository":
			ref.Repository = m[i]
		case "number":
			ref.ID = m[i]
		default:
			continue
		}
	}

	return ref, true
}

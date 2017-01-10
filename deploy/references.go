package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	// TODO: support anchors in GitHub URLs along with query string parameters
	pullRequestReferenceRegexes = []*regexp.Regexp{
		regexp.MustCompile("^(?P<repository>\\S+/\\S+)#(?P<number>\\d+)[^A-Za-z]?$"),                        // octocat/helloworld#12
		regexp.MustCompile("^https?://github.com/(?P<repository>\\S+/\\S+)/pull/(?P<number>\\d+)(?:\\?|$)"), // https://github.com/octocat/helloworld/pull/12
	}
)

type PullRequestReference struct {
	ID         string
	Repository string
}

func FindPullRequestReferences(s string) []PullRequestReference {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	var refs []PullRequestReference
	for scanner.Scan() {
		word := scanner.Text()

		for _, re := range pullRequestReferenceRegexes {
			if ref, ok := extractPullRequestReference(word, re); ok {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}

func extractPullRequestReference(s string, re *regexp.Regexp) (ref PullRequestReference, ok bool) {
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

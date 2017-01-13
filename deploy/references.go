package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	pullRequestReferenceRegexes = []*regexp.Regexp{
		regexp.MustCompile("^(?P<repository>\\S+/\\S+)#(?P<number>\\d+)[^A-Za-z]?$"),                          // octocat/helloworld#12
		regexp.MustCompile("^https?://github.com/(?P<repository>\\S+/\\S+)/pull/(?P<number>\\d+)(?:\\?|#|$)"), // https://github.com/octocat/helloworld/pull/12
	}
	userReferenceRegexes = []*regexp.Regexp{
		// Usernames can be up to 21 characters long. They can contain lowercase letters a to z (without accents),
		// numbers 0 to 9, hyphens, periods, and underscores.
		//
		// See https://get.slack.help/hc/en-us/articles/216360827-Change-your-username
		regexp.MustCompile("^@(?P<username>[A-Za-z0-9\\._-]+)"),
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

type UserReference struct {
	Name string
}

func FindUserReferences(s string) []UserReference {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	var refs []UserReference
	for scanner.Scan() {
		word := scanner.Text()

		for _, re := range userReferenceRegexes {
			if ref, ok := extractUserReference(word, re); ok {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}

func extractUserReference(s string, re *regexp.Regexp) (ref UserReference, ok bool) {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return ref, false
	}

	for i, name := range re.SubexpNames() {
		switch name {
		case "username":
			ref.Name = m[i]
		default:
			continue
		}
	}

	return ref, true
}

package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	pullRequestReferenceRegexes = []*regexp.Regexp{
		regexp.MustCompile("^(?P<repository>[A-Za-z0-9\\._-]+/[A-Za-z0-9\\._-]+)#(?P<number>\\d+)[^A-Za-z]?$"),    // octocat/helloworld#12
		regexp.MustCompile("^<?https?://github.com/(?P<repository>\\S+/\\S+)/pull/(?P<number>\\d+)(?:[\\?#>]|$)"), // https://github.com/octocat/helloworld/pull/12
	}
	userReferenceRegexes = []*regexp.Regexp{
		// Usernames can be up to 21 characters long. They can contain lowercase letters a to z (without accents),
		// numbers 0 to 9, hyphens, periods, and underscores.
		//
		// See https://get.slack.help/hc/en-us/articles/216360827-Change-your-username
		regexp.MustCompile("^@(?P<username>[A-Za-z0-9\\._-]+)"),                           // unescaped, i.e. @user1
		regexp.MustCompile("^<@(?P<userid>[A-Z0-9]+)\\|(?P<username>[A-Za-z0-9\\._-]+)>"), // escaped, i.e. <@U1|user1>
	}
)

type PullRequestReference struct {
	ID         string
	Repository string
}

func FindPullRequestReferences(s string) []PullRequestReference {
	var refs []PullRequestReference
	findReferences(s, pullRequestReferenceRegexes, func(matches map[string]string) {
		refs = append(refs, PullRequestReference{Repository: matches["repository"], ID: matches["number"]})
	})

	return refs
}

type UserReference struct {
	ID   string
	Name string
}

func FindUserReferences(s string) []UserReference {
	var refs []UserReference
	findReferences(s, userReferenceRegexes, func(matches map[string]string) {
		refs = append(refs, UserReference{ID: matches["userid"], Name: matches["username"]})
	})

	return refs
}

func findReferences(s string, regexes []*regexp.Regexp, f func(map[string]string)) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()

		for _, re := range regexes {
			m := re.FindStringSubmatch(word)
			if m == nil {
				continue
			}

			matches := make(map[string]string, len(re.SubexpNames()))
			for i, name := range re.SubexpNames() {
				matches[name] = m[i]
			}

			f(matches)
		}
	}
}

package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	pullRequestReferenceRegex = regexp.MustCompile("^(?P<repository>\\S+/\\S+)#(?P<number>\\d+)[^A-Za-z]?$")
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
		if m := pullRequestReferenceRegex.FindStringSubmatch(word); len(m) > 0 {
			refs = append(refs, Reference{
				ID:         m[2],
				Repository: m[1],
			})
		}
	}

	return refs
}

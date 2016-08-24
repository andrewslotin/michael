package deploy

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	referenceRegexes = []*regexp.Regexp{
		regexp.MustCompile("^(?P<repository>\\S+/\\S+)#(?P<number>\\d+)[^A-Za-z]?$"),
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
			m := re.FindStringSubmatch(word)
			if m == nil {
				continue
			}

			var ref Reference
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
			refs = append(refs, ref)
		}
	}

	return refs
}

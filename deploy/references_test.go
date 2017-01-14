package deploy_test

import (
	"testing"

	"github.com/andrewslotin/michael/deploy"
	"github.com/stretchr/testify/assert"
)

func TestFindPullRequestReferences_Short(t *testing.T) {
	s := "user project 123 user/ /project user/#123 /project#123 /# #/123 use/project user/project#123 user/project#1a"

	refs := deploy.FindPullRequestReferences(s)
	if assert.Len(t, refs, 1) {
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "123", Repository: "user/project"})
	}
}

func TestFindPullRequestReferences_GitHubLink(t *testing.T) {
	s := "" +
		"https://github.com/user/project/pull/1?w=1#comment-123 " +
		"https://github.com/user/project/issues/2 " +
		"https://github.com/user/project/pulls " +
		"https://bitbucket.org/user/project/pull/3"

	refs := deploy.FindPullRequestReferences(s)
	if assert.Len(t, refs, 1) {
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "1", Repository: "user/project"})
	}
}

func TestFindPullRequestReferences_Escaped(t *testing.T) {
	s := "<https://github.com/user/project/pull/1?w=1#comment-123> <user/project#123>"

	refs := deploy.FindPullRequestReferences(s)
	if assert.Len(t, refs, 1) {
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "1", Repository: "user/project"})
	}
}

func TestFindPullRequestReferences_Mixed_Multiple(t *testing.T) {
	s := "" +
		"userA/projectA#1, " +
		"https://github.com/userB/projectB/pull/2 " +
		"and userC/projectC#3"

	refs := deploy.FindPullRequestReferences(s)
	if assert.Len(t, refs, 3) {
		// userA/projectA#1
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "1", Repository: "userA/projectA"})

		// userB/projectB#2
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "2", Repository: "userB/projectB"})

		// userC/projectC#3
		assert.Contains(t, refs, deploy.PullRequestReference{ID: "3", Repository: "userC/projectC"})
	}
}

func TestFindUserReferences(t *testing.T) {
	s := "" +
		"hello @person_1, my email is writeme@gmail.com, see you @ the bar. " +
		"if you see @person.2 please send him to @me"

	refs := deploy.FindUserReferences(s)
	if assert.Len(t, refs, 3) {
		assert.Contains(t, refs, deploy.UserReference{Name: "person_1"})
		assert.Contains(t, refs, deploy.UserReference{Name: "person.2"})
		assert.Contains(t, refs, deploy.UserReference{Name: "me"})
	}
}

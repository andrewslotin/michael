package deploy_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindReferences(t *testing.T) {
	s := "user project 123 user/ /project user/#123 /project#123 /# #/123 use/project user/project#123 user/project#1a"

	refs := deploy.FindReferences(s)
	require.Len(t, refs, 1)

	ref := refs[0]
	assert.Equal(t, "123", ref.ID)
	assert.Equal(t, "user/project", ref.Repository)
}

func TestFindReferences_Multiple(t *testing.T) {
	s := "userA/projectA#1, userB/projectB#2 and userC/projectC#3"

	refs := deploy.FindReferences(s)
	require.Len(t, refs, 3)

	// userA/projectA#1
	assert.Equal(t, "1", refs[0].ID)
	assert.Equal(t, "userA/projectA", refs[0].Repository)

	// userB/projectB#2
	assert.Equal(t, "2", refs[1].ID)
	assert.Equal(t, "userB/projectB", refs[1].Repository)

	// userC/projectC#3
	assert.Equal(t, "3", refs[2].ID)
	assert.Equal(t, "userC/projectC", refs[2].Repository)
}

package auth_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/andrewslotin/slack-deploy-command/auth/authtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneTimeTokenAuthorizer_IssueToken(t *testing.T) {
	src := new(authtest.TokenSourceMock)
	src.On("Generate", 10).Return("abcdef1234")
	src.On("Generate", 5).Return("xyz12")

	authorizer := auth.NewOneTimeTokenAuthorizer(src)

	if token, err := authorizer.IssueToken(10); assert.NoError(t, err) {
		assert.Equal(t, "abcdef1234", token)
	}

	if token, err := authorizer.IssueToken(5); assert.NoError(t, err) {
		assert.Equal(t, "xyz12", token)
	}

	src.AssertExpectations(t)
}

func TestOneTimeTokenAuthorizer_IssueToken_Uniqueness(t *testing.T) {
	authorizer := auth.NewOneTimeTokenAuthorizer(authtest.StaticTokenSource("token1"))

	token, err := authorizer.IssueToken(1)
	require.NoError(t, err)
	require.Equal(t, "token1", token)

	_, err = authorizer.IssueToken(1)
	assert.Error(t, err)
}

func TestOneTimeTokenAuthorizer_Authorize(t *testing.T) {
	authorizer := auth.NewOneTimeTokenAuthorizer(authtest.StaticTokenSource("token1"))
	assert.False(t, authorizer.Authorize("token1"))

	token, err := authorizer.IssueToken(1)
	require.NoError(t, err)
	require.Equal(t, "token1", token)

	assert.True(t, authorizer.Authorize(token))
	assert.False(t, authorizer.Authorize(token))
}

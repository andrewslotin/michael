package auth_test

import (
	"testing"

	"github.com/andrewslotin/michael/auth"
	"github.com/andrewslotin/michael/auth/authtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneTimeTokenAuthenticator_IssueToken(t *testing.T) {
	src := new(authtest.TokenSourceMock)
	src.On("Generate", 10).Return("abcdef1234")
	src.On("Generate", 5).Return("xyz12")

	authenticator := auth.NewOneTimeTokenAuthenticator(src)

	if token, err := authenticator.IssueToken(10); assert.NoError(t, err) {
		assert.Equal(t, "abcdef1234", token)
	}

	if token, err := authenticator.IssueToken(5); assert.NoError(t, err) {
		assert.Equal(t, "xyz12", token)
	}

	src.AssertExpectations(t)
}

func TestOneTimeTokenAuthenticator_IssueToken_Uniqueness(t *testing.T) {
	authenticator := auth.NewOneTimeTokenAuthenticator(authtest.StaticTokenSource("token1"))

	token, err := authenticator.IssueToken(1)
	require.NoError(t, err)
	require.Equal(t, "token1", token)

	_, err = authenticator.IssueToken(1)
	assert.Error(t, err)
}

func TestOneTimeTokenAuthenticator_Authenticate(t *testing.T) {
	authenticator := auth.NewOneTimeTokenAuthenticator(authtest.StaticTokenSource("token1"))
	assert.False(t, authenticator.Authenticate("token1"))

	token, err := authenticator.IssueToken(1)
	require.NoError(t, err)
	require.Equal(t, "token1", token)

	assert.True(t, authenticator.Authenticate(token))
	assert.False(t, authenticator.Authenticate(token))
}

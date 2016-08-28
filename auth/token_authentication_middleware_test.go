package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/andrewslotin/slack-deploy-command/auth/authtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthenticationMiddleware_ValidToken(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?token="+token, nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	handler.On("ServeHTTP", recorder, req).Return().Once()
	authenticator.On("Authenticate", token).Return(true)

	auth.TokenAuthenticationMiddleware(handler, authenticator).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	handler.AssertExpectations(t)
	authenticator.AssertExpectations(t)
}

func TestTokenAuthenticationMiddleware_InvalidToken(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?token="+token, nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	authenticator.On("Authenticate", token).Return(false)

	auth.TokenAuthenticationMiddleware(handler, authenticator).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), strings.TrimSpace(recorder.Body.String()))

	handler.AssertExpectations(t)
	authenticator.AssertExpectations(t)
}

func TestTokenAuthenticationMiddleware_MissingToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)

	auth.TokenAuthenticationMiddleware(handler, authenticator).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "Missing token", strings.TrimSpace(recorder.Body.String()))

	handler.AssertExpectations(t)
}

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andrewslotin/michael/auth"
	"github.com/andrewslotin/michael/auth/authtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthenticationMiddleware_ValidToken(t *testing.T) {
	secret := []byte("secret")
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/channel1?token="+token, nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	authenticator.On("Authenticate", token).Return(true)

	auth.TokenAuthenticationMiddleware(handler, authenticator, secret).ServeHTTP(recorder, req)

	authenticator.AssertExpectations(t)

	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Equal(t, "/channel1", recorder.Header().Get("Location"))

	if cookieString := recorder.Header().Get("Set-Cookie"); assert.Regexp(t, "^Auth=[^;]+[;$]", cookieString) {
		parts := strings.SplitN(cookieString, ";", 2)

		token := strings.TrimPrefix(parts[0], "Auth=")
		if claims, err := auth.ParseChannelAccessTokenClaims(token, secret); assert.NoError(t, err) {
			assert.WithinDuration(t, time.Now().Add(auth.ChannelAccessTokenExpirationPeriod), claims.Channels["channel1"], time.Second)
		}
	}
}

func TestTokenAuthenticationMiddleware_InvalidToken(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/channel1?token="+token, nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	handler.On("ServeHTTP", recorder, req).Return().Once()
	authenticator.On("Authenticate", token).Return(false)

	auth.TokenAuthenticationMiddleware(handler, authenticator, []byte("secret")).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	assert.Empty(t, recorder.Header().Get("Set-Cookie"))
	_, err = req.Cookie("Auth")
	assert.Equal(t, err, http.ErrNoCookie)

	handler.AssertExpectations(t)
	authenticator.AssertExpectations(t)
}

func TestTokenAuthenticationMiddleware_NoChannelID(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?token="+token, nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	handler.On("ServeHTTP", recorder, req).Return().Once()

	auth.TokenAuthenticationMiddleware(handler, authenticator, []byte("secret")).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	assert.Empty(t, recorder.Header().Get("Set-Cookie"))
	_, err = req.Cookie("Auth")
	assert.Equal(t, err, http.ErrNoCookie)

	handler.AssertExpectations(t)
}

func TestTokenAuthenticationMiddleware_NoToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	var (
		handler       authtest.HandlerMock
		authenticator authtest.TokenAuthenticatorMock
	)
	handler.On("ServeHTTP", recorder, req).Return().Once()

	auth.TokenAuthenticationMiddleware(handler, authenticator, nil).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	assert.Empty(t, recorder.Header().Get("Set-Cookie"))
	_, err = req.Cookie("Auth")
	assert.Equal(t, err, http.ErrNoCookie)

	handler.AssertExpectations(t)
}

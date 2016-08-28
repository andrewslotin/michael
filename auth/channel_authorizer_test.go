package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/andrewslotin/slack-deploy-command/auth/authtest"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelAuthorizerMiddleware_ValidToken_WithChannelAccess(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	var handler authtest.HandlerMock
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	handler.On("ServeHTTP", recorder, req).Return().Once()

	auth.ChannelAuthorizerMiddleware(handler, jwtSecret).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	if !handler.AssertExpectations(t) {
		t.Logf("Response: %q", recorder.Body)
	}
}

func TestChannelAuthorizerMiddleware_ValidToken_NoChannelAccess(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Channels: map[string]time.Time{
			"channel2": time.Now().Add(time.Hour),
		},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	var handler authtest.HandlerMock
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(handler, jwtSecret).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "No channel access", strings.TrimSpace(recorder.Body.String()))
}

func TestChannelAuthorizerMiddleware_ValidToken_ExpiredChannelAccess(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(-1 * time.Minute),
		},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	var handler authtest.HandlerMock
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(handler, jwtSecret).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "Channel access expired", strings.TrimSpace(recorder.Body.String()))
}

func TestChannelAuthorizerMiddleware_ExpiredToken(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(-1 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(authtest.HandlerMock{}, jwtSecret).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "Token is expired", strings.TrimSpace(recorder.Body.String()))
}

func TestChannelAuthorizerMiddleware_SignatureMismatch(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(1 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(authtest.HandlerMock{}, []byte("a very secret secret")).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "Response: %q", recorder.Body)
}

func TestChannelAuthorizerMiddleware_UnexpectedSigningMethod(t *testing.T) {
	pkey, err := rsa.GenerateKey(rand.Reader, 512)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	})

	signedToken, err := token.SignedString(pkey)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/channel1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(authtest.HandlerMock{}, []byte("test secret")).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Response: %q", recorder.Body)
}

func TestChannelAuthorizerMiddleware_NoChannelID(t *testing.T) {
	jwtSecret := []byte("test secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Channels: map[string]time.Time{},
	})

	signedToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	var handler authtest.HandlerMock
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: signedToken,
	})

	auth.ChannelAuthorizerMiddleware(handler, jwtSecret).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

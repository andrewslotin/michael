package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelAccessTokenFromRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "Auth",
		Value: "cookie1",
	})
	assert.Equal(t, "cookie1", auth.ChannelAccessTokenFromRequest(req))
}

func TestChannelAccessTokenFromRequest_NoAuthCookie(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	assert.Equal(t, "", auth.ChannelAccessTokenFromRequest(req))
}

func TestParseChannelAccessTokenClaims(t *testing.T) {
	secret := []byte("test secret")

	claims := auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	require.NoError(t, err)

	parsedClaims, err := auth.ParseChannelAccessTokenClaims(tokenString, secret)
	require.NoError(t, err)

	assert.Equal(t, claims.IssuedAt, parsedClaims.IssuedAt)
	assert.Equal(t, claims.ExpiresAt, parsedClaims.ExpiresAt)
	if assert.Len(t, parsedClaims.Channels, 1) {
		assert.WithinDuration(t, claims.Channels["channel1"].UTC(), parsedClaims.Channels["channel1"].UTC(), time.Second)
	}
}

func TestParseChannelAccessTokenClaims_ExpiredToken(t *testing.T) {
	secret := []byte("test secret")

	claims := auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-10 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(-1 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	require.NoError(t, err)

	_, err = auth.ParseChannelAccessTokenClaims(tokenString, secret)
	assert.Equal(t, auth.ErrExpiredToken, err)
}

func TestParseChannelAccessTokenClaims_SignatureMismatch(t *testing.T) {
	secret := []byte("test secret")

	claims := auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(1 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("a very secret secret"))
	require.NoError(t, err)

	_, err = auth.ParseChannelAccessTokenClaims(tokenString, secret)
	assert.Equal(t, auth.ErrInvalidToken, err)
}

func TestParseChannelAccessTokenClaims_UnexpectedSigningMethod(t *testing.T) {
	pkey, err := rsa.GenerateKey(rand.Reader, 512)
	require.NoError(t, err)

	claims := auth.JWTChannelClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		Channels: map[string]time.Time{
			"channel1": time.Now().Add(time.Hour),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(pkey)
	require.NoError(t, err)

	_, err = auth.ParseChannelAccessTokenClaims(tokenString, pkey)
	assert.Equal(t, auth.ErrInvalidSigningMethod, err)
}

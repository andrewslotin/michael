package auth

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// ChannelAccessTokenExpirationPeriod is the default channel access expiration period.
// Once granted channel access needs to be renewed after 30 days.
const ChannelAccessTokenExpirationPeriod = 30 * 24 * time.Hour

// ChannelAccessTokenFromRequest reads and returns signed JWT from request. If the request doesn't contain
// access token this method returns an empty string.
func ChannelAccessTokenFromRequest(r *http.Request) string {
	authCookie, err := r.Cookie("Auth")
	if err != nil {
		return ""
	}

	return authCookie.Value
}

// ParseChannelAccessTokenClaims verifies and parses signed JWT string and returns encoded JWTChannelClaims.
// Most of the time the returned error is of type auth.Error.
func ParseChannelAccessTokenClaims(tokenString string, key interface{}) (claims *JWTChannelClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTChannelClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return key, nil
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			switch {
			case validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				return nil, ErrExpiredToken
			case validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, ErrInvalidToken
			case validationErr.Errors&jwt.ValidationErrorUnverifiable != 0:
				return nil, ErrInvalidSigningMethod
			}
		}

		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTChannelClaims)
	if !ok {
		return nil, ErrInvalidTokenFormat
	}

	return claims, nil
}

// StoreChannelAccessToken writes tokenString into Auth= cookie expiring in expTime.
func StoreChannelAccessToken(w http.ResponseWriter, tokenString string, expTime time.Time) {
	cookie := &http.Cookie{
		Name:    "Auth",
		Value:   tokenString,
		Expires: expTime,
	}

	http.SetCookie(w, cookie)
}

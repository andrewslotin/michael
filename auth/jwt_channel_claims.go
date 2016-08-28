package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTChannelClaims struct {
	jwt.StandardClaims

	Channels map[string]time.Time `json:"channels"`
}

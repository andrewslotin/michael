package auth

import "net/http"

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrInvalidSigningMethod = Error{Message: "Unsupported token signing method", Code: http.StatusBadRequest}
	ErrExpiredToken         = Error{Message: "Token is expired", Code: http.StatusUnauthorized}
	ErrInvalidToken         = Error{Message: "Invalid token", Code: http.StatusUnauthorized}
	ErrInvalidTokenFormat   = Error{Message: "Invalid token format", Code: http.StatusBadRequest}
	ErrNoChannelAccess      = Error{Message: "No channel access", Code: http.StatusUnauthorized}
	ErrExpiredChannelAccess = Error{Message: "Channel access expired", Code: http.StatusUnauthorized}
)

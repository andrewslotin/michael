package auth

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ChannelAuthorizer struct {
	handler http.Handler
	secret  []byte
}

// ChannelAuthorizerMiddleware calls an undelying http.Handler once and only there is a valid JWT
// provided in Authorization header.
func ChannelAuthorizerMiddleware(h http.Handler, jwtSecret []byte) *ChannelAuthorizer {
	return &ChannelAuthorizer{
		handler: h,
		secret:  jwtSecret,
	}
}

func (h *ChannelAuthorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var channelID string
	if len(r.URL.Path) > 1 {
		channelID = r.URL.Path[1:]
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	tokenString := h.tokenFromRequest(r)
	if tokenString == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := h.checkAccess(channelID, tokenString)
	if err != nil {
		if authError, ok := err.(Error); ok {
			http.Error(w, authError.Message, authError.Code)
		} else {
			log.Printf("failed to check channel access: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	if h.handler != nil {
		h.handler.ServeHTTP(w, r)
	} else {
		w.Write([]byte("OK"))
	}
}

func (h *ChannelAuthorizer) tokenFromRequest(r *http.Request) string {
	authCookie, err := r.Cookie("Auth")
	if err != nil {
		return ""
	}

	return authCookie.Value
}

func (h *ChannelAuthorizer) checkAccess(channelID, signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &JWTChannelClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return h.secret, nil
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			switch {
			case validationErr.Errors&jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet != 0:
				return ErrExpiredToken
			}
		}

		return err
	}
	if !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTChannelClaims)
	if !ok {
		return ErrInvalidTokenFormat
	}

	expiresAt, ok := claims.Channels[channelID]
	switch {
	case !ok:
		return ErrNoChannelAccess
	case time.Now().After(expiresAt):
		return ErrExpiredChannelAccess
	default:
		return nil
	}
}

package auth

import (
	"net/http"
	"time"

	"github.com/andrewslotin/slack-deploy-command/dashboard"
	jwt "github.com/dgrijalva/jwt-go"
)

type ChannelAuthenticator struct {
	handler http.Handler
	auth    TokenAuthenticator
	secret  []byte
}

// TokenAuthenticationMiddleware wraps an http.Handler and checks if the request contains token parameter
// which value can be authenticated by given authenticator. If the token is authenticated CahnnelAuthenticator
// grants access to requested channel. If there was no token provided, the request gets passed further leaving
// the underlying handler to deal with authorization.
func TokenAuthenticationMiddleware(h http.Handler, authenticator TokenAuthenticator, jwtSecret []byte) http.Handler {
	return &ChannelAuthenticator{
		handler: h,
		auth:    authenticator,
		secret:  jwtSecret,
	}
}

func (h *ChannelAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelID := dashboard.ChannelIDFromRequest(r)
	if channelID == "" {
		h.handler.ServeHTTP(w, r)
		return
	}

	if token := r.FormValue("token"); token == "" || !h.auth.Authenticate(token) {
		h.handler.ServeHTTP(w, r)
		return
	}

	h.grantChannelAccess(channelID, w, r)

	url := *r.URL
	q := url.Query()
	q.Del("token")
	url.RawQuery = q.Encode()

	http.Redirect(w, r, url.String(), http.StatusFound)
}

func (h *ChannelAuthenticator) grantChannelAccess(channelID string, w http.ResponseWriter, r *http.Request) error {
	if channelID == "" {
		return nil
	}

	var claims *JWTChannelClaims
	if token := ChannelAccessTokenFromRequest(r); token != "" {
		if existingClaims, err := ParseChannelAccessTokenClaims(token, h.secret); err == nil { // TODO: inject key
			claims = existingClaims
		}
	}

	if claims == nil {
		claims = new(JWTChannelClaims)
		claims.Channels = make(map[string]time.Time, 1)
	}

	issueTime := time.Now()
	expirationTime := issueTime.Add(ChannelAccessTokenExpirationPeriod)

	claims.IssuedAt = issueTime.Unix()
	claims.ExpiresAt = expirationTime.Unix()
	claims.Channels[channelID] = expirationTime

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.secret)
	if err != nil {
		return err
	}

	StoreChannelAccessToken(w, tokenString, expirationTime)

	return nil
}

package auth

import "net/http"

// TokenAuthenticationMiddleware wraps an http.Handler and checks if the request contains token parameter
// which value can be authenticated by given authenticator. If token cannot be authenticated TokenAuthenticationMiddleware
// responds with HTTP 403 Unautorized without calling the wrapped handler.
func TokenAuthenticationMiddleware(h http.Handler, authenticator TokenAuthenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		if !authenticator.Authenticate(token) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

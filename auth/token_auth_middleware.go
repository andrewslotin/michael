package auth

import "net/http"

// TokenAuthMiddleware wraps an http.Handler and checks if the request contains token parameter
// which value can be authorized by given authorizer. If token cannot be authorized TokenAuthMiddleware
// responds with HTTP 403 Unautorized without calling the wrapped handler.
func TokenAuthMiddleware(h http.Handler, authorizer TokenAuthorizer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		if !authorizer.Authorize(token) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

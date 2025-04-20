package sesh

import (
	"context"
	"net/http"
)

// SessionChecker checks if a session cookie exists, gets the session if it exists and adds it to the request context.
func SessionChecker[T any](next http.Handler, store *sessionStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the session if it exists
		var data T
		err := store.GetWithCookie(w, r, &data)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// create a new request with the the session data added to the context
		r = r.WithContext(context.WithValue(r.Context(), store.config.ContextName, data))
		next.ServeHTTP(w, r)
	})
}

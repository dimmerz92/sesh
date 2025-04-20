package sesh

import (
	"net/http"
	"time"
)

// NewHttp creates a new session, sets a new session cookie on the response writer, and returns the session ID.
func (s *sessionStore) NewWithCookie(w http.ResponseWriter, data any) (string, error) {
	// create a new session
	sessionId, err := s.New(data)
	if err != nil {
		return "", err
	}

	// add a session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     s.config.CookieName,
		Path:     s.config.CookiePath,
		Value:    sessionId,
		Expires:  time.Now().Add(s.config.SessionLength),
		HttpOnly: s.config.CookieHttpOnly,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
	})

	return sessionId, nil
}

// Deletes the session from the session store and invalidates the session cookie on the response writer.
func (s *sessionStore) DeleteWithCookie(w http.ResponseWriter, r *http.Request) error {
	// get the session ID
	cookie, err := r.Cookie(s.config.CookieName)
	if err != nil {
		return err
	}

	// delete the session
	err = s.Delete(cookie.Value)
	if err != nil {
		return err
	}

	// invalidate the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     s.config.CookieName,
		Path:     s.config.CookiePath,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: s.config.CookieHttpOnly,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
	})

	return nil
}

func (s *sessionStore) GetWithCookie(w http.ResponseWriter, r *http.Request, v any) error {
	// get the session ID
	cookie, err := r.Cookie(s.config.CookieName)
	if err != nil {
		return err
	}

	// get the session
	err = s.Get(cookie.Value, v)
	if err != nil {
		return err
	}

	// add an updated cookie
	http.SetCookie(w, &http.Cookie{
		Name:     s.config.CookieName,
		Path:     s.config.CookiePath,
		Value:    cookie.Value,
		Expires:  time.Now().Add(s.config.SessionLength),
		HttpOnly: s.config.CookieHttpOnly,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
	})

	return nil
}

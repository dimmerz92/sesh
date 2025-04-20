package sesh

import (
	"net/http"
	"time"
)

type Config struct {
	// the directory in which the persistent session store is located.
	Dir string
	// if true, the session store will only be in memory, not persistent.
	InMemory bool

	// the duration of a single session.
	SessionLength time.Duration
	// automatically renew sessions on successful validation.
	ExtendSessions bool

	// the name of the session cookie
	CookieName string
	// the path of the cookie
	CookiePath string
	// if true, the cookie is designated http only
	CookieHttpOnly bool
	// if true, the cookie is designated secure only
	CookieSecure bool
	// sets the same site protocol for the cookie
	CookieSameSite http.SameSite

	// sets the name of the session on the context
	ContextName string
}

// DefaultConfig returns a session store config with sensible defaults.
// Modify these to suit using the `WithX` chain methods.
func DefaultConfig() Config {
	return Config{
		Dir: "./session_data",

		SessionLength:  time.Hour,
		ExtendSessions: true,

		CookieName:     "session",
		CookiePath:     "/",
		CookieHttpOnly: true,
		CookieSecure:   true,
		CookieSameSite: http.SameSiteStrictMode,

		ContextName: "session",
	}
}

// WithDir returns a new Config with Dir set to the given value.
//
// The default value of Dir is set to `/data`.
func (c Config) WithDir(path string) Config {
	c.Dir = path
	return c
}

// WithInMemory returns a new Config with InMemory set to the given value.
//
// The default value of InMemory is set to false.
func (c Config) WithInMemory(val bool) Config {
	c.InMemory = val
	return c
}

// WithSessionLength returns a new Config with SessionLength set to the given duration.
//
// The default length of a session is set to 1 hour.
func (c Config) WithSessionLength(length time.Duration) Config {
	c.SessionLength = length
	return c
}

// WithExtendSessions returns a new Config with ExtendSessions set to the given value.
//
// The default value of ExtendSessions is true.
func (c Config) WithExtendSessions(val bool) Config {
	c.ExtendSessions = val
	return c
}

// WithCookieName returns a new Config with CookieName set to the given value.
//
// The default value of CookieName is session.
func (c Config) WithCookieName(name string) Config {
	c.CookieName = name
	return c
}

// WithCookiePath returns a new Config with CookiePath set to the given value.
//
// The default value of CookiePath is /.
func (c Config) WithCookiePath(path string) Config {
	c.CookiePath = path
	return c
}

// WithCookieHttpOnly returns a new Config with CookieHttpOnly set to the given value.
//
// The default value of CookieHttpOnly is true.
func (c Config) WithCookieHttpOnly(val bool) Config {
	c.CookieHttpOnly = val
	return c
}

// WithCookieSecure returns a new Config with CookieSecure set to the given value.
//
// The default value of CookieSecure is true.
func (c Config) WithCookieSecure(val bool) Config {
	c.CookieSecure = val
	return c
}

// WithCookieSameSite returns a new Config with CookieSameSite set to the given value.
//
// The default value of CookieSameSite is Strict.
func (c Config) WithCookieSameSite(val http.SameSite) Config {
	c.CookieSameSite = val
	return c
}

// WithContextName returns a new Config with ContextName set to the given value.
//
// The default value of ContextName is session.
func (c Config) WithContextName(val http.SameSite) Config {
	c.CookieSameSite = val
	return c
}

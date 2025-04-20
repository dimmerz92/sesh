package sesh

import "time"

type Config struct {
	// the directory in which the persistent session store is located.
	Dir string
	// if true, the session store will only be in memory, not persistent.
	InMemory bool

	// the duration of a single session.
	SessionLength time.Duration
	// automatically renew sessions on successful validation.
	ExtendSessions bool
}

// DefaultConfig returns a session store config with sensible defaults.
// Modify these to suit using the `WithX` chain methods.
func DefaultConfig() Config {
	return Config{
		Dir: "./session_data",

		SessionLength:  time.Hour,
		ExtendSessions: true,
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

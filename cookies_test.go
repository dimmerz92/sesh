package sesh

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var store *sessionStore
var cookie *http.Cookie

func TestNewWithCookie(t *testing.T) {
	s, err := NewSessionStore(DefaultConfig())
	if err != nil {
		t.Fatalf("failed new session store: %v", err)
	}
	store = s

	rr := httptest.NewRecorder()

	sessionId, err := store.NewWithCookie(rr, "hello")
	sessionCookie := rr.Result().Cookies()
	if len(sessionCookie) != 1 {
		t.Fatal("failed 1 new cookie")
	} else if sessionCookie[0].Name != "session" || sessionCookie[0].Value != sessionId {
		t.Fatal("failed new cookie values")
	}

	cookie = sessionCookie[0]
}

func TestGetWithCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cookie)

	rr := httptest.NewRecorder()

	var v string
	err := store.GetWithCookie(rr, r, &v)
	if err != nil {
		t.Fatalf("failed to get session from cookie: %v", err)
	}

	if v != "hello" {
		t.Fatalf("failed get session value from cookie, wanted: hello, got: %s", v)
	}
}

func TestDeleteWithCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cookie)

	rr := httptest.NewRecorder()

	err := store.DeleteWithCookie(rr, r)
	if err != nil {
		t.Fatalf("failed to delete with cookie: %v", err)
	}

	sessionCookie := rr.Result().Cookies()
	if len(sessionCookie) != 1 {
		t.Fatalf("failed delete with cookie, no invalidated cookie")
	} else if sessionCookie[0].Value != "" || sessionCookie[0].MaxAge != -1 {
		t.Fatalf("failed to delete with cookie, value and max age not set correctly")
	}

	var v string
	err = store.GetWithCookie(rr, r, &v)
	if err == nil || v != "" {
		t.Fatalf("failed to delete with cookie, value still exists: %v", err)
	}

	store.Close()
}

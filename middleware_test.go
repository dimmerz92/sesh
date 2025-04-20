package sesh

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionChecker(t *testing.T) {
	store, err := NewSessionStore(DefaultConfig())
	if err != nil {
		t.Fatalf("failed new session store: %v", err)
	}

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("failed to make new request: %v", err)
	}

	rr := httptest.NewRecorder()
	sessionId, err := store.NewWithCookie(rr, "hello")
	if err != nil {
		t.Fatalf("failed to create new session with cookie: %v", err)
	}

	sessionCookie := rr.Result().Cookies()
	if len(sessionCookie) != 1 {
		t.Fatal("failed 1 new cookie")
	} else if sessionCookie[0].Name != "session" || sessionCookie[0].Value != sessionId {
		t.Fatal("failed new cookie values")
	}

	mwNoCookie := SessionChecker[string](
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(store.config.ContextName) != nil {
				t.Fatalf("failed session checker, context should be empty")
			}
		}),
		store,
	)

	mwNoCookie.ServeHTTP(rr, r)

	mwWithCookie := SessionChecker[string](
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v, ok := r.Context().Value(store.config.ContextName).(string)
			if !ok {
				t.Fatal("failed session checker, context empty")
			} else if v != "hello" {
				t.Fatalf("failed session checker, context value not the same, wanted: hello, got: %s", v)
			}
		}),
		store,
	)

	r.AddCookie(sessionCookie[0])

	mwWithCookie.ServeHTTP(rr, r)

	store.Close()
}

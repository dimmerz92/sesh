package sesh

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type SessionStore struct {
	config Config
	db     *badger.DB
}

// NewSessionStore returns a session store database connection.
func NewSessionStore(config Config) (*SessionStore, error) {
	var db *badger.DB
	var err error

	// create or open the badgerDB session store
	if config.InMemory {
		db, err = badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	} else {
		db, err = badger.Open(badger.DefaultOptions(config.Dir).WithLogger(nil))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create a new session store: %w", err)
	}

	return &SessionStore{
		config: config,
		db:     db,
	}, nil
}

// closes the connection to the session store database.
func (s *SessionStore) Close() error {
	return s.db.Close()
}

// New adds a new session to the store along with the given data and returns a session ID.
func (s *SessionStore) New(data any) (string, error) {
	// generate a uuid session ID
	sessionId := uuid.NewString()

	// encode the given data
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return "", fmt.Errorf("failed to encode data: %w", err)
	}

	// add the session to the database
	err = s.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(sessionId), buf.Bytes()).WithTTL(s.config.SessionLength)
		return txn.SetEntry(entry)
	})
	if err != nil {
		return "", fmt.Errorf("failed to add session to store: %w", err)
	}

	return sessionId, nil
}

// Get retrieves the given session ID and stores the session data in the value pointer v.
func (s *SessionStore) Get(sessionId string, v any) error {
	var session []byte
	// retrieve the session and update ttl if it exists
	err := s.db.View(func(txn *badger.Txn) error {
		// get the session
		item, err := txn.Get([]byte(sessionId))
		if err == badger.ErrKeyNotFound {
			return fmt.Errorf("session not found: %w", err)
		} else if err != nil {
			return fmt.Errorf("failed to retrieve session: %w", err)
		}

		// get the session data
		err = item.Value(func(val []byte) error {
			session = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to get the session data: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// update the ttl if configured.
	if s.config.ExtendSessions {
		err = s.db.Update(func(txn *badger.Txn) error {
			// update the ttl
			entry := badger.NewEntry([]byte(sessionId), session).WithTTL(s.config.SessionLength)
			err = txn.SetEntry(entry)
			if err != nil {
				return fmt.Errorf("failed to update session ttl: %w", err)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	// decode the session data
	err = gob.NewDecoder(bytes.NewReader(session)).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode session data: %w", err)
	}

	return nil
}

// Delete removes a session from the session store if it exists.
func (s *SessionStore) Delete(sessionId string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(sessionId))
	})
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

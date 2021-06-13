package sessions

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/vestlog/nix/pkg/models"
)

type SessionStore struct {
	Store sessions.Store
	Name  string
}

func (s *SessionStore) CreateSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		if err := session.Save(r, w); err != nil {
			return fmt.Errorf("error: could not save session in store: %w", err)
		}
	}
	return nil
}

func (s *SessionStore) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		return fmt.Errorf("error: could not get session: %w", err)
	}
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error: could not save session: %w", err)
	}
	return nil
}

func (s *SessionStore) SaveData(w http.ResponseWriter, r *http.Request, name string, value interface{}) error {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		return fmt.Errorf("error: could not get session from store: %w", err)
	}
	session.Values[name] = value
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error: could not save session to store: %w", err)
	}
	return nil
}

func (s *SessionStore) GetData(r *http.Request, name string) (interface{}, error) {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		return "", fmt.Errorf("error: could not get session from store: %w", err)
	}
	return session.Values[name], nil
}

func (s *SessionStore) SaveString(w http.ResponseWriter, r *http.Request, name string, value string) error {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		return fmt.Errorf("error: could not get session from store: %w", err)
	}
	session.Values[name] = value
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error: could not save session to store: %w", err)
	}
	return nil
}

func (s *SessionStore) GetString(r *http.Request, name string) (string, error) {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		return "", fmt.Errorf("error: could not get session from store: %w", err)
	}
	value, ok := session.Values[name].(string)
	if !ok {
		return "", fmt.Errorf("error: value is not a string")
	}
	return value, nil
}

func CreateCookieStore(name string, keyPairs ...[]byte) *SessionStore {
	gob.Register(&models.User{})
	return &SessionStore{
		Store: sessions.NewCookieStore(keyPairs...),
		Name:  name,
	}
}

// ?
func GetSessionError(err error) error {
	return fmt.Errorf("error: could not get session from store: %w", err)
}

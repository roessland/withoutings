package sessions

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/roessland/withoutings/withingsapi"
	"net/http"
)

func init() {
	gob.Register(withingsapi.Token{})
}

type Manager struct {
	*sessions.CookieStore
	Name string
}

func NewManager(secret []byte) *Manager {
	return &Manager{
		CookieStore: sessions.NewCookieStore(secret),
		Name:        "_sess",
	}
}

func (m *Manager) New(r *http.Request) (*Session, error) {
	sess, err := m.CookieStore.New(r, m.Name)
	return &Session{sess}, err
}

// Get returns a new session and an error if the session exists but could
// not be decoded.
func (m *Manager) Get(r *http.Request) (*Session, error) {
	sess, err := m.CookieStore.Get(r, m.Name)
	return &Session{sess}, err
}

type Session struct {
	*sessions.Session
}

func (sess *Session) SetState(nonce string) {
	sess.Values["state"] = nonce
}

func (sess *Session) State() string {
	nonce, ok := sess.Values["state"].(string)
	if !ok {
		return ""
	}
	return nonce
}

func (sess *Session) SetToken(token *withingsapi.Token) {
	if token != nil {
		sess.Values["token"] = *token
	}
}

func (sess *Session) Token() *withingsapi.Token {
	token, ok := sess.Values["token"].(withingsapi.Token)
	if !ok {
		return nil
	}
	return &token
}

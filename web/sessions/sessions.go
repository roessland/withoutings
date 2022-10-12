package sessions

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/roessland/withoutings/withings"
	"net/http"
)

func init() {
	gob.Register(withings.Token{})
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

func (m *Manager) Get(r *http.Request) (*Session, error) {
	sess, err := m.CookieStore.Get(r, m.Name)
	if err != nil {
		return nil, err
	}
	return &Session{sess}, nil
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

func (sess *Session) SetToken(token *withings.Token) {
	if token != nil {
		sess.Values["token"] = *token
	}
}

func (sess *Session) Token() *withings.Token {
	token, ok := sess.Values["token"].(withings.Token)
	if !ok {
		return nil
	}
	return &token
}

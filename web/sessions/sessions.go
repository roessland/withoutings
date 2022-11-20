package sessions

import (
	"context"
	"encoding/base32"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgtype"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/withingsapi"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	gob.Register(withingsapi.Token{})
}

type Manager struct {
	sessions.Store
	Name string
}

func NewManager(secret []byte) *Manager {
	return &Manager{
		Store: sessions.NewCookieStore(secret),
		Name:  "_sess",
	}
}

func (m *Manager) New(r *http.Request) (*Session, error) {
	sess, err := m.Store.New(r, m.Name)
	return &Session{sess}, err
}

// Get returns a new session and an error if the session exists but could
// not be decoded.
func (m *Manager) Get(r *http.Request) (*Session, error) {
	sess, err := m.Store.Get(r, m.Name)
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

func (sess *Session) SetAccountID(accountID int64) {
	sess.Values["account_id"] = accountID
}

func (sess *Session) AccountID() int64 {
	accountID, ok := sess.Values["account_id"].(int64)
	if !ok {
		return -1
	}
	return accountID
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewDatabaseStore(path string, keyPairs ...[]byte) *DatabaseStore {
	if path == "" {
		path = os.TempDir()
	}
	fs := &DatabaseStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		path: path,
	}

	fs.MaxAge(fs.Options.MaxAge)
	return fs
}

// DatabaseStore stores sessions in the database.
type DatabaseStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options // default configuration
	path    string
	queries *db.Queries
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new DatabaseStore is 4096.
func (s *DatabaseStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *DatabaseStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *DatabaseStore) New(r *http.Request, name string) (*sessions.Session, error) {
	ctx := r.Context()
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		session.ID = c.Value
		err = s.load(ctx, session)
		if err == nil {
			session.IsNew = false
		}
	}
	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store path. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *DatabaseStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	ctx := r.Context()
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.erase(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(ctx, session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *DatabaseStore) MaxAge(age int) {
	s.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// save writes encoded session.Values to database.
func (s *DatabaseStore) save(ctx context.Context, session *sessions.Session) error {
	var data pgtype.JSONB
	err := data.Set(session.Values)
	if err != nil {
		return err
	}
	sessionUuid, err := s.queries.CreateSession(ctx, data)
	session.ID = sessionUuid.String()
	if err != nil {
		return err
	}
	return nil
}

// load reads a session from DB and decodes its content into session.Values.
func (s *DatabaseStore) load(ctx context.Context, session *sessions.Session) error {
	sessionUuid, err := uuid.Parse(session.ID)
	if err != nil {
		return err
	}
	dbsession, err := s.queries.GetSession(ctx, sessionUuid)
	if err != nil {
		return err
	}
	return dbsession.Data.Scan(&session.Values)
}

// delete session from database
func (s *DatabaseStore) erase(session *sessions.Session) error {
	filename := filepath.Join(s.path, "session_"+session.ID)

	err := os.Remove(filename)
	return err
}

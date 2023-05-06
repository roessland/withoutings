package flash

import (
	"context"
	"github.com/alexedwards/scs/v2"
)

type contextKey struct{}

var contextKeyFlash contextKey = struct{}{}

func GetMsgFromContext(ctx context.Context) string {
	msg, ok := ctx.Value(contextKeyFlash).(string)
	if !ok {
		return ""
	}
	return msg
}

func AddMsgToContext(ctx context.Context, msg string) context.Context {
	return context.WithValue(ctx, contextKeyFlash, msg)
}

type Manager struct {
	Sessions *scs.SessionManager
}

func NewManager(sessions *scs.SessionManager) *Manager {
	if sessions == nil {
		panic("sessions is nil")
	}
	return &Manager{
		Sessions: sessions,
	}
}

func (m *Manager) PutMsg(ctx context.Context, msg string) {
	m.Sessions.Put(ctx, "flash", msg)
}

func (m *Manager) PopMsg(ctx context.Context) string {
	return m.Sessions.PopString(ctx, "flash")
}

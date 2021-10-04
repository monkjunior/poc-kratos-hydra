package context

import (
	"context"
)

const (
	sessionKey         privateKey = "is_session_active"
	logoutKey          privateKey = "log_out_url"
)

type privateKey string

// SetSession is used to set the session value to context
func SetSession(ctx context.Context, session bool, logoutURL string) context.Context {
	ctx = context.WithValue(ctx, sessionKey, session)
	ctx = context.WithValue(ctx, logoutKey, logoutURL)
	return ctx
}

// GetSession is used to get the kratos session value from context
func GetSession(ctx context.Context) (bool, string) {
	s := ctx.Value(sessionKey)
	l := ctx.Value(logoutKey)
	if s == nil || l == nil {
		return false, ""
	}

	if session, okS := s.(bool); okS {
		if logoutURL, okL := l.(string); okL {
			return session, logoutURL
		}
	}
	return false, ""
}

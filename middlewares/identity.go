package middlewares

import (
	builtInCtx "context"
	"net/http"

	"github.com/monkjunior/poc-kratos-hydra/context"
	kratosClient "github.com/ory/kratos-client-go"
)

// Identity middleware checks if current session of received request is active, it then saves the result and logoutURL
// to the request context.
type Identity struct {
	KratosClient *kratosClient.APIClient
}

// Apply logs request before passing it to http.Handler
func (mw *Identity) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will try to set to Session Active value and appropriate LogoutURL
// to request context before passing it to the HandleFunc.
func (mw *Identity) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ory_kratos_session")
		if err != nil || cookie == nil {
			next(w, r)
			return
		}

		session, res, err := mw.KratosClient.V0alpha1Api.ToSession(builtInCtx.Background()).Cookie(r.Header.Get("Cookie")).Execute()
		if err != nil || res == nil || res.StatusCode != http.StatusOK {
			next(w, r)
			return
		}

		if !session.GetActive() {
			next(w, r)
			return
		}

		logoutURL, res, err := mw.KratosClient.V0alpha1Api.CreateSelfServiceLogoutFlowUrlForBrowsers(builtInCtx.Background()).Cookie(r.Header.Get("Cookie")).Execute()
		if err != nil || res == nil || res.StatusCode != http.StatusOK {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.SetSession(ctx, session.GetActive(), logoutURL.GetLogoutUrl())
		r = r.WithContext(ctx)
		next(w, r)
		return
	}
}

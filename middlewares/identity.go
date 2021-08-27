package middlewares

import (
	builtInCtx "context"
	"github.com/monkjunior/poc-kratos-hydra/context"
	"log"
	"net/http"

	kratosClient "github.com/ory/kratos-client-go"
)

type Identity struct {
	KratosClient *kratosClient.APIClient
}

func (mw *Identity) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will try to set to Session Active value and appropriate LogoutURL
// to request context before passing it to the HandleFunc.
func (mw *Identity) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ory_kratos_session")
		if err != nil || cookie == nil {
			log.Println("cookie is nil")
			next(w, r)
			return
		}
		log.Println("ory_kratos_session: ", cookie.Value)

		session, res, err := mw.KratosClient.V0alpha1Api.ToSession(builtInCtx.Background()).Cookie(r.Header.Get("Cookie")).Execute()
		if err != nil || res == nil || res.StatusCode != http.StatusOK {
			log.Println("can not got is_session_active", err)
			next(w, r)
			return
		}

		if !session.GetActive() {
			log.Println("is_session_active: ", false)
			next(w, r)
			return
		}

		logoutURL, res, err := mw.KratosClient.V0alpha1Api.CreateSelfServiceLogoutFlowUrlForBrowsers(builtInCtx.Background()).Cookie(r.Header.Get("Cookie")).Execute()
		if err != nil || res == nil || res.StatusCode != http.StatusOK {
			log.Println("can not got logoutURL", err)
			next(w, r)
			return
		}
		log.Println("session is active and logoutURL is ", logoutURL.GetLogoutUrl())
		ctx := r.Context()
		ctx = context.SetSession(ctx, session.GetActive(), logoutURL.GetLogoutUrl())
		r = r.WithContext(ctx)
		next(w, r)
		return
	}
}

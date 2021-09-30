package controllers

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/monkjunior/poc-kratos-hydra/pkg/rand"
	"github.com/monkjunior/poc-kratos-hydra/pkg/views"
	hydraSDK "github.com/ory/hydra-client-go/client"
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	hydraModel "github.com/ory/hydra-client-go/models"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	hydraLoginState string
)

func NewHydra(k *kratosClient.APIClient, hCli *hydraSDK.OryHydra, hAdm *hydraSDK.OryHydra) *Hydra {
	return &Hydra{
		ConsentView:  views.NewView("bootstrap", "consent"),
		kratosClient: k,
		hydraClient:  hCli,
		hydraAdmin:   hAdm,
	}
}

// Hydra controller will handler flows relate to Hydra integration: login with Hydra flow, and so on
// It interacts with Ory Kratos, an opensource Identity Provider, and Ory Hydra, an opensource OAuth2/OIDC provider.
type Hydra struct {
	ConsentView  *views.View
	kratosClient *kratosClient.APIClient
	hydraClient  *hydraSDK.OryHydra
	hydraAdmin   *hydraSDK.OryHydra
}

// GetHydraLogin
// GET /auth/hydra/login
func (h *Hydra) GetHydraLogin(w http.ResponseWriter, r *http.Request) {
	loginChallenge := r.URL.Query().Get("login_challenge")
	if loginChallenge == "" {
		log.Println("Missing login_challenge parameter")
		redirectToLogin(w, r)
		return
	}
	params := &hydraAdmin.GetLoginRequestParams{
		LoginChallenge: loginChallenge,
		Context:        r.Context(),
	}
	isOK, err := h.hydraAdmin.Admin.GetLoginRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra login info with login_challenge =", loginChallenge, err)
		redirectToLogin(w, r)
		return
	}
	payload := isOK.GetPayload()

	if *payload.Skip {
		// We can do some logic here, for example
		// update the number of times the user logged in.
		// We can also deny if there is something went wrong.
		h.acceptLogin(w, r, loginChallenge)
		return
	}

	state := r.URL.Query().Get("hydra_login_state")
	if state == "" {
		log.Println("Got empty hydra login state")
		redirectToLogin(w, r)
		return
	}

	kratosSessionCookie, err := r.Cookie("ory_kratos_session")
	if err != nil {
		log.Println("Failed to get ory_kratos_session", err)
		redirectToLogin(w, r)
		return
	}
	if kratosSessionCookie.Value == "" {
		log.Println("No kratos login session was set")
		redirectToLogin(w, r)
		return
	}
	// TODO: Need to enhance the way we validate this param to prevent conflicts
	if state != hydraLoginState {
		log.Println("Mismatch hydra login state")
		redirectToLogin(w, r)
		return
	}

	h.acceptLogin(w, r, loginChallenge)
}

// ConsentForm stores consent form data to render consent page
type ConsentForm struct {
	// TODO: implement csrf protection using gorilla csrf
	Subject                 string
	ConsentChallenge        string   `schema:"consent_challenge"`
	Scopes                  []string `schema:"scopes"`
	Remember                bool     `schema:"remember"`
	Accept                  string   `schema:"accept"`
	AccessTokenCustomClaims AccessTokenCustomClaims
}

// AccessTokenCustomClaims defines and stores some data from Kratos session.
// We use it to enrich information when perform token introspection.
// Then, these information will be set to HTTP header by using Oathkeeper mutator.
type AccessTokenCustomClaims struct {
	UserUUID string
	Email    string
}

// GetHydraConsent
// GET /auth/hydra/consent
func (h *Hydra) GetHydraConsent(w http.ResponseWriter, r *http.Request) {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		log.Println("Missing consent_challenge parameter")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	params := hydraAdmin.NewGetConsentRequestParams()
	params.ConsentChallenge = consentChallenge
	isOK, err := h.hydraAdmin.Admin.GetConsentRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra consent info with consent_challenge =", consentChallenge, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	payload := isOK.GetPayload()
	form := &ConsentForm{
		Subject:          payload.Subject,
		ConsentChallenge: consentChallenge,
		Scopes:           strings.Split(payload.Client.Scope, " "),
	}
	form.AccessTokenCustomClaims = h.getAccessTokenCustomClaims(r)
	if payload.Skip {
		h.acceptConsent(w, r, *form)
		return
	}

	// intercept consent if login from first party app
	if payload.Client.ClientID == config.Cfg.Hydra.Client.ID {
		h.acceptConsent(w, r, *form)
		return
	}

	data := views.Data{
		Yield: form,
	}
	h.ConsentView.Render(w, r, data)
}

// PostHydraConsent
// POST /auth/hydra/consent
func (h *Hydra) PostHydraConsent(w http.ResponseWriter, r *http.Request) {
	var form ConsentForm
	if err := parseForm(r, &form); err != nil {
		log.Println("Could not parse consent form", err)
		http.Error(w, "Could not parse consent form", http.StatusInternalServerError)
		return
	}
	log.Printf("Consent form: %+v\n", form)
	if form.Accept == "deny" {
		h.denyConsent(w, r, form)
		return
	}
	h.acceptConsent(w, r, form)
}

// getAccessTokenCustomClaims enriches custom claims info by requesting
// to Kratos "whoami" endpoint for more information about the identity.
func (h *Hydra) getAccessTokenCustomClaims(r *http.Request) AccessTokenCustomClaims {
	session, res, err := h.kratosClient.V0alpha1Api.ToSession(r.Context()).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK || !session.GetActive() {
		return AccessTokenCustomClaims{}
	}
	userUUID := session.Identity.GetId()
	var userEmail string
	if identityTraits, ok := session.Identity.Traits.(map[string]interface{}); ok {
		userEmail, _ = identityTraits["email"].(string)
	}
	return AccessTokenCustomClaims{
		UserUUID: userUUID,
		Email:    userEmail,
	}

}

// acceptLogin will redirect to return endpoint if the process is successful
// or generate an error page if an error occurred
func (h *Hydra) acceptLogin(w http.ResponseWriter, r *http.Request, loginChallenge string) {
	session, res, err := h.kratosClient.V0alpha1Api.ToSession(r.Context()).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK || !session.GetActive() {
		log.Println("You did not log in")
		redirectToLogin(w, r)
		return
	}

	identityID := session.Identity.GetId()
	identityTraits := session.Identity.Traits
	sessionID := session.GetId()
	isSessionActive := session.GetActive()

	log.Printf(`Info of logged in user
UserID: %v
SessionID: %v
IsActive: %v
UserInfo %v
`, identityID, sessionID, isSessionActive, identityTraits)

	loginReqParams := &hydraAdmin.AcceptLoginRequestParams{}
	loginReqParams.WithLoginChallenge(loginChallenge)
	loginReqParams.WithBody(&hydraModel.AcceptLoginRequest{
		Subject:     &identityID,
		Remember:    true,
		RememberFor: 3600,
	})
	loginReqParams.WithContext(r.Context())
	acceptRes, err := h.hydraAdmin.Admin.AcceptLoginRequest(loginReqParams)
	if err != nil {
		log.Println("Failed to accept hydra login request", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *acceptRes.Payload.RedirectTo, http.StatusFound)
}

// acceptConsent fetches consent info of current consent challenge, uses both consent info
// and post form form the HTTP request to accept the consent challenge.
func (h *Hydra) acceptConsent(w http.ResponseWriter, r *http.Request, form ConsentForm) {
	params := hydraAdmin.NewGetConsentRequestParams()
	params.ConsentChallenge = form.ConsentChallenge
	isOK, err := h.hydraAdmin.Admin.GetConsentRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra consent info with consent_challenge =", form.ConsentChallenge, err)
		http.Error(w, "Some thing went wrong", http.StatusInternalServerError)
		return
	}
	payload := isOK.GetPayload()

	consentParams := &hydraAdmin.AcceptConsentRequestParams{
		ConsentChallenge: form.ConsentChallenge,
		Body: &hydraModel.AcceptConsentRequest{
			GrantScope:               form.Scopes,
			GrantAccessTokenAudience: payload.RequestedAccessTokenAudience,
			Remember:                 form.Remember,
			RememberFor:              3600,
			Session: &hydraModel.ConsentRequestSession{
				AccessToken: form.AccessTokenCustomClaims,
				IDToken:     form.AccessTokenCustomClaims,
			},
		},
		Context: r.Context(),
	}
	consentOK, err := h.hydraAdmin.Admin.AcceptConsentRequest(consentParams)
	if err != nil {
		log.Println("Could not accept consent challenge ", err)
		http.Error(w, "Some thing went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *consentOK.Payload.RedirectTo, http.StatusFound)
}

// denyConsent rejects the consent request, it usually happens when user presses DENY button
func (h *Hydra) denyConsent(w http.ResponseWriter, r *http.Request, form ConsentForm) {
	consentParams := &hydraAdmin.RejectConsentRequestParams{
		Context:          r.Context(),
		ConsentChallenge: form.ConsentChallenge,
		Body: &hydraModel.RejectRequest{
			Error:            "User denied access",
			ErrorDescription: "Put some description about the error later!",
			ErrorHint:        "Error hint: ...",
			ErrorDebug:       "Error debug: ...",
			StatusCode:       0,
		},
	}
	rejectOK, err := h.hydraAdmin.Admin.RejectConsentRequest(consentParams)
	if err != nil {
		log.Println("Could not reject consent request ", err)
		http.Error(w, "Some thing went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *rejectOK.Payload.RedirectTo, http.StatusFound)
}

// redirectToLogin redirect to login endpoint to perform login
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	state, err := rand.GenerateHydraState()
	hydraLoginState = state
	if err != nil {
		log.Println("Failed to generate hydra state", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	v := url.Values{}
	v.Add("login_challenge", r.URL.Query().Get("login_challenge"))
	v.Add("hydra_login_state", hydraLoginState)
	returnToString := config.Cfg.BaseURL + "/auth/hydra/login?" + url.QueryEscape(v.Encode())
	redirectUrl := KratosPublicURL + KratosSSLoginBrowserPath + "?refresh=true&return_to=" + returnToString
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

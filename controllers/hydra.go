package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/monkjunior/poc-kratos-hydra/rand"
	"github.com/monkjunior/poc-kratos-hydra/views"
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
	params := hydraAdmin.NewGetLoginRequestParams()
	params.LoginChallenge = loginChallenge
	isOK, err := h.hydraAdmin.Admin.GetLoginRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra login info with login_challenge =", loginChallenge, err)
		redirectToLogin(w, r)
		return
	}
	payload := isOK.GetPayload()

	// TODO: need to handle this case
	// skip is true often happens when your session is still valid after a previous succeed login challenge
	if *payload.Skip {
		fmt.Fprintln(w, "Skip is true, we should accept this login request from Hydra", http.StatusOK)
		return
	}

	state := r.URL.Query().Get("hydra_login_state")
	log.Println("hydra_login_state=", state)
	if state == "" {
		log.Println("Got empty hydra login state")
		redirectToLogin(w, r)
		return
	}

	kratosSessionCookie, err := r.Cookie("ory_kratos_session")
	log.Println("ory_kratos_session=", kratosSessionCookie)
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

	loginReqBody := &hydraModel.AcceptLoginRequest{
		Subject:     &identityID,
		Remember:    true,
		RememberFor: 3600,
	}
	loginReqParams := &hydraAdmin.AcceptLoginRequestParams{}
	loginReqParams.WithLoginChallenge(loginChallenge)
	loginReqParams.WithBody(loginReqBody)
	loginReqParams.WithContext(r.Context())
	acceptRes, err := h.hydraAdmin.Admin.AcceptLoginRequest(loginReqParams)
	if err != nil {
		log.Println("Failed to accept hydra login request", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *acceptRes.Payload.RedirectTo, http.StatusFound)
}

type ConsentForm struct {
	Subject          string
	ConsentChallenge string   `schema:"consent_challenge"`
	Scopes           []string `schema:"scopes"`
	Remember         string   `schema:"remember"`
	Accept           string   `schema:"accept"`
}

// GetHydraConsent
// GET /auth/hydra/consent
func (h *Hydra) GetHydraConsent(w http.ResponseWriter, r *http.Request) {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		fmt.Fprintln(w, "Missing consent_challenge parameter")
		return
	}

	params := hydraAdmin.NewGetConsentRequestParams()
	params.ConsentChallenge = consentChallenge
	isOK, err := h.hydraAdmin.Admin.GetConsentRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra consent info with consent_challenge =", consentChallenge, err)
		fmt.Fprintln(w, "Failed to fetch consent info")
		return
	}
	payload := isOK.GetPayload()
	if payload.Skip {
		fmt.Fprintln(w, "Skip is true, we should accept this consent request", http.StatusOK)
		return
	}

	data := views.Data{
		Yield: ConsentForm{
			Subject:          payload.Subject,
			ConsentChallenge: consentChallenge,
			Scopes:           strings.Split(payload.Client.Scope, " "),
		},
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
	fmt.Fprintf(w, "Consent form: %+v", form)
	w.WriteHeader(http.StatusOK)
}

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
	returnToString := "http://127.0.0.1:4455/auth/hydra/login?" + url.QueryEscape(v.Encode())
	redirectUrl := KratosPublicBaseURL + "/self-service/login/browser?refresh=true&return_to=" + returnToString
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

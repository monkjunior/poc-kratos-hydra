package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/monkjunior/poc-kratos-hydra/common"
	"github.com/monkjunior/poc-kratos-hydra/rand"
	"github.com/monkjunior/poc-kratos-hydra/views"
	hydraSDK "github.com/ory/hydra-client-go/client"
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	hydraModel "github.com/ory/hydra-client-go/models"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	// TODO: this should be readable from config
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"

	hydraLoginState string
)

func NewUsers(k *kratosClient.APIClient, hCli *hydraSDK.OryHydra, hAdm *hydraSDK.OryHydra) *Users {
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     k,
		hydraClient:      hCli,
		hydraAdmin:       hAdm,
	}
}

type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	kratosClient     *kratosClient.APIClient
	hydraClient      *hydraSDK.OryHydra
	hydraAdmin       *hydraSDK.OryHydra
}

// LoginForm stores data for rendering Login form and submit a Login flow
type LoginForm struct {
	SubmitMethod string
	Action       string
	CsrfToken    string `schema:"csrf_token"`
	FlowID       string
	Email        string `schema:"password_identifier"`
	Password     string `schema:"password"`
}

// GetLogin requires flow params, if the flow is not set, it will redirect to Kratos to browse a new one.
// Kratos will create a new flow and redirect back to /auth/login with the param was set in the URL.
// GetLogin will use this id to fetch data from Kratos to render submit form.
//
// GET /auth/login/?flow=<flow_id>
func (u *Users) GetLogin(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/login/browser", http.StatusFound)
		return
	}
	flowObject, res, err := u.kratosClient.V0alpha1Api.GetSelfServiceLoginFlow(r.Context()).Id(flow).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		common.LogOnError(err, res)
		return
	}
	common.PrintJSONPretty(flowObject)
	data := views.Data{
		Yield: RegistrationForm{
			CsrfToken:    flowObject.Ui.GetNodes()[0].Attributes.UiNodeInputAttributes.Value.(string),
			FlowID:       flow,
			SubmitMethod: flowObject.Ui.Method,
			Action:       flowObject.Ui.Action,
		},
	}
	u.LoginView.Render(w, r, data)
}

// GetHydraLogin
// GET /auth/hydra/login
func (u *Users) GetHydraLogin(w http.ResponseWriter, r *http.Request) {
	loginChallenge := r.URL.Query().Get("login_challenge")
	if loginChallenge == "" {
		log.Println("Missing login_challenge parameter")
		redirectToLogin(w, r)
		return
	}
	params := hydraAdmin.NewGetLoginRequestParams()
	params.LoginChallenge = loginChallenge
	isOK, err := u.hydraAdmin.Admin.GetLoginRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra login info with login_challenge =", loginChallenge, err)
		redirectToLogin(w, r)
		return
	}
	payload := isOK.GetPayload()
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

	session, res, err := u.kratosClient.V0alpha1Api.ToSession(r.Context()).Cookie(r.Header.Get("Cookie")).Execute()
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
	acceptRes, err := u.hydraAdmin.Admin.AcceptLoginRequest(loginReqParams)
	if err != nil {
		log.Println("Failed to accept hydra login request", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *acceptRes.Payload.RedirectTo, http.StatusFound)
}

// GetHydraConsent
// GET /auth/hydra/consent
func (u *Users) GetHydraConsent(w http.ResponseWriter, r *http.Request) {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		fmt.Fprintln(w, "Missing consent_challenge parameter")
		return
	}
	fmt.Fprintln(w, "consent_challenge =", consentChallenge)

}

// RegistrationForm stores data for rendering Registration form and submit a Registration flow
type RegistrationForm struct {
	RegistrationMethod string
	SubmitMethod       string
	Action             string
	CsrfToken          string `schema:"csrf_token"`
	FlowID             string
	Email              string `schema:"traits.email"`
	Password           string `schema:"password"`
}

// GetRegistration requires flow params to render Registration screen
// if flow param is not found, it will redirect to Kratos /self-service/registration/browser
// to browse a new flow_id.
// Kratos then redirect back to this path with a flow param in the URL.
//
// GET /auth/registration/?flow=<flow_id>
func (u *Users) GetRegistration(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/registration/browser", http.StatusFound)
		return
	}
	flowObject, res, err := u.kratosClient.V0alpha1Api.GetSelfServiceRegistrationFlow(r.Context()).Id(flow).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		common.LogOnError(err, res)
		return
	}
	//kratos.PrintJSONPretty(flowObject)
	data := views.Data{
		Yield: RegistrationForm{
			CsrfToken:    flowObject.Ui.GetNodes()[0].Attributes.UiNodeInputAttributes.Value.(string),
			FlowID:       flow,
			SubmitMethod: flowObject.Ui.Method,
			Action:       flowObject.Ui.Action,
		},
	}
	u.RegistrationView.Render(w, r, data)
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

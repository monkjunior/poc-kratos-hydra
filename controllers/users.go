package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/monkjunior/poc-kratos-hydra/context"
	"github.com/monkjunior/poc-kratos-hydra/kratos"
	"github.com/monkjunior/poc-kratos-hydra/rand"
	"github.com/monkjunior/poc-kratos-hydra/views"
	hydraSDK "github.com/ory/hydra-client-go/client"
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	// TODO: this should be moved to kratos package
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"
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
		kratos.LogOnError(err, res)
		return
	}
	kratos.PrintJSONPretty(flowObject)
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
		http.Error(w, "Missing login_challenge parameter", http.StatusForbidden)
		return
	}
	params := hydraAdmin.NewGetLoginRequestParams()
	params.LoginChallenge = loginChallenge
	isOK, err := u.hydraAdmin.Admin.GetLoginRequest(params)
	if err != nil || isOK == nil {
		log.Println("Failed to fetch hydra login info with login_challenge =", loginChallenge, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	payload := isOK.GetPayload()
	if *payload.Skip {
		fmt.Fprintln(w, "Skip is true, we should accept this login request from Hydra", http.StatusOK)
		return
	}

	hydraLoginState := r.URL.Query().Get("hydra_login_state")
	log.Println("hydra_login_state=",hydraLoginState)
	if hydraLoginState == "" {
		log.Println("Got empty hydra login state, redirect to login page")
		redirectToLogin(w, r)
		return
	}

	kratosSessionCookie, err := r.Cookie("ory_kratos_session")
	log.Println("ory_kratos_session=",kratosSessionCookie)
	if err != nil {
		log.Println("Failed to get ory_kratos_session", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if kratosSessionCookie.Value == "" {
		log.Println("No kratos login session was set, we should redirect to login page")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		redirectToLogin(w, r)
		return
	}

	// TODO: How to validate OIDC state received from the request
	log.Println("login_hint ", payload.OidcContext.LoginHint)
	if hydraLoginState != context.GetHydraLoginState(r.Context()) {
		log.Println("Mismatch hydra login state, we should redirect to login page")
		log.Println("Query param: ", hydraLoginState)
		log.Println("Value from context ", context.GetHydraLoginState(r.Context()))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Now you should figure out the user and accept login request", http.StatusOK)
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
		kratos.LogOnError(err, res)
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
	if err != nil {
		log.Println("Failed to generate hydra state", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	context.SetHydraLoginState(ctx, state)
	r = r.WithContext(ctx)

	v := url.Values{}
	v.Add("login_challenge", r.URL.Query().Get("login_challenge"))
	v.Add("hydra_login_state", state)
	returnToString := "http://127.0.0.1:4455/auth/hydra/login?" + url.QueryEscape(v.Encode())
	redirectUrl := KratosPublicBaseURL + "/self-service/login/browser?refresh=true&return_to=" + returnToString
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

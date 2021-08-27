package controllers

import (
	"net/http"

	"github.com/monkjunior/poc-kratos-hydra/kratos"
	"github.com/monkjunior/poc-kratos-hydra/views"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	// TODO: this should be moved to kratos package
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"
)

func NewUsers(k *kratosClient.APIClient) *Users {
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     k,
	}
}

type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	kratosClient     *kratosClient.APIClient
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
	kratos.PrintJSONPretty(flowObject)
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

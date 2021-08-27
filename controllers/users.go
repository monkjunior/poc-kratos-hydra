package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/monkjunior/poc-kratos-hydra/kratos"
	"github.com/monkjunior/poc-kratos-hydra/views"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	// TODO: this should be moved to kratos package
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"

	CfgKratos = kratosClient.Configuration{
		Host:   "oathkeeper:4455",
		Scheme: "http",
		Debug:  true,
		Servers: []kratosClient.ServerConfiguration{
			{
				URL: "/.ory/kratos/public",
			},
		},
	}
)

func NewUsers() *Users {
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     kratosClient.NewAPIClient(&CfgKratos),
	}
}

type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	kratosClient     *kratosClient.APIClient
}

func (u *Users) GetLogin(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/login/browser", http.StatusFound)
		return
	}
	u.LoginView.Render(w, r, nil)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
	Remember bool   `schema:"remember"`
}

func (u *Users) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var form LoginForm
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(form)
	// TODO: implement login
	//http.Redirect(w, r, KratosPublicBaseURL + "/self-service/login/browser", http.StatusFound)
}

type RegistrationForm struct {
	RegistrationMethod string
	SubmitMethod       string
	Action             string
	CsrfToken          string `schema:"csrf_token"`
	FlowID             string
	Email              string `schema:"traits.email"`
	Password           string `schema:"password"`
}

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

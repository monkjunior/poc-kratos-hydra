package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/monkjunior/poc-kratos-hydra/kratos"
	"github.com/monkjunior/poc-kratos-hydra/views"
)

var (
	// TODO: this should be moved to kratos package
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"
)

func NewUsers() *Users {
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     kratos.NewClient(),
	}
}

type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	kratosClient     kratos.ClientService
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

func (u *Users) GetRegistration(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/registration/browser", http.StatusFound)
		return
	}
	u.RegistrationView.Render(w, r, nil)
}

type RegistrationForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var form RegistrationForm
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(&form, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// TODO: implement register
	//http.Redirect(w, r, KratosPublicBaseURL+"/self-service/register/browser", http.StatusFound)
}

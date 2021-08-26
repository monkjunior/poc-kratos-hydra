package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/monkjunior/poc-kratos-hydra/views"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	KratosPublicBaseURL = "http://127.0.0.1:4455/.ory/kratos/public"
)

func NewUsers() *Users {
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     kratosClient.NewConfiguration(),
	}
}

type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	kratosClient     *kratosClient.Configuration
}

func (u *Users) GetLogin(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Println("GET /auth/login | flow not found")
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/login/browser", http.StatusFound)
		return
	}
	log.Printf("GET /auth/login | flow = %s\n", flow)
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(form)
	// TODO: implement login
	//http.Redirect(w, r, KratosPublicBaseURL + "/self-service/login/browser", http.StatusFound)
}

func (u *Users) GetRegistration(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Println("GET /auth/registration | flow not found")
		http.Redirect(w, r, KratosPublicBaseURL+"/self-service/registration/browser", http.StatusFound)
		return
	}
	log.Printf("GET /auth/registration | flow = %s\n", flow)
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
	log.Printf("POST /auth/registration | form=%+v\n", form)
	// TODO: implement register
	//http.Redirect(w, r, KratosPublicBaseURL+"/self-service/register/browser", http.StatusFound)
}

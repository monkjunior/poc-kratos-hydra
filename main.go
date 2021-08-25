package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/controllers"
)

func main() {
	protectedSites := controllers.NewProtectedSites()
	userC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", protectedSites.Dashboard)
	r.HandleFunc("/auth/login", userC.GetLogin).Methods("GET")
	r.HandleFunc("/auth/login", userC.PostLogin).Methods("POST")
	r.HandleFunc("/auth/registration", userC.GetRegistration).Methods("GET")
	r.HandleFunc("/auth/registration", userC.PostRegistration).Methods("POST")
	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", r))
}
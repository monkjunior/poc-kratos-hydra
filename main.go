package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/controllers"
	"github.com/monkjunior/poc-kratos-hydra/middlewares"
)

func main() {
	protectedSites := controllers.NewProtectedSites()
	userC := controllers.NewUsers()

	logMw := middlewares.EntryLog{}

	r := mux.NewRouter()
	r.Handle("/", logMw.Apply(protectedSites.Dashboard))
	r.HandleFunc("/auth/login", logMw.ApplyFn(userC.GetLogin)).Methods("GET")
	r.HandleFunc("/auth/login", logMw.ApplyFn(userC.PostLogin)).Methods("POST")
	r.HandleFunc("/auth/registration", logMw.ApplyFn(userC.GetRegistration)).Methods("GET")
	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", r))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/controllers"
	"github.com/monkjunior/poc-kratos-hydra/middlewares"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
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

func main() {
	k := kratosClient.NewAPIClient(&CfgKratos)

	publicSites := controllers.NewPublicSites()
	protectedSites := controllers.NewProtectedSites()
	userC := controllers.NewUsers(k)

	logMw := middlewares.EntryLog{}
	identityMw := middlewares.Identity{KratosClient: k}

	r := mux.NewRouter()

	r.Handle("/", publicSites.Home)
	r.Handle("/dashboard", protectedSites.Dashboard)

	r.HandleFunc("/auth/login", userC.GetLogin).Methods("GET")
	r.HandleFunc("/auth/registration", userC.GetRegistration).Methods("GET")

	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", logMw.Apply(identityMw.Apply(r))))
}

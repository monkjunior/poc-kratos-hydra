package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/controllers"
	"github.com/monkjunior/poc-kratos-hydra/middlewares"
	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

func main() {
	k := kratosSDK.NewAPIClient(&CfgKratos)
	hCli := hydraSDK.NewHTTPClientWithConfig(nil, &ConfigHydraClient)
	hAdm := hydraSDK.NewHTTPClientWithConfig(nil, &ConfigHydraAdmin)

	publicSites := controllers.NewPublicSites()
	protectedSites := controllers.NewProtectedSites()
	userC := controllers.NewUsers(k)
	hydraC := controllers.NewHydra(k, hCli, hAdm)

	logMw := middlewares.EntryLog{}
	identityMw := middlewares.Identity{KratosClient: k}

	r := mux.NewRouter()

	r.Handle("/", publicSites.Home)
	r.Handle("/dashboard", protectedSites.Dashboard)

	r.HandleFunc("/auth/login", userC.GetLogin).Methods("GET")
	r.HandleFunc("/auth/registration", userC.GetRegistration).Methods("GET")
	r.HandleFunc("/auth/hydra/login", hydraC.GetHydraLogin).Methods("GET")
	r.HandleFunc("/auth/hydra/consent", hydraC.GetHydraConsent).Methods("GET")

	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", logMw.Apply(identityMw.Apply(r))))
}

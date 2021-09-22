package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/monkjunior/poc-kratos-hydra/controllers"
	"github.com/monkjunior/poc-kratos-hydra/middlewares"

	"github.com/gorilla/mux"
	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Set to true in production. This ensures that a config.json file is provided before the application start")
	flag.Parse()
	kratosCfg, hPubCfg, hAdmCfg := LoadConfig(*boolPtr)
	k := kratosSDK.NewAPIClient(&kratosCfg)
	hCli := hydraSDK.NewHTTPClientWithConfig(nil, &hPubCfg)
	hAdm := hydraSDK.NewHTTPClientWithConfig(nil, &hAdmCfg)

	publicSites := controllers.NewPublicSites()
	protectedSites := controllers.NewProtectedSites()
	userC := controllers.NewUsers(k)
	hydraC := controllers.NewHydra(k, hCli, hAdm)

	logMw := middlewares.EntryLog{}
	identityMw := middlewares.Identity{KratosClient: k}

	r := mux.NewRouter()

	r.Handle("/", publicSites.Home)
	r.Handle("/dashboard", protectedSites.Dashboard)

	r.HandleFunc("/callback", userC.GetCallback).Methods("GET")
	r.HandleFunc("/auth/login", userC.GetLogin).Methods("GET")
	r.HandleFunc("/auth/registration", userC.GetRegistration).Methods("GET")
	r.HandleFunc("/auth/hydra/login", hydraC.GetHydraLogin).Methods("GET")
	r.HandleFunc("/auth/hydra/consent", hydraC.GetHydraConsent).Methods("GET")
	r.HandleFunc("/auth/hydra/consent", hydraC.PostHydraConsent).Methods("POST")

	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", logMw.Apply(identityMw.Apply(r))))
}

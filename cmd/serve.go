package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/monkjunior/poc-kratos-hydra/pkg/controllers"
	"github.com/monkjunior/poc-kratos-hydra/pkg/middlewares"
	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start an HTTP server which handles login and consent endpoint",
	Run:   runServeCmd,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) {
	kratosCfg, hPubCfg, hAdmCfg := config.GetAuthStackCfg()
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

	r.HandleFunc("/auth/login", userC.GetLogin).Methods("GET")
	r.HandleFunc("/auth/registration", userC.GetRegistration).Methods("GET")
	r.HandleFunc("/auth/hydra/login", hydraC.GetHydraLogin).Methods("GET")
	r.HandleFunc("/auth/hydra/consent", hydraC.GetHydraConsent).Methods("GET")
	r.HandleFunc("/auth/hydra/consent", hydraC.PostHydraConsent).Methods("POST")
	r.HandleFunc("/user/change-password", userC.PostChangePassword).Methods("POST")

	// Assets
	assetsHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

	fmt.Println("Listening at port 4435 ...")
	log.Fatal(http.ListenAndServe(":4435", logMw.Apply(identityMw.Apply(r))))
}

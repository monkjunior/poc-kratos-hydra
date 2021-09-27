package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/pkg/controllers"
	"github.com/monkjunior/poc-kratos-hydra/pkg/middlewares"
	kratosSDK "github.com/ory/kratos-client-go"
	"github.com/spf13/cobra"
)

// mockUICmd represents the mockUI command
var mockUICmd = &cobra.Command{
	Use:   "ui",
	Short: "Run a mock 1-st party app to perform login with Hydra.",
	Run:   runMockUICmd,
}

func init() {
	mockCmd.AddCommand(mockUICmd)
}

func runMockUICmd(cmd *cobra.Command, args []string) {
	kratosCfg, _, _ := GetAuthStackCfg()
	k := kratosSDK.NewAPIClient(&kratosCfg)

	mockUISites := controllers.NewMockUISites()
	userC := controllers.NewUsers(k)

	logMw := middlewares.EntryLog{}
	identityMw := middlewares.Identity{KratosClient: k}

	r := mux.NewRouter()

	r.HandleFunc("/", mockUISites.GetHome)
	r.HandleFunc("/callback", userC.GetCallback).Methods("GET")

	// Assets
	assetsHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

	fmt.Println("Listening at port 4436 ...")
	log.Fatal(http.ListenAndServe(":4436", logMw.Apply(identityMw.Apply(r))))
}

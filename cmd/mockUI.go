package cmd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/pkg/controllers"
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
	mockUISites := controllers.NewMockUISites()

	r := mux.NewRouter()

	r.HandleFunc("/", mockUISites.GetHome)
	r.HandleFunc("/callback", mockUISites.GetCallback).Methods("GET")

	// Assets
	assetsHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

	fmt.Println("Listening at port 4436 ...")
	err := http.ListenAndServe(":4436", r)
	if err != nil {
		panic(err)
	}
}

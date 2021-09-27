package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/monkjunior/poc-kratos-hydra/pkg/controllers"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// mockAPICmd represents the mockAPI command
var mockAPICmd = &cobra.Command{
	Use:   "api",
	Short: "Run a mock HTTP server to test HTTP Headers received from Oathkeeper mutator.",
	Run:   runMockAPICmd,
}

func init() {
	mockCmd.AddCommand(mockAPICmd)
}

func runMockAPICmd(cmd *cobra.Command, args []string) {
	mockAPI := controllers.NewMockAPI()

	r := mux.NewRouter()

	r.HandleFunc("/mock/api", mockAPI.GetAPI).Methods("GET")

	fmt.Println("Listening at port 4437 ...")
	log.Fatal(http.ListenAndServe(":4437", r))
}

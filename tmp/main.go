package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"

	kratosClient "github.com/ory/kratos-client-go"
)

func main() {
	// TODO: What is cookie jar
	cj, _ := cookiejar.New(nil)

	apiClient := kratosClient.NewAPIClient(&kratosClient.Configuration{
		Host:   "127.0.0.1:4455",
		Scheme: "http",
		DefaultHeader: map[string]string{
			"Env": "dev",
		},
		UserAgent: "monk_junior",
		Debug:     true,
		Servers: []kratosClient.ServerConfiguration{
			kratosClient.ServerConfiguration{
				URL: "/.ory/kratos/public",
			},
		},
		HTTPClient: &http.Client{Jar: cj},
	})

	ctx := context.Background()

	flow, res, err := apiClient.V0alpha1Api.InitializeSelfServiceRegistrationFlowForBrowsers(ctx).Execute()
	ExitOnError(err, res)
	PrintJSONPretty(flow)

	flowCSRF := flow.Ui.GetNodes()[0].Attributes.UiNodeInputAttributes.Value.(string)
	result, res, err := apiClient.V0alpha1Api.SubmitSelfServiceRegistrationFlow(ctx).Flow(
		flow.GetId(),
	).SubmitSelfServiceRegistrationFlowBody(
		kratosClient.SubmitSelfServiceRegistrationFlowBody{
			SubmitSelfServiceRegistrationFlowWithPasswordMethodBody: &kratosClient.SubmitSelfServiceRegistrationFlowWithPasswordMethodBody{
				CsrfToken: &flowCSRF,
				Method:    "password",
				Password:  "sonvn123",
				Traits: map[string]interface{}{
					"email": "admin@gmail.com",
				},
			},
		},
	).Execute()
	ExitOnError(err, res)
	PrintJSONPretty(result)
}

func ExitOnError(err error, res *http.Response) {
	if err == nil {
		return
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		out, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%s\n\nAn error occurred: %+v\n", out, err)
		os.Exit(1)
	}
	body, _ := json.MarshalIndent(json.RawMessage(bodyBytes), "", "  ")
	out, _ := json.MarshalIndent(err, "", "  ")
	fmt.Printf("%s\n\nAn error occurred: %+v\nbody: %s\n", out, err, body)
	os.Exit(1)
}

func PrintJSONPretty(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

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
	kratosClientClient "github.com/ory/kratos-client-go/client"
)

func main() {
	// Cookies can not be set automatically in testing environment because we do not request through a browser.
	// So we need to use cookie jar instead.
	//
	// More detailed information: https://stackoverflow.com/questions/31270461/what-is-the-difference-between-cookie-and-cookiejar
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

	flow2, res2, err2 := apiClient.V0alpha1Api.GetSelfServiceRegistrationFlow(ctx).Id(flow.GetId()).Execute()
	ExitOnError(err2, res2)
	PrintJSONPretty(flow2)

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

	kratosClientClient.Get
}

func ExitOnError(err error, res *http.Response) {
	if err == nil {
		return
	}
	var bodyBytes []byte
	if res != nil {
		bodyBytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			out, _ := json.MarshalIndent(err, "", "  ")
			fmt.Printf("%s\n\nAn error occurred: %+v\n", out, err)
			os.Exit(1)
		}
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

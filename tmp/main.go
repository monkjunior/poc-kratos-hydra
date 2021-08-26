package main

import (
	"context"
	"fmt"

	kratosClient "github.com/ory/kratos-client-go"
)

func main() {
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
	})
	request := apiClient.MetadataApi.IsAlive(context.Background())
	inlineRes, res, err := request.Execute()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Inline response", inlineRes)
	fmt.Println("Response", res)
}

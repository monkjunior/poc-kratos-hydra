package main

import (
	"fmt"

	hydraClient "github.com/ory/hydra-client-go/client"
)

func main() {
	fmt.Println("Hydra is awesome!")
	c := hydraClient.NewHTTPClientWithConfig(nil, &hydraClient.TransportConfig{
		Host: "127.0.0.1:4444",
		BasePath: "/",
		Schemes: []string{"http"},
	})
	isOK, err := c.Public.IsInstanceReady(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(isOK)
}

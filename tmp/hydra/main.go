package main

import (
	"fmt"

	"github.com/monkjunior/poc-kratos-hydra/rand"
	"golang.org/x/oauth2"
)

func main() {
	fmt.Println("Hydra is awesome!")
	//c := hydraClient.NewHTTPClientWithConfig(nil, &hydraClient.TransportConfig{
	//	Host:     "127.0.0.1:4444",
	//	BasePath: "/",
	//	Schemes:  []string{"http"},
	//})
	//isOK, err := c.Public.IsInstanceReady(nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(isOK)
	//
	//wellKnown, err := c.Public.WellKnown(nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(wellKnown)
	//
	//discoverOpenIDOK, err := c.Public.DiscoverOpenIDConfiguration(nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(discoverOpenIDOK)
	//
	// Init an OIDC authorization code flow
	// Should not use their own implementation
	// Use "golang.org/x/oauth2" package instead

	oauth2Config := oauth2.Config{
		ClientID:     "auth-code-client",
		ClientSecret: "secret",
		RedirectURL:  "http://127.0.0.1:4455/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://127.0.0.1:4444/oauth2/auth",
			TokenURL: "http://127.0.0.1:4444/oauth2/token",
		},

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{"openid"},
	}
	state, _ := rand.GenerateHydraState()
	authCodeURL := oauth2Config.AuthCodeURL(state)
	fmt.Println(authCodeURL)
}

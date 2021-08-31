package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
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
	// Use github.com/coreos/go-oidc package instead
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://127.0.0.1:4444/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", provider.Endpoint())

	oauth2Config := oauth2.Config{
		ClientID:     "auth-code-client",
		ClientSecret: "secret",
		RedirectURL:  "http://127.0.0.1:4455/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID},
	}

	url := oauth2Config.AuthCodeURL("a-random-string-for-state")
	fmt.Println(url)
}

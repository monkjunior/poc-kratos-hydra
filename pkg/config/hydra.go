package config

import (
	"strings"

	"github.com/monkjunior/poc-kratos-hydra/pkg/rand"
	hydraModels "github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
)

type HydraConfig struct {
	BrowserURL string               `json:"browserURL"`
	Public     HydraTransportConfig `json:"public"`
	Admin      HydraTransportConfig `json:"admin"`
	Client     HydraClient          `json:"client"`
}

type HydraTransportConfig struct {
	Host     string   `json:"host"`
	BasePath string   `json:"basePath"`
	Schemes  []string `json:"schemes"`
}

type HydraClient struct {
	ID            string   `json:"id"`
	Secret        string   `json:"secret"`
	GrantTypes    []string `json:"grantTypes"`
	ResponseTypes []string `json:"responseTypes"`
	Scopes        string   `json:"scopes"`
	CallbacksURL  []string `json:"callbacksURL"`
}

// GetBrowserAuthCodeURL generate authentication URL from Hydra client config
// This URL is used to init AuthZ code login flow.
// Docs: https://www.ory.sh/hydra/docs/concepts/login/
// In case we have a list of callback URLs, this function will use the first URL in the array.
func (h Config) GetBrowserAuthCodeURL() (url, state string) {
	oauth2Config := oauth2.Config{
		ClientID:     h.Hydra.Client.ID,
		ClientSecret: h.Hydra.Client.Secret,
		RedirectURL:  h.Hydra.Client.CallbacksURL[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.Hydra.BrowserURL + "/oauth2/auth",
			TokenURL: h.Hydra.BrowserURL + "/oauth2/token",
		},
		Scopes: strings.Split(h.Hydra.Client.Scopes, " "),
	}
	state, _ = rand.GenerateHydraState()
	authCodeURL := oauth2Config.AuthCodeURL(state)
	return authCodeURL, state
}

// GetInternalHydraOAuth2Config is used to export oauth2 config
// for internal purpose, for example: exchange token.
// In case we have a list of callback URLs, this function will use the first URL in the array.
func (h *Config) GetInternalHydraOAuth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     h.Hydra.Client.ID,
		ClientSecret: h.Hydra.Client.Secret,
		RedirectURL:  h.Hydra.Client.CallbacksURL[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.Hydra.Public.Schemes[0] + "://" + h.Hydra.Public.Host + "/oauth2/auth",
			TokenURL: h.Hydra.Public.Schemes[0] + "://" + h.Hydra.Public.Host + "/oauth2/token",
		},
		Scopes: strings.Split(h.Hydra.Client.Scopes, " "),
	}
}

// GetHydraOauth2Config take config from file then return *hydraModels.OAuth2Client
// Currently, we use this result to register an Oauth2 client in Hydra.
func (h *Config) GetHydraOauth2Config() *hydraModels.OAuth2Client {
	return &hydraModels.OAuth2Client{
		ClientID:      h.Hydra.Client.ID,
		ClientSecret:  h.Hydra.Client.Secret,
		GrantTypes:    h.Hydra.Client.GrantTypes,
		ResponseTypes: h.Hydra.Client.ResponseTypes,
		Scope:         h.Hydra.Client.Scopes,
		RedirectUris:  h.Hydra.Client.CallbacksURL,
	}
}

package config

import (
	"github.com/monkjunior/poc-kratos-hydra/pkg/rand"
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
	ID          string   `json:"id"`
	Secret      string   `json:"secret"`
	CallbackURL string   `json:"callbackURL"`
	Scopes      []string `json:"scopes"`
}

// GetBrowserAuthCodeURL generate authentication URL from Hydra client config
// This URL is used to init AuthZ code login flow.
// Docs: https://www.ory.sh/hydra/docs/concepts/login/
func (h Config) GetBrowserAuthCodeURL() (url, state string) {
	oauth2Config := oauth2.Config{
		ClientID:     h.Hydra.Client.ID,
		ClientSecret: h.Hydra.Client.Secret,
		RedirectURL:  h.Hydra.Client.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.Hydra.BrowserURL + "/oauth2/auth",
			TokenURL: h.Hydra.BrowserURL + "/oauth2/token",
		},
		Scopes: h.Hydra.Client.Scopes,
	}
	state, _ = rand.GenerateHydraState()
	authCodeURL := oauth2Config.AuthCodeURL(state)
	return authCodeURL, state
}

// GetInternalHydraOAuth2Config is used to export oauth2 config
// for internal purpose, for example: exchange token.
func (h *Config) GetInternalHydraOAuth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     h.Hydra.Client.ID,
		ClientSecret: h.Hydra.Client.Secret,
		RedirectURL:  h.Hydra.Client.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.Hydra.Public.Schemes[0] + "://" + h.Hydra.Public.Host + "/oauth2/auth",
			TokenURL: h.Hydra.Public.Schemes[0] + "://" + h.Hydra.Public.Host + "/oauth2/token",
		},
		Scopes: h.Hydra.Client.Scopes,
	}
}

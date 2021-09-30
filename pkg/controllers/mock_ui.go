package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/monkjunior/poc-kratos-hydra/pkg/views"
	"golang.org/x/oauth2"
)

var (
	hydraLoginURL string
	oauthState    string //TODO: this is just for test, need to reimplement the way we validate oauthState
)

func NewMockUISites() *MockUISites {
	return &MockUISites{
		Home: views.NewView("bootstrap", "mock_ui_home"),
	}
}

// MockUISites is a list of sites that our fake UI requires.
type MockUISites struct {
	Home *views.View
}

// MockSiteData stores auth code login URL
type MockSiteData struct {
	HydraLoginURL string
}

// GetHome just contain a login button to perform login with hydra
func (f *MockUISites) GetHome(w http.ResponseWriter, r *http.Request) {
	hydraLoginURL, oauthState = config.Cfg.GetBrowserAuthCodeURL()
	data := views.Data{
		Yield: MockSiteData{
			HydraLoginURL: hydraLoginURL,
		},
	}
	f.Home.Render(w, r, data)
}

// GetCallback receive authorization code and exchange token with Hydra, our OAuth2.0/OIDC server
// then it render token, and other result to viewer.
// GET /callback
func (u *Users) GetCallback(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("error")) > 0 {
		data := views.Data{
			Yield: CallbackForm{
				Error: &CallbackError{
					Name:        r.URL.Query().Get("error"),
					Description: r.URL.Query().Get("error_description"),
					Hint:        r.URL.Query().Get("error_hint"),
					Debug:       r.URL.Query().Get("error_debug"),
				},
			},
		}
		u.CallbackView.Render(w, r, data)
		return
	}

	// TODO: validate if states is matched
	isStatesMatched := true
	if !isStatesMatched {
		data := views.Data{
			Yield: CallbackForm{
				Error: &CallbackError{
					Name:        "States does not match",
					Description: "Expect A but received B",
				},
			},
		}
		u.CallbackView.Render(w, r, data)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := exchangeToken(r.Context(), code)
	if err != nil {
		data := views.Data{
			Yield: CallbackForm{
				Error: &CallbackError{
					Name:        "Failed to exchange token",
					Description: err.Error(),
				},
			},
		}
		u.CallbackView.Render(w, r, data)
		return
	}

	idToken := token.Extra("id_token")
	data := views.Data{
		Yield: CallbackForm{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry.Format(time.RFC1123),
			IDToken:      fmt.Sprintf("%v", idToken),
		},
	}
	u.CallbackView.Render(w, r, data)
	return
}

func exchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	oauth2Config := config.Cfg.GetInternalHydraOAuth2Config()
	return oauth2Config.Exchange(ctx, code)
}

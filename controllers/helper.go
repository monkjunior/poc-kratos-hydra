package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/spf13/viper"
)

var (
	KratosPublicURL                 = viper.GetString("baseUrl") + viper.GetString("kratos.publicBasePath")
	KratosSSLoginBrowserPath        = "/self-service/login/browser"
	KratosSSRegistrationBrowserPath = "/self-service/registration/browser"
)

// parseForm populates r.PostForm
//
// For all POST requests, ParseForm parses the raw data form from the request and updates
// r.PostForm.
func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

// parseValues decodes a map[string][]string to a struct.
func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}

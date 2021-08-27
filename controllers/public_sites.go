package controllers

import "github.com/monkjunior/poc-kratos-hydra/views"

func NewPublicSites() *PublicSites {
	return &PublicSites{
		Home: views.NewView("bootstrap", "home"),
	}
}

// PublicSites is a list of sites that do not require use to log in.
type PublicSites struct {
	Home *views.View
}

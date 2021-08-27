package controllers

import "github.com/monkjunior/poc-kratos-hydra/views"

func NewPublicSites() *PublicSites {
	return &PublicSites{
		Home: views.NewView("bootstrap", "home"),
	}
}

type PublicSites struct {
	Home *views.View
}

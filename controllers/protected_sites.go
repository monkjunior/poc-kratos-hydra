package controllers

import "github.com/monkjunior/poc-kratos-hydra/views"

func NewProtectedSites() *ProtectedSites {
	return &ProtectedSites{
		Dashboard: views.NewView("bootstrap", "dashboard"),
	}
}

type ProtectedSites struct {
	Dashboard *views.View
}

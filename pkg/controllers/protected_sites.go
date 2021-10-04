package controllers

import "github.com/monkjunior/poc-kratos-hydra/pkg/views"

func NewProtectedSites() *ProtectedSites {
	return &ProtectedSites{
		Dashboard: views.NewView("bootstrap", "dashboard"),
	}
}

// ProtectedSites is a list of sites that requires user logged in.
// Current we are use Oathkeeper to authenticate the session of requests coming.
type ProtectedSites struct {
	Dashboard *views.View
}

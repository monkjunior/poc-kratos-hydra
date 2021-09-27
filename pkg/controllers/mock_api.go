package controllers

import (
	"fmt"
	"net/http"
)

func NewMockAPI() *MockAPI {
	return &MockAPI{}
}

// MockAPI is used to test Oathkeeper mutator function
type MockAPI struct{}

// GetAPI prints out received HTTP headers
// GET /mock/api
func (u *MockAPI) GetAPI(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		fmt.Fprintln(w, name, values)
	}
}

package views

// Data is the top level structure that views expect data
// to come in.
type Data struct {
	IsSessionActive bool
	LogoutURL       string
	LogoutToken     string
	Yield           interface{}
}

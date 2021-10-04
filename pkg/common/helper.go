package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func LogOnError(err error, res *http.Response) {
	if err == nil {
		return
	}
	if res == nil {
		out, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%s\n\nAn error occurred: %+v\nbody: <nil>\n", out, err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		out, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%s\n\nAn error occurred: %+v\n", out, err)
		return
	}
	body, _ := json.MarshalIndent(json.RawMessage(bodyBytes), "", "  ")
	out, _ := json.MarshalIndent(err, "", "  ")
	fmt.Printf("%s\n\nAn error occurred: %+v\nbody: %s\n", out, err, body)
	return
}

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

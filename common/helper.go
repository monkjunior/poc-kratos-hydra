package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func PrintJSONPretty(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

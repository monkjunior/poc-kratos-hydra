package kratos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ExitOnError(err error, res *http.Response) {
	if err == nil {
		return
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		out, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%s\n\nAn error occurred: %+v\n", out, err)
		os.Exit(1)
	}
	body, _ := json.MarshalIndent(json.RawMessage(bodyBytes), "", "  ")
	out, _ := json.MarshalIndent(err, "", "  ")
	fmt.Printf("%s\n\nAn error occurred: %+v\nbody: %s\n", out, err, body)
	os.Exit(1)
}

func PrintJSONPretty(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

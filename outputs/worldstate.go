package outputs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/kataras/muxie"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

// IndexHandler test
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	fmt.Fprintf(w, "This is the <strong>%s</strong>", "index page")
}

// ProfileHandler test 2
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	name := muxie.GetParam(w, "platform")
	value, _ := intMap[name]
	w.Write(datasources.Apidata[value])
}

// ProfileHandler2 test 2
func ProfileHandler2(w http.ResponseWriter, r *http.Request) {
	lang := muxie.GetParam(w, "lang")
	name := muxie.GetParam(w, "platform")

	value, _ := intMap[name]
	messageJSON, _ := json.Marshal(parser.Testdata[value][lang])

	w.Write(messageJSON)
}

// ProfileHandler3 test 3
func ProfileHandler3(w http.ResponseWriter, r *http.Request) {
	lang := muxie.GetParam(w, "lang")
	name := muxie.GetParam(w, "platform")

	value, _ := intMap[name]
	messageJSON, _ := json.Marshal(parser.Testdata2[value][lang])

	w.Write(messageJSON)
}

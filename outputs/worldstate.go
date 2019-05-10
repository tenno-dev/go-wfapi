package outputs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/gorilla/mux"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

// IndexHandler test
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	fmt.Fprintf(w, "This is the <strong>%s</strong>", "index page")
}

// ProfileHandler test 2
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	value, _ := intMap[platform]
	w.Write(datasources.Apidata[value])
}

// ProfileHandler2 test 2
func ProfileHandler2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	messageJSON, _ := json.Marshal(parser.Testdata[value][token])

	w.Write(messageJSON)
}

// ProfileHandler3 test 3
func ProfileHandler3(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	messageJSON, _ := json.Marshal(parser.Testdata2[value][token])

	w.Write(messageJSON)
}

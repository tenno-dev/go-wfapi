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

// Everything test 2
func Everything(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	value, _ := intMap[platform]
	w.Header().Set("Content-Type", "application/json")
	w.Write(datasources.Apidata[value])
}

// DarvoDeals DarvoDeals
func DarvoDeals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	token = token[0:2]
	value, _ := intMap[platform]
	messageJSON, _ := json.Marshal(parser.Darvodata[value][token])

	w.Write(messageJSON)
}

// News Newsdata
func News(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	token = token[0:2]
	w.Header().Set("Content-Type", "application/json")
	messageJSON, _ := json.Marshal(parser.Newsdata[value][token])

	w.Write(messageJSON)
}

// Alerts Alertsdata
func Alerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	token = token[0:2]
	w.Header().Set("Content-Type", "application/json")
	messageJSON, _ := json.Marshal(parser.Alertsdata[value][token])

	w.Write(messageJSON)
}

// Fissures Fissuresdata
func Fissures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	token = token[0:2]
	w.Header().Set("Content-Type", "application/json")
	messageJSON, _ := json.Marshal(parser.Fissuresdata[value][token])

	w.Write(messageJSON)
}

// Nightwave Fissuresdata
func Nightwave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	platform := vars["platform"]
	token := r.Header.Get("Accept-Language")
	value, _ := intMap[platform]
	token = token[0:2]
	w.Header().Set("Content-Type", "application/json")
	messageJSON, _ := json.Marshal(parser.Nightwavedata[value][token])

	w.Write(messageJSON)
}

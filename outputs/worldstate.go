package outputs

import (
	"fmt"
	"net/http"

	"github.com/bitti09/go-wfapi/datasources"
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

package helper

import (
	"strings"

	"github.com/bitti09/go-wfapi/datasources"
)

// Langtranslate2 translate mission types - returns array  of 2
func Langtranslate2(src string, lang string) (ret [2]string) {
	var x1 [2]string
	x1[0] = src
	x1[1] = src
	src = strings.ToLower(src)
	result, ok := datasources.Languages[lang][src].(map[string]interface{})
	if ok != false {
		x1[0] = result["value"].(string)
		x1[1] = result["desc"].(string)

	}
	ret = x1
	return ret
}

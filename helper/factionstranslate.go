package helper

import (
	"github.com/bitti09/go-wfapi/datasources"
)

// Factionstranslate translate mission types
func Factionstranslate(src string, lang string) (ret string) {
	var x1 string
	x1 = src
	result, ok := datasources.FactionsData[lang][src].(map[string]interface{})
	if ok {
		x1 = result["value"].(string)
	}
	ret = x1
	return ret
}

package helper

import (
	"fmt"

	"github.com/bitti09/go-wfapi/datasources"
)

// Voidtranslate translate Fissure data
func Voidtranslate(src string, lang string) (ret [2]string) {
	var x1 [2]string
	x1[0] = src
	x1[1] = src

	result, ok := datasources.FissureModifiers[lang][src].(map[string]interface{})
	if ok {
		x1[0] = result["value"].(string)
		x1[1] = fmt.Sprintf("%.0f", result["num"].(float64))

	}
	ret = x1
	return ret
}

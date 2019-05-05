package helper

import "fmt"

var fissureModifiers map[string]map[string]interface{}

// Voidtranslate translate Fissure data
func Voidtranslate(src string, lang string) (ret [2]string) {
	var x1 [2]string
	x1[0] = src
	x1[1] = src

	result, ok := fissureModifiers[lang][src].(map[string]interface{})
	if ok != false {
		x1[0] = result["value"].(string)
		x1[1] = fmt.Sprintf("%.0f", result["num"].(float64))

	}
	ret = x1
	return ret
}

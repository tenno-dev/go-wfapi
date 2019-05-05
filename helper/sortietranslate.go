package helper

import "fmt"

var sortieloc map[string]map[string]interface{}
var sortiemodtypes map[string]map[string]interface{}
var sortiemoddesc map[string]map[string]interface{}
var sortiemodbosses map[string]map[string]interface{}

// Sortietranslate translate sortie data
func Sortietranslate(src string, langtype string, lang string) (ret [2]string) {
	if langtype == "sortiemod" {
		var x1 [2]string
		x1[0] = src
		x1[1] = src

		result, ok := sortiemodtypes[lang][src]
		if ok != false {
			x1[0] = result.(string)
		}
		result2, ok := sortiemoddesc[lang][src]
		if ok != false {
			x1[1] = result2.(string)
		}

		ret = x1
	}

	if langtype == "sortiemodboss" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src

		result, ok := sortiemodbosses[lang][src].(map[string]interface{})
		if ok != false {
			x1[0] = result["faction"].(string)
			x1[1] = result["name"].(string)

		}

		ret = x1

	}

	if langtype == "sortieloc" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src
		result, ok := sortieloc[lang][src].(map[string]interface{})
		fmt.Println("test2", sortieloc[lang][src])

		if ok != false {
			x1[0] = result["value"].(string)
			x1[1] = result["enemy"].(string)

		}
		/**/
		ret = x1

	}
	/**/
	return ret
}

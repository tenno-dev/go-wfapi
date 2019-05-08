package helper

import (
	"fmt"

	"github.com/bitti09/go-wfapi/datasources"
)

// Sortietranslate translate sortie data
func Sortietranslate(src string, langtype string, lang string) (ret [2]string) {
	if langtype == "sortiemod" {
		var x1 [2]string
		x1[0] = src
		x1[1] = src

		result, ok := datasources.Sortiemodtypes[lang][src]
		if ok != false {
			x1[0] = result.(string)
		}
		result2, ok := datasources.Sortiemoddesc[lang][src]
		if ok != false {
			x1[1] = result2.(string)
		}

		ret = x1
	}

	if langtype == "sortiemodboss" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src

		result, ok := datasources.Sortiemodbosses[lang][src].(map[string]interface{})
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
		result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})

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

// Sortietranslate2 - Sortie Rewards "translate"
func Sortietranslate2(src string, lang string) (ret string) {
	var x1 string

	x1 = src
	result := datasources.SortieRewards
	fmt.Println(string(result))
	x1 = string(result)
	ret = x1

	return ret

}

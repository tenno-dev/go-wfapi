package helper

import (
	"github.com/bitti09/go-wfapi/datasources"
	"github.com/tidwall/gjson"
)

// Sortietranslate translate sortie data
func Sortietranslate(src string, langtype string, lang string) (ret [2]string) {
	if langtype == "sortiemod" {
		var x1 [2]string
		x1[0] = src
		x1[1] = src

		result, ok := datasources.Sortiemodtypes[lang][src]
		if ok {
			x1[0] = result.(string)
		}
		result2, ok := datasources.Sortiemoddesc[lang][src]
		if ok {
			x1[1] = result2.(string)
		}

		ret = x1
	}

	if langtype == "sortiemodboss" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src

		result, ok := datasources.Sortiemodbosses[lang][src].(map[string]interface{})
		if ok {
			x1[0] = result["faction"].(string)
			x1[1] = result["name"].(string)

		}

		ret = x1

	}
	return ret
}

// Regiontranslate - Region "translate "
func Regiontranslate(src string, lang string) (ret [5]string) {

	var x1 [5]string
	nodesearch := "ExportRegions.#(uniqueName==" + "\"" + src + "\"" + ")" + ".name"
	planetsearch := "ExportRegions.#(uniqueName==" + "\"" + src + "\"" + ")" + ".systemName"

	Nodename := gjson.GetBytes(datasources.Regiondata[lang], nodesearch).String()
	Planetname := gjson.GetBytes(datasources.Regiondata[lang], planetsearch).String()

	// fmt.Println("test2:", Nodename)
	//	fmt.Println("test21:", string(datasources.Regiondata["en"]))

	x1[3] = src
	result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})
	if ok {
		x1[0] = Nodename
		x1[1] = Planetname
		x1[2] = result["enemy"].(string)
		x1[3] = result["type"].(string)
		x1[4] = result["value"].(string)

	}
	/**/
	ret = x1
	/**/
	return ret
}

// Sortietranslate2 - Sortie Rewards "translate "
func Sortietranslate2(src string, lang string) (ret []byte) {
	result := datasources.SortieRewards
	// fmt.Println(result)
	ret = result

	return ret

}

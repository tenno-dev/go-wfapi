package helper

import (
	"encoding/json"
	"fmt"

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
func Regiontranslate(src string, lang string) (ret [4]string) {

	var x1 [4]string
	//nodesearch := "ExportRegions.#(uniqueName==" + "\"" + src + "\"" + ")" + ".name"
	//planetsearch := "ExportRegions.#(uniqueName==" + "\"" + src + "\"" + ")" + ".systemName"
	//fmt.Println(src)
	//Nodename := gjson.GetBytes(datasources.Regiondata[lang], nodesearch).String()
	//Planetname := gjson.GetBytes(datasources.Regiondata[lang], planetsearch).String()

	// fmt.Println("test2:", Nodename)
	//	fmt.Println("test21:", string(datasources.Regiondata["en"]))

	x1[3] = src
	result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})
	if ok {
		//x1[0] = Nodename
		//x1[1] = Planetname
		x1[0] = result["enemy"].(string)
		x1[1] = result["type"].(string)
		x1[2] = result["value"].(string)

	}
	/**/
	//fmt.Println(x1)

	ret = x1
	/**/
	return ret
}

type SortieRewards1 struct {
	Id       string `json:"_id"`
	ItemName string `json:"itemName"`
	Rarity   string `json:"rarity"`
	Chance   int    `json:"chance"`
}
type SortieRewards2 struct {
	Items []SortieRewards1
}

// Sortietranslate2 - Sortie Rewards "translate "
func Sortietranslate2(src string, lang string) (ret []byte) {

	msg := []SortieRewards1{}
	//fmt.Println(datasources.SortieRewards[lang]["sortieRewards"])
	x2 := gjson.GetBytes(datasources.SortieRewards[lang], "sortieRewards").Raw
	//fmt.Println(x2)
	json.Unmarshal([]byte(x2), &msg)
	fmt.Println(msg)

	//ret1 := "out"
	//fmt.Println(ret1)
	//ret = result

	return []byte(x2)

}

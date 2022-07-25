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

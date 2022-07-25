package helper

import "github.com/bitti09/go-wfapi/datasources"

// Regiontranslate - Region "translate "
func Regiontranslate(src string, lang string) (ret string) {

	var x1 string
	x1 = src
	result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})
	if ok {
		x1 = result["value"].(string)

	}
	/**/
	//fmt.Println(x1)

	ret = x1
	/**/
	return ret
}

// Regiontranslate - Region "translate "
func Typetranslate(src string, lang string) (ret string) {

	var x1 string
	x1 = src
	result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})
	if ok {
		x1 = result["type"].(string)

	}
	/**/
	//fmt.Println(x1)

	ret = x1
	/**/
	return ret
}

// Regiontranslate - Region "translate "
func Enemytranslate(src string, lang string) (ret string) {

	var x1 string
	x1 = src
	result, ok := datasources.Sortieloc[lang][src].(map[string]interface{})
	if ok {
		x1 = result["enemy"].(string)

	}
	/**/
	//fmt.Println(x1)

	ret = x1
	/**/
	return ret
}

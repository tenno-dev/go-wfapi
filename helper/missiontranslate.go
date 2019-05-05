package helper

var missionTypes map[string]map[string]interface{}

// Missiontranslate translate mission types
func Missiontranslate(src string, lang string) (ret string) {
	var x1 string
	x1 = src

	result, ok := missionTypes[lang][src].(map[string]interface{})
	if ok != false {
		x1 = result["value"].(string)

	}
	ret = x1
	return ret
}

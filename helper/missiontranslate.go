package helper

import (
	"github.com/bitti09/go-wfapi/datasources"
)

// Missiontranslate translate mission types
func Missiontranslate(src string, lang string) (ret string) {
	var x1 string
	x1 = src
	result, ok := datasources.MissionTypes[lang][src].(map[string]interface{})
	if ok {
		x1 = result["value"].(string)
	}
	ret = x1
	return ret
}

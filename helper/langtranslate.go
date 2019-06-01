package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bitti09/go-wfapi/datasources"
)

// Langtranslate2 translate mission types - returns array  of 2
func Langtranslate2(src string, lang string) (ret [2]string) {
	var x1 [2]string
	x1[0] = src
	x1[1] = src
	src = strings.ToLower(src)
	result, ok := datasources.Languages[lang][src].(map[string]interface{})
	if ok != false {
		x1[0] = result["value"].(string)
		x1[1] = result["desc"].(string)

	}
	ret = x1
	return ret
}

// Langtranslate1 translate mission types - returns 1 result
func Langtranslate1(src string, lang string) (ret string) {
	var x1 string
	x1 = src
	src = strings.ToLower(src)
	result, ok := datasources.Languages[lang][src].(map[string]interface{})
	if ok != false {
		x1 = result["value"].(string)
	} else {
		src1 := strings.Replace(src, "storeitems/", "", -1)
		fmt.Println("translate error2", src1)
		result, ok := datasources.Languages[lang][src1].(map[string]interface{})
		if ok != false {
			x1 = result["value"].(string)
		} else {
			s := strings.Split(x1, "/")
			key := s[len(s)-1]
			re := regexp.MustCompile(`[A-Z]?[^A-Z]*`)
			split := re.FindAllString(key, -1)
			x1 = strings.Join(split, " ")
		}
	}
	ret = x1
	return ret
}

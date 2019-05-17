package helper

import (
	// "github.com/bitti09/go-wfapi/datasources"
)

// Getitemdetails test
func Getitemdetails(src string) (ret string) {
	var x1 string
	x1 = src
	/*result, ok := datasources.Nexusdata[src].(map[string]interface{})
	if ok != false {
		x1 = result["value"].(string)

	}*/
	ret = x1
	return ret
}

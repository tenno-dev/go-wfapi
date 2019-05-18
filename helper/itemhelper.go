package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	// "github.com/bitti09/go-wfapi/datasources"
)

// Getitemdetails test
func Getitemdetails(src string) (ret string) {
	client := &http.Client{}
	src = strings.ToLower(src)
	src = url.PathEscape(src)
	url := "http://localhost:8000/warframe/v1/items/" + src
	fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		ret = "error: not found"
				fmt.Println(err)

		return ret
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	/*result, ok := datasources.Nexusdata[src].(map[string]interface{})
	if ok != false {
		x1 = result["value"].(string)

	}*/
	ret = string(body[:])
	return ret
}

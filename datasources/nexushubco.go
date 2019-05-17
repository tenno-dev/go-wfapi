package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Nexusdata Result of LoadNexusdata
var Nexusdata []byte

// LoadNexusdata loads data from nexushub.co api
func LoadNexusdata(id1 string, id2 int) (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "https://api.nexushub.co/warframe/v1/items"

	fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Nexusdata = body[:]
	return
}

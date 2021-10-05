package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Kuvadata Kuvadata result
var Kuvadata []byte

// LoadKuvadata  Load Kuvadata from https://10o.io/kuvalog.json
func LoadKuvadata() (ret []byte) {
	client := &http.Client{}

	url := "https://10o.io/kuvalog.json"
	//fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Kuvadata = body[:]
	return
}

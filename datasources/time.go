package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Valistime Result of LoadTime
var Valistime []byte

// Cetustime Result of LoadTime
var Cetustime []byte

// Earthtime Result of LoadTime
var Earthtime []byte

// LoadTime loads time data from wfstat.us
func LoadTime() (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "https://api.warframestat.us/pc/vallisCycle"
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Valistime = body[:]

	//cetus
	url = "https://api.warframestat.us/pc/cetusCycle"
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Cetustime = body[:]

	//earth
	url = "https://api.warframestat.us/pc/earthCycle"
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Earthtime = body[:]
	return
}

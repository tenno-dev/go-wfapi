package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var CetusTimedata [4][]byte
var VallisTimedata [4][]byte
var CambionTimedata [4][]byte

func LoadCetusTimedata(id1 string, id2 int) (ret []byte) {

	client := &http.Client{}
	url := "https://api.warframestat.us/" + id1 + "/cetusCycle"
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
	CetusTimedata[id2] = body[:]
	return
}
func LoadVallisTimedata(id1 string, id2 int) (ret []byte) {
	client := &http.Client{}
	url := "https://api.warframestat.us/" + id1 + "/vallisCycle"
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
	VallisTimedata[id2] = body[:]
	return
}
func LoadCambionTimedata(id1 string, id2 int) (ret []byte) {

	client := &http.Client{}
	url := "https://api.warframestat.us/" + id1 + "/cambionCycle"
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
	CambionTimedata[id2] = body[:]

	return
}

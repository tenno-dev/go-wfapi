package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Apidata Result of LoadApidata
var Apidata [4][]byte

// Regiondata Result of LoadRegiondata
//var Regiondata [10][]byte
var Regiondata = make(map[string][]byte)

// LoadApidata loads data from Warframe.com api
func LoadApidata(id1 string, id2 int) (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "http://content.warframe.com/dynamic/worldState.php"
	if id1 != "pc" {
		url = "http://content." + id1 + ".warframe.com/dynamic/worldState.php"
	}
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
	Apidata[id2] = body[:]
	return
}

// LoadRegiondata loads data from Warframe.com api
func LoadRegiondata(id1 string, id2 int) (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "http://content.warframe.com/MobileExport/Manifest/ExportRegions.json"
	if id1 != "en" {
		url = "http://content.warframe.com/MobileExport/Manifest/ExportRegions_" + id1 + ".json"
	}
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
	Regiondata[id1] = body[:]
	return
}

package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Apidata Result of LoadApidata
var Apidata [4][]byte
var Timestamp [4]string

// Regiondata Result of LoadRegiondata
var Regiondata = make(map[string][]byte)

// Resourcedata Result of LoadItemdata
var Resourcedata = make(map[string][]byte)

// Upgradesdata Result of LoadItemdata
var Upgradesdata = make(map[string][]byte)

// LoadApidata loads data from Warframe.com api
func LoadApidata(id1 string, id2 int) (ret []byte) {
	client := &http.Client{}

	url := "http://content.warframe.com/dynamic/worldState.php"
	if id1 != "pc" {
		url = "http://content." + id1 + ".warframe.com/dynamic/worldState.php"
	}
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
	Apidata[id2] = body[:]
	Timestamp[id2] = time.Now().UTC().String()
	return
}

// LoadRegiondata loads data from Warframe.com api
func LoadRegiondata(id1 string, id2 int, wg *sync.WaitGroup) (ret []byte) {
	// WF API Source
	defer wg.Done()

	client := &http.Client{}

	url := "http://content.warframe.com/MobileExport/Manifest/ExportRegions.json"
	if id1 != "en" {
		url = "http://content.warframe.com/MobileExport/Manifest/ExportRegions_" + id1 + ".json"
	}
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
	Regiondata[id1] = body[:]
	return
}

// LoadResourcedata loads data from Warframe.com api
func LoadResourcedata(id1 string, id2 int, wg *sync.WaitGroup) (ret []byte) {
	// WF API Source
	defer wg.Done()

	client := &http.Client{}

	url := "http://content.warframe.com/MobileExport/Manifest/ExportResources.json"
	if id1 != "en" {
		url = "http://content.warframe.com/MobileExport/Manifest/ExportResources_" + id1 + ".json"
	}
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
	Resourcedata[id1] = body[:]
	return
}

// LoadUpgradesdata loads data from Warframe.com api
func LoadUpgradesdata(id1 string, id2 int, wg *sync.WaitGroup) (ret []byte) {
	// WF API Source
	defer wg.Done()

	client := &http.Client{}

	url := "http://content.warframe.com/MobileExport/Manifest/ExportResources.json"
	if id1 != "en" {
		url = "http://content.warframe.com/MobileExport/Manifest/ExportResources_" + id1 + ".json"
	}
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
	Upgradesdata[id1] = body[:]
	return
}

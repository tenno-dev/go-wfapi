package datasources

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// LangMap type fore json unmarshall
type LangMap map[string]interface{}

// LangMap2 d
type LangMap2 interface{}

// Sortieloc --
var Sortieloc = make(map[string]map[string]interface{})

// Sortiemodtypes Mod Type Lang strings
var Sortiemodtypes = make(map[string]map[string]interface{})

// Sortiemoddesc Mod Desc lang strings
var Sortiemoddesc = make(map[string]map[string]interface{})

// Sortiemodbosses Boss lang strings
var Sortiemodbosses = make(map[string]map[string]interface{})

// Sortielang --
var Sortielang = make(map[string]map[string]interface{}) // temp
// Sortielang1 --
var Sortielang1 map[string]interface{} // temp
// Sortielang2 --
var Sortielang2 map[string]string // temp
// FissureModifiers FissureModifiers lang strings
var FissureModifiers = make(map[string]map[string]interface{})

// MissionTypes MissionTypes lang strings
var MissionTypes = make(map[string]map[string]interface{})

// Loadlangdata load lang string from warframestat.us repo
func Loadlangdata(id1 string, id2 int) {

	client := &http.Client{}
	/*
		// arcanesData

		url := "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/arcanesData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/arcanesData.json"
		}
		// fmt.Println("url:", url)
		req, _ := http.NewRequest("GET", url, nil)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		arcanesData[id1] = string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// conclaveData
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/conclaveData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/conclaveData.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		conclaveData[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// eventsData
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/eventsData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/eventsData.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		eventsData[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// factionsData
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/factionsData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/factionsData.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		factionsData[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)


		// languages
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/languages.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/languages.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		languages[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// MissionTypes
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/MissionTypes.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/MissionTypes.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		MissionTypes[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// operationTypes
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/operationTypes.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/operationTypes.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		operationTypes[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// persistentEnemyData
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/persistentEnemyData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/persistentEnemyData.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		persistentEnemyData[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)
	*/
	// solNodes
	url := "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/" + id1 + "/solNodes.json"
	if id1 == "en" {
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/solNodes.json"
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
	fmt.Println(id1)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	Sortieloc[id1] = result

	_, _ = io.Copy(ioutil.Discard, res.Body)

	// sortieData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/" + id1 + "/sortieData.json"
	if id1 == "en" {
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/sortieData.json"
	}
	//fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &Sortielang)
	fmt.Println("test2")
	//		Sortielang = body.(map[string]interface{})
	Sortiemodtypes[id1] = make(LangMap)
	fmt.Println("test3")

	Sortiemodtypes[id1] = Sortielang["modifierTypes"]
	//	 fmt.Println("test2",Sortiemodtypes[id1]["SORTIE_MODIFIER_ARMOR"])
	Sortiemoddesc[id1] = Sortielang["modifierDescriptions"]
	Sortiemodbosses[id1] = Sortielang["bosses"]
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// FissureModifiers
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/" + id1 + "/FissureModifiers.json"
	if id1 == "en" {
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/FissureModifiers.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &result)
	FissureModifiers[id1] = result
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// MissionTypes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/" + id1 + "/MissionTypes.json"
	if id1 == "en" {
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/MissionTypes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &result)
	MissionTypes[id1] = result
	_, _ = io.Copy(ioutil.Discard, res.Body)

	/*
		// syndicatesData
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/syndicatesData.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/syndicatesData.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		syndicatesData[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// synthTargets
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/synthTargets.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/synthTargets.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		synthTargets[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// upgradeTypes
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/upgradeTypes.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/upgradeTypes.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		upgradeTypes[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// warframes
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/warframes.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/warframes.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		warframes[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)

		// weapons
		url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/weapons.json"
		if (id1 =="en"){
				url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/weapons.json"
		}
		// fmt.Println("url:", url)
		req, _ = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		weapons[id1]= string(body[:])
		_, _ = io.Copy(ioutil.Discard, res.Body)
	*/
	return
}

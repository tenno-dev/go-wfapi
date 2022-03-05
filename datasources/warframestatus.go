package datasources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	git "github.com/go-git/go-git/v5"
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

// FactionsData Factions lang strings
var FactionsData = make(map[string]map[string]interface{})

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

// Languages General Lang strings
var Languages = make(map[string]map[string]interface{})

// SortieRewards General Lang strings
var SortieRewards []byte

// Dirpath path to lang  git dir
var Dirpath = "./langsource/"

// InitLangDir test new ways
func InitLangDir() {
	_, err := os.Stat("./langsource/")

	if os.IsNotExist(err) {
		fmt.Println(err) //Shows error if file not exists
		_, _ = git.PlainClone("./langsource/", false, &git.CloneOptions{
			URL:      "https://github.com/WFCD/warframe-worldstate-data",
			Progress: os.Stdout,
		})
	}
	if !os.IsNotExist(err) {
		fmt.Println("path exist") // Shows success message like file is there
		r, _ := git.PlainOpen("./langsource")
		w, _ := r.Worktree()
		err3 := w.Pull(&git.PullOptions{
			Progress: os.Stdout,
			Force:    true,
		})
		fmt.Println(err3)

	}
}

// Loadlangdata load lang string from warframestat.us repo
func Loadlangdata(id1 string, id2 int, wg *sync.WaitGroup) {
	defer wg.Done()
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
	url := Dirpath + "data/" + id1 + "/solNodes.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "solNodes.json"
	}
	req, err := os.Open(url)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer req.Close()
	body, _ := ioutil.ReadAll(req)
	fmt.Println(id1)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	Sortieloc[id1] = result

	// sortieData
	url = Dirpath + "data/" + id1 + "/sortieData.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "sortieData.json"
	}
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	err = json.Unmarshal(body, &Sortielang)
	Sortiemodtypes[id1] = make(LangMap)
	Sortiemodtypes[id1] = Sortielang["modifierTypes"]
	Sortiemoddesc[id1] = Sortielang["modifierDescriptions"]
	Sortiemodbosses[id1] = Sortielang["bosses"]

	// FissureModifiers
	url = Dirpath + "data/" + id1 + "/fissureModifiers.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "fissureModifiers.json"
	}
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	json.Unmarshal([]byte(body), &result)
	FissureModifiers[id1] = result

	// MissionTypes
	url = Dirpath + "data/" + id1 + "/missionTypes.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "missionTypes.json"
	}
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	json.Unmarshal([]byte(body), &result)
	MissionTypes[id1] = result

	// languages
	url = Dirpath + "data/" + id1 + "/languages.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "languages.json"
	}
	// fmt.Println("url:", url)
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	json.Unmarshal([]byte(body), &result)
	Languages[id1] = result

	// FactionsData
	url = Dirpath + "data/" + id1 + "/factionsData.json"
	if id1 == "en" {
		url = Dirpath + "data/" + "factionsData.json"
	}
	// fmt.Println("url:", url)
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	json.Unmarshal([]byte(body), &result)
	FactionsData[id1] = result

	// sortieRewards
	url = Dirpath + "data/" + id1 + "/sortieRewards.json"
	if id1 != "en" {
		// url = "https://drops.warframestat.us/data/sortieRewards.json"
		return
	}
	// fmt.Println("url:", url)
	req, err = os.Open(url)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Close()
	body, _ = ioutil.ReadAll(req)
	var result2 string

	json.Unmarshal([]byte(body), &result2)
	SortieRewards = body

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

package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gorilla/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)

var missiontypelang = []byte(`{}`)
var factionslang = []byte(`{}`)
var locationlang = []byte(`{}`)
var apidata = []byte(`{}`)
var encjson = []byte(`{}`)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func loadapidata() {
	// Missiontypess
	url := "http://content.warframe.com/pc/dynamic/worldState.php"
	wfClient := http.Client{
		Timeout: time.Second * 20, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := wfClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	_, _ = io.Copy(ioutil.Discard, res.Body)
	apidata = body
}
func loadlangs() {
	// Missiontypes
	url := "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/missionTypes.json"
	wfClient := http.Client{
		Timeout: time.Second * 20, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := wfClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	_, _ = io.Copy(ioutil.Discard, res.Body)

	missiontypelang = body
	// Factions
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/factionsData.json"
	wfClient = http.Client{
		Timeout: time.Second * 20, // Maximum of 2 secs
	}

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr = wfClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()

	body, readErr = ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	_, _ = io.Copy(ioutil.Discard, res.Body)

	factionslang = body
	// Locations
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/solNodes.json"
	wfClient = http.Client{
		Timeout: time.Second * 20, // Maximum of 2 secs
	}

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr = wfClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()

	body, readErr = ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	_, _ = io.Copy(ioutil.Discard, res.Body)

	locationlang = body
}
func translatetest(src string, langtype string) (ret string) {
	var m map[string]interface{}

	if langtype == "faction" {
		err := json.Unmarshal(factionslang, &m)
		if err != nil {
			panic(err)
		}
		x1 := m[src].(map[string]interface{})["value"].(string)
		ret = string(x1)
	}
	if langtype == "missiontype" {
		err := json.Unmarshal(missiontypelang, &m)
		if err != nil {
			panic(err)
		}
		x1 := m[src].(map[string]interface{})["value"].(string)
		ret = string(x1)
	}
	if langtype == "location" {
		err := json.Unmarshal(locationlang, &m)
		if err != nil {
			panic(err)
		}
		x1 := m[src].(map[string]interface{})["value"].(string)
		ret = string(x1)
	}
	return ret

}

func main() {
	loadlangs()
	loadapidata()
	parseData()
	gofasion.SetJsonParser(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal, jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal)

	router := httprouter.New()
	router.GET("/wf", ShowData)

	log.Fatal(http.ListenAndServe(":8000", context.ClearHandler(router)))
}

func parseData() {
	type Eventmessage struct {
		LanguageCode string
		Message      string
	}
	type Event struct {
		ID         string
		Message    []Eventmessage
		URL        string
		Date       string
		priority   bool
		mobileonly bool
	}
	type Alerts struct {
		ID                  string
		Started             int
		Ends                int
		MissionType         string
		MissionFaction      string
		MissionLocation     string
		MinEnemyLevel       int
		MaxEnemyLevel       int
		EnemyWaves          int `json:",omitempty"`
		RewardCredits       int
		RewardItemMany      string `json:",omitempty"`
		RewardItemManyCount int    `json:",omitempty"`
		RewardItem          string `json:",omitempty"`
	}
	type Sortievariant struct {
		MissionType     string
		MissionMod      string
		MissionLocation string
		MissionTileset  string
	}
	type Sortie struct {
		ID       string
		Started  int
		Ends     int
		Boss     string
		Reward   string
		Variants []Sortievariant
		Twitter  bool
	}
	type SyndicateJobs struct {
		Jobtype            string
		Rewards            string
		MasterrankRequired int
		MinEnemyLevel      int
		MaxEnemyLevel      int
		XpReward           []int `json:"XPRewards"`
	}
	type SyndicateMissions struct {
		ID        string
		Started   int
		Ends      int
		Syndicate string
		Jobs      []SyndicateJobs
	}
	type Main struct {
		WorldSeed         string
		Alerts            []Alerts
		Events            []Event
		Sortie            []Sortie
		SyndicateMissions []SyndicateMissions
		Goals             string // As of  13.02.19 : no Goals reported from WF api
		Testresult        string // Test String
	}
	mainStruct := Main{Testresult: "Success"}

	fsion := gofasion.NewFasion(string(apidata[:]))
	// Event Section
	Eventarray := fsion.Get("Events").Array()
	fsion.ValueDefaultStr("-")
	fsion.ValueDefaultInt(0)
	for _, v := range Eventarray {
		id := v.Get("_id").Get("$oid").ValueStr()
		messagearray := v.Get("Messages").Array()
		var test []Eventmessage
		for i := range messagearray {
			test = append(test, Eventmessage{
				LanguageCode: messagearray[i].Get("LanguageCode").ValueStr(),
				Message:      messagearray[i].Get("Message").ValueStr()})
		}
		url := v.Get("Prop").ValueStr()
		date := v.Get("Date").Get("$date").Get("$numberLong").ValueStr()
		priority := v.Get("Priority").ValueBool()
		mobileonly := v.Get("MobileOnly").ValueBool()
		w := Event{ID: id, Message: test, URL: url, Date: date, priority: priority, mobileonly: mobileonly}
		mainStruct.Events = append(mainStruct.Events, w)
	}

	//Goals Section
	mainStruct.Goals = "No Goals"

	// WorldSeed Section
	mainStruct.WorldSeed = fsion.Get("WorldSeed").ValueStr()

	// Alerts Section
	Alertarray := fsion.Get("Alerts").Array()
	for _, v := range Alertarray {
		rewarditemsmany := ""
		rewarditem := ""
		rewarditemsmanycount := 0
		enemywaves := 0
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		missiontype := translatetest(v.Get("MissionInfo").Get("missionType").ValueStr(), "missiontype")
		missionfaction := translatetest(v.Get("MissionInfo").Get("faction").ValueStr(), "faction")
		missionlocation := translatetest(v.Get("MissionInfo").Get("location").ValueStr(), "location")
		minEnemyLevel := v.Get("MissionInfo").Get("minEnemyLevel").ValueInt()
		maxEnemyLevel := v.Get("MissionInfo").Get("maxEnemyLevel").ValueInt()
		enemywaves = v.Get("MissionInfo").Get("maxWaveNum").ValueInt()
		rewardcredits := v.Get("MissionInfo").Get("missionReward").Get("credits").ValueInt()
		rewarditemsmanyarray := v.Get("MissionInfo").Get("missionReward").Get("countedItems").Array()
		if len(rewarditemsmanyarray) != 0 {
			rewarditemsmany = rewarditemsmanyarray[0].Get("ItemType").ValueStr()
			rewarditemsmanycount = rewarditemsmanyarray[0].Get("ItemCount").ValueInt()
		}
		rewarditemarray := v.Get("MissionInfo").Get("missionReward").Get("items").Array()
		if len(rewarditemarray) != 0 {
			rewarditem = rewarditemarray[0].Get("items").ValueStr()
		}
		if int32(ended) > int32(time.Now().Unix()) {
			w := Alerts{id, started,
				ended, missiontype,
				missionfaction, missionlocation,
				minEnemyLevel, maxEnemyLevel, enemywaves,
				rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
			mainStruct.Alerts = append(mainStruct.Alerts, w)
		}
	}

	// Sorties Section
	Sortiearray := fsion.Get("Sorties").Array()

	for _, v := range Sortiearray {
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		boss := v.Get("Boss").ValueStr()
		reward := v.Get("Reward").ValueStr()
		variantarray := v.Get("Variants").Array()
		var variants []Sortievariant
		for i := range variantarray {
			variants = append(variants, Sortievariant{
				MissionType:     translatetest(variantarray[i].Get("missionType").ValueStr(), "missiontype"),
				MissionMod:      variantarray[i].Get("modifierType").ValueStr(),
				MissionLocation: translatetest(variantarray[i].Get("node").ValueStr(), "location"),
				MissionTileset:  variantarray[i].Get("tileset").ValueStr(),
			})
		}
		twitter := v.Get("Twitter").ValueBool()

		if int32(ended) > int32(time.Now().Unix()) {
			w := Sortie{ID: id, Started: started,
				Ends: ended, Boss: boss,
				Reward: reward, Variants: variants,
				Twitter: twitter}
			mainStruct.Sortie = append(mainStruct.Sortie, w)
		}
	}

	//SyndicateMissions Section
	SyndicateMissionsarray := fsion.Get("SyndicateMissions").Array()
	for _, v := range SyndicateMissionsarray {
		faction := v.Get("Tag").ValueStr()
		if faction == "SolarisSyndicate" || faction == "CetusSyndicate" {
			id := v.Get("_id").Get("$oid").ValueStr()
			started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
			ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
			syndicate := faction
			jobarray := v.Get("Jobs").Array()
			var jobs []SyndicateJobs
			for i := range jobarray {
				xparray := jobarray[i].Get("xpAmounts").Array()
				//xp1 :=
				jobs = append(jobs, SyndicateJobs{
					Jobtype:            jobarray[i].Get("jobType").ValueStr(),
					Rewards:            jobarray[i].Get("rewards").ValueStr(),
					MasterrankRequired: jobarray[i].Get("masteryReq").ValueInt(),
					MinEnemyLevel:      jobarray[i].Get("minEnemyLevel").ValueInt(),
					MaxEnemyLevel:      jobarray[i].Get("maxEnemyLevel").ValueInt(),
					XpReward:           []int{int(xparray[0].ValueInt()), int(xparray[1].ValueInt()), int(xparray[2].ValueInt())},
				})
			}
			w := SyndicateMissions{
				ID:        id,
				Started:   started,
				Ends:      ended,
				Syndicate: syndicate,
				Jobs:      jobs}
			mainStruct.SyndicateMissions = append(mainStruct.SyndicateMissions, w)
		}
	}
	encjson, _ = json.Marshal(mainStruct)
}

// ShowData
func ShowData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(encjson)
}

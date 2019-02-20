package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Anderson-Lu/gofasion/gofasion"
	emitter "github.com/emitter-io/go"
	jsoniter "github.com/json-iterator/go"
)

//current supported lang
var langid = map[string]int{
	"en": 0,
}
var platforms = [4]string{"pc", "ps4", "xb1", "swi"}
var missiontypelang = make([]byte, 1)
var factionslang = make([]byte, 1)
var locationlang = make([]byte, 1)
var apidata = make([][]byte, 4)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func loadapidata(id1 int) {
	// WF API Source
	url := "http://content.warframe.com/" + platforms[id1] + "/dynamic/worldState.php"
	fmt.Println("url:", url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	apidata[id1] = body
}
func loadlangs() {
	// Missiontypes EN
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
	// Factions EN
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

	// Locations EN
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
func translatetest(src string, langtype string, lang string) (ret string) {
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
	var wg sync.WaitGroup
	wg.Add(len(platforms))
	gofasion.SetJsonParser(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal, jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal)

	// mqtt client start
	o := emitter.NewClientOptions().AddBroker("tcp://0.0.0.0:8090")
	c := emitter.NewClient(o)
	sToken := c.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		panic("Error on Client.Connect(): " + sToken.Error().Error())
	}
	c.Subscribe("yyy", "test/")
	//mqtt client end

	loadlangs()
	for i := 0; i < len(platforms); i++ {

		loadapidata(i)
		parseInvasions(i, c)
		wg.Done()
	}
	wg.Wait()
	PrintMemUsage()
	/*		parseAlerts(apidata[x], c)
			parseNews(apidata[x], c)
			parseSorties(apidata[x], c)
			parseSyndicateMissions(apidata[x], c)
	parseActiveMissions(apidata[0], c)
	parseActiveMissions(apidata[1], c)
	parseActiveMissions(apidata[2], c)
	parseActiveMissions(apidata[3], c)*/

	PrintMemUsage()

}
func parseAlerts(apidata []byte, e emitter.Emitter) {
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
	fsion := gofasion.NewFasion(string(apidata[:][:]))
	var alerts []Alerts
	lang := string("en")
	fmt.Println(fsion.Get("WorldSeed").ValueStr())

	Alertarray := fsion.Get("Alerts").Array()
	for _, v := range Alertarray {
		rewarditemsmany := ""
		rewarditem := ""
		rewarditemsmanycount := 0
		enemywaves := 0
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		missiontype := translatetest(v.Get("MissionInfo").Get("missionType").ValueStr(), "missiontype", lang)
		missionfaction := translatetest(v.Get("MissionInfo").Get("faction").ValueStr(), "faction", lang)
		missionlocation := translatetest(v.Get("MissionInfo").Get("location").ValueStr(), "location", lang)
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
		if ended > int(time.Now().Unix()) {
			w := Alerts{id, started,
				ended, missiontype,
				missionfaction, missionlocation,
				minEnemyLevel, maxEnemyLevel, enemywaves,
				rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
			alerts = append(alerts, w)
		}

	}
	fmt.Println(len(alerts))
}
func parseNews(apidata []byte, e emitter.Emitter) {
	type Newsmessage struct {
		LanguageCode string
		Message      string
	}
	type News struct {
		ID         string
		Message    []Newsmessage
		URL        string
		Date       string
		priority   bool
		mobileonly bool
	}
	fsion := gofasion.NewFasion(string(apidata[:][:]))
	var news []News
	lang := string("en")
	Newsarray := fsion.Get("Events").Array()
	fsion.ValueDefaultStr("-")
	fsion.ValueDefaultInt(0)
	for _, v := range Newsarray {
		id := v.Get("_id").Get("$oid").ValueStr()
		messagearray := v.Get("Messages").Array()
		var test []Newsmessage
		for i := range messagearray {
			if messagearray[i].Get("LanguageCode").ValueStr() == lang {
				test = append(test, Newsmessage{
					LanguageCode: messagearray[i].Get("LanguageCode").ValueStr(),
					Message:      messagearray[i].Get("Message").ValueStr()})
			}
			// remove duplicate Items
			if len(test) > 1 {
				test = append(test[:1])
			}
		}
		url := v.Get("Prop").ValueStr()
		date := v.Get("Date").Get("$date").Get("$numberLong").ValueStr()
		priority := v.Get("Priority").ValueBool()
		mobileonly := v.Get("MobileOnly").ValueBool()
		w := News{ID: id, Message: test, URL: url, Date: date, priority: priority, mobileonly: mobileonly}
		if len(test) != 0 {
			news = append(news, w)
		}
	}
}
func parseSorties(apidata []byte, e emitter.Emitter) {
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
	fsion := gofasion.NewFasion(string(apidata[:][:]))
	var sortie []Sortie
	lang := string("en")
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
				MissionType:     translatetest(variantarray[i].Get("missionType").ValueStr(), "missiontype", lang),
				MissionMod:      variantarray[i].Get("modifierType").ValueStr(),
				MissionLocation: translatetest(variantarray[i].Get("node").ValueStr(), "location", lang),
				MissionTileset:  variantarray[i].Get("tileset").ValueStr(),
			})
		}

		twitter := v.Get("Twitter").ValueBool()

		w := Sortie{ID: id, Started: started,
			Ends: ended, Boss: boss,
			Reward: reward, Variants: variants,
			Twitter: twitter}
		sortie = append(sortie, w)
	}

}
func parseSyndicateMissions(apidata []byte, e emitter.Emitter) {
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

	fsion := gofasion.NewFasion(string(apidata[:]))
	var syndicates []SyndicateMissions
	//lang := string("en")
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
			syndicates = append(syndicates, w)
		}
	}

}
func parseActiveMissions(apidata []byte, e emitter.Emitter) {
	type ActiveMissions struct {
		ID          string
		Started     int
		Ends        int
		Region      int
		Node        string
		MissionType string
		Modifier    string
	}
	fsion := gofasion.NewFasion(string(apidata[:]))
	var mission []ActiveMissions
	lang := string("en")
	ActiveMissionsarray := fsion.Get("ActiveMissions").Array()
	fmt.Println(len(ActiveMissionsarray))

	for _, v := range ActiveMissionsarray {
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		region := v.Get("Region").ValueInt()
		node := translatetest(v.Get("Node").ValueStr(), "location", lang)
		missiontype := translatetest(v.Get("MissionType").ValueStr(), "missiontype", lang)
		modifier := v.Get("Modifier").ValueStr()

		w := ActiveMissions{
			ID:          id,
			Started:     started,
			Ends:        ended,
			Region:      region,
			Node:        node,
			MissionType: missiontype,
			Modifier:    modifier,
		}
		mission = append(mission, w)
	}
	fmt.Println(len(mission))

}
func parseInvasions(id12 int, e emitter.Emitter) {
	type Invasion struct {
		ID                  string
		Location            string
		MissionType         string
		Completed           bool
		Started             int
		AttackerRewardItem  string `json:",omitempty"`
		AttackerRewardCount int    `json:",omitempty"`
		AttackerMissionInfo string `json:",omitempty"`
		DefenderRewardItem  string `json:",omitempty"`
		DefenderRewardCount int    `json:",omitempty"`
		DefenderMissionInfo string `json:",omitempty"`
	}

	fsion := gofasion.NewFasion(string(apidata[id12][:]))
	var invasions []Invasion
	lang := string("en")
	fmt.Println(fsion.Get("WorldSeed").ValueStr())
	Invasionarray := fsion.Get("Invasions").Array()
	for _, v := range Invasionarray {
		attackeritem := ""
		attackeritemcount := 0
		defenderitem := ""
		defenderitemcount := 0
		id := v.Get("_id").Get("$oid").ValueStr()                                        //k
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000 //k
		location := translatetest(v.Get("Node").ValueStr(), "location", lang)
		missiontype := v.Get("LocTag").ValueStr()
		completed := v.Get("Completed").ValueBool()
		attackerrewardarray := v.Get("AttackerReward").Get("countedItems").Array()
		if len(attackerrewardarray) != 0 {
			attackeritem = attackerrewardarray[0].Get("ItemType").ValueStr()
			attackeritemcount = attackerrewardarray[0].Get("ItemCount").ValueInt()
		}
		attackerfaction := translatetest(v.Get("AttackerMissionInfo").Get("faction").ValueStr(), "faction", lang)
		defenderrewardarray := v.Get("DefenderReward").Get("countedItems").Array()
		if len(defenderrewardarray) != 0 {
			defenderitem = defenderrewardarray[0].Get("ItemType").ValueStr()
			defenderitemcount = defenderrewardarray[0].Get("ItemCount").ValueInt()
		}
		defenderfaction := translatetest(v.Get("DefenderMissionInfo").Get("faction").ValueStr(), "faction", lang)

		w := Invasion{id, location, missiontype, completed, started,
			attackeritem, attackeritemcount, attackerfaction,
			defenderitem, defenderitemcount, defenderfaction}
		invasions = append(invasions, w)

	}
	fmt.Println(len(invasions))
}

/*
func parseData(lang string, e emitter.Emitter) (ret []byte) {

	type Main struct {
		WorldSeed  string
		Goals      string // As of  13.02.19 : no Goals reported from WF api
		Testresult string // Test String
	}
	mainStruct := Main{Testresult: "Success"}

	fsion := gofasion.NewFasion(string(apidata[:]))

	//Goals Section
	mainStruct.Goals = "No Goals"

	// WorldSeed Section
	mainStruct.WorldSeed = fsion.Get("WorldSeed").ValueStr()

	// Alerts Section
	/* test publish to local mqtt broker
	var test, _ = json.Marshal(mainStruct.Alerts)
	var test2 = string(test[:])
	var test3 = `{ "msg":` + test2 + "}"
	//fmt.Println(test2)
	e.Publish("xx", "test/", test3)
	// test publish end
	// Sorties Section

	//SyndicateMissions Section
	// ActiveMissions Section

	ret, _ = json.Marshal(mainStruct)
	return ret
}
*/

// PrintMemUsage - only fo debug
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

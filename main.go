package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Anderson-Lu/gofasion/gofasion"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
var languageslang = make([]byte, 5486)
var apidata = make([][]byte, 4)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func loadapidata(id1 string) (ret []byte) {
	// WF API Source
	client := &http.Client{}
	url := "http://content.warframe.com/" + id1 + "/dynamic/worldState.php"
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

	return body
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

	// Languages EN
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/languages.json"
	wfClient = http.Client{
		Timeout: time.Second * 60, // Maximum of 2 secs
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
	fmt.Println("langfile loaded")
	languageslang = body
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
	if langtype == "languages" {
		var obj interface{} // Parse json into an interface{}
		err := json.Unmarshal(languageslang, &obj)
		if err != nil {
			panic(err)
		}
		m := obj.(map[string]interface{}) // Important: to access property
		foomap := m[strings.ToLower(src)]
		x1 := src
		if foomap != nil {
			x1 = foomap.(map[string]interface{})["value"].(string)
		}

		ret = x1
	}
	return ret

}

func main() {

	gofasion.SetJsonParser(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal, jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal)
	// mqtt client start
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("gotrivial")
	//opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	//opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe("test/topic", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	time.Sleep(1 * time.Second)
	//mqtt client end

	loadlangs()
	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		apidata[x] = loadapidata(v)
	}
	PrintMemUsage()

	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		//parseTests(apidata[x], v, nil)
		parseAlerts(apidata[x], v, c)
		parseNews(apidata[x], v, c)
		parseSorties(apidata[x], v, c)
		parseSyndicateMissions(apidata[x], v, c)
		parseActiveMissions(apidata[x], v, c)
		parseInvasions(apidata[x], v, c)
	}
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			for x, v := range platforms {
				fmt.Println("x:", x)
				fmt.Println("v:", v)
				//parseTests(apidata[x], v, nil)
				parseAlerts(apidata[x], v, c)
				parseNews(apidata[x], v, c)
				parseSorties(apidata[x], v, c)
				parseSyndicateMissions(apidata[x], v, c)
				parseActiveMissions(apidata[x], v, c)
				parseInvasions(apidata[x], v, c)

			}
			PrintMemUsage()

		}
	}()
	time.Sleep(time.Second * 180)

	//ticker.Stop()
	fmt.Println("Ticker stopped")
	/*		parseAlerts(apidata[x], c)
			parseNews(apidata[x], c)
			parseSorties(apidata[x], c)
			parseSyndicateMissions(apidata[x], c)
	parseActiveMissions(apidata[0], c)
	parseActiveMissions(apidata[1], c)
	parseActiveMissions(apidata[2], c)
	parseActiveMissions(apidata[3], c)*/

}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
func parseTests(apidata []byte, platform string, c mqtt.Client) {
	fsion := gofasion.NewFasion(string(apidata[:]))
	fmt.Println(fsion.Get("WorldSeed").ValueStr())
	topicf := "/wf/" + platform + "/tests"
	fmt.Println(topicf)
}

func parseAlerts(apidata []byte, platform string, c mqtt.Client) {
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
	fsion := gofasion.NewFasion(string(apidata[:]))
	var alerts []Alerts
	lang := string("en")
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
			rewarditemsmany = translatetest(rewarditemsmanyarray[0].Get("ItemType").ValueStr(), "languages", "en")
			rewarditemsmanycount = rewarditemsmanyarray[0].Get("ItemCount").ValueInt()
		}
		rewarditemarray := v.Get("MissionInfo").Get("missionReward").Get("items").Array()
		if len(rewarditemarray) != 0 {
			rewarditem = translatetest(rewarditemarray[0].Get("items").ValueStr(), "languages", "en")
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
	topicf := "/wf/" + platform + "/alerts"
	messageJSON, _ := json.Marshal(alerts)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

	fmt.Println(len(alerts))
}
func parseNews(apidata []byte, platform string, c mqtt.Client) {
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
		Image      string
		mobileonly bool
	}
	fsion := gofasion.NewFasion(string(apidata[:]))
	var news []News
	lang := string("en")
	Newsarray := fsion.Get("Events").Array()
	fsion.ValueDefaultStr("-")
	fsion.ValueDefaultInt(0)
	for _, v := range Newsarray {
		image := "http://n9e5v4d8.ssl.hwcdn.net/uploads/e0b4d18d3330bb0e62dcdcb364d5f004.png"
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
		if strings.HasPrefix(v.Get("ImageUrl").ValueStr(), "http") {
			image = v.Get("ImageUrl").ValueStr()
		}
		priority := v.Get("Priority").ValueBool()
		mobileonly := v.Get("MobileOnly").ValueBool()
		w := News{ID: id, Message: test, URL: url, Date: date, Image: image, priority: priority, mobileonly: mobileonly}
		if len(test) != 0 {
			news = append(news, w)
		}
	}
	topicf := "/wf/" + platform + "/news"
	messageJSON, _ := json.Marshal(news)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseSorties(apidata []byte, platform string, c mqtt.Client) {
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
	fsion := gofasion.NewFasion(string(apidata[:]))
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
	topicf := "/wf/" + platform + "/sorties"
	messageJSON, _ := json.Marshal(sortie)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseSyndicateMissions(apidata []byte, platform string, c mqtt.Client) {
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
					Jobtype:            translatetest(jobarray[i].Get("jobType").ValueStr(), "languages", "en"),
					Rewards:            translatetest(jobarray[i].Get("rewards").ValueStr(), "languages", "en"),
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
	topicf := "/wf/" + platform + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseActiveMissions(apidata []byte, platform string, c mqtt.Client) {
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
	topicf := "/wf/" + platform + "/missions"
	messageJSON, _ := json.Marshal(mission)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseInvasions(apidata []byte, platform string, c mqtt.Client) {
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

	fsion := gofasion.NewFasion(string(apidata[:]))
	var invasions []Invasion
	lang := string("en")
	//	text := fsion.Get("WorldSeed").ValueStr()

	Invasionarray := fsion.Get("Invasions").Array()
	for _, v := range Invasionarray {
		attackeritem := ""
		attackeritemcount := 0
		defenderitem := ""
		defenderitemcount := 0
		id := v.Get("_id").Get("$oid").ValueStr()
		//k
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
	topicf := "/wf/" + platform + "/invasions"
	messageJSON, _ := json.Marshal(invasions)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

// PrintMemUsage - only for debug
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

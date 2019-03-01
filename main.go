package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/buger/jsonparser"

	"github.com/Anderson-Lu/gofasion/gofasion"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/profile"
)

//current supported lang
var langid = map[string]int{
	"en": 0,
}
var platforms = [4]string{"pc", "ps4", "xb1", "swi"}
var missiontypelang map[string]interface{}
var factionslang map[string]interface{}
var locationlang map[string]interface{}
var sortiemodtypes map[string]interface{}
var sortiemoddesc map[string]interface{}
var sortiemodbosses map[string]interface{}
var sortieloc map[string]interface{}
var sortielang map[string]interface{}

var languageslang map[string]interface{}
var apidata = make([][]byte, 4)
var sortierewards = ""
var json = jsoniter.ConfigCompatibleWithStandardLibrary

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func loadapidata(id1 string) (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "https://api.warframestat.us/" + id1 + "/"
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
	return body[:]
}
func main() {
	defer profile.Start(profile.MemProfile).Stop()
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
	//mqtt client end

	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		apidata[x] = loadapidata(v)
		parseAlerts(x, v, c)

	}
	PrintMemUsage()
	/*
		c1 := cron.New()
		c1.AddFunc("@every 1m1s", func() {

			fmt.Println("Tick")
			for x, v := range platforms {
				fmt.Println("x:", x)
				fmt.Println("v:", v)
				apidata[x] = loadapidata(v)
				parseTests(x, v, nil)
				/*
					parseNews(x, v, c)
					parseSyndicateMissions(x, v, c)
					parseActiveMissions(x, v, c)
					parseInvasions(x, v, c)
					parseSorties(x, v, c)

			}
		})
		c1.Start()
	*/
	PrintMemUsage()

	// just for debuging - printing  full warframe api response
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/1", sayHello1)
	http.HandleFunc("/2", sayHello2)
	http.HandleFunc("/3", sayHello3)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}

}
func sayHello(w http.ResponseWriter, r *http.Request) {
	message := apidata[0][:]

	w.Write([]byte(message))
}
func sayHello1(w http.ResponseWriter, r *http.Request) {
	message := apidata[1][:]

	w.Write([]byte(message))
}
func sayHello2(w http.ResponseWriter, r *http.Request) {
	message := apidata[2][:]

	w.Write([]byte(message))
}
func sayHello3(w http.ResponseWriter, r *http.Request) {
	message := apidata[3][:]

	w.Write([]byte(message))
}
func parseTests(nor int, platform string, c mqtt.Client) {
	fsion := gofasion.NewFasion(string(apidata[nor]))
	fmt.Println(fsion.Get("WorldSeed").ValueStr())
	topicf := "/wf/" + platform + "/tests"
	fmt.Println(topicf)
}

func parseAlerts(platformno int, platform string, c mqtt.Client) {
	type Alerts struct {
		ID                  string
		Started             string
		Ends                string
		MissionType         string
		MissionFaction      string
		MissionLocation     string
		MinEnemyLevel       int64
		MaxEnemyLevel       int64
		EnemyWaves          int64 `json:",omitempty"`
		RewardCredits       int64
		RewardItemMany      string `json:",omitempty"`
		RewardItemManyCount int64  `json:",omitempty"`
		RewardItem          string `json:",omitempty"`
	}
	data := apidata[platformno]
	var alerts []Alerts
	_, _, _, erralert := jsonparser.Get(data, "alerts")
	fmt.Println(erralert)
	if erralert != nil {
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fmt.Println(jsonparser.GetString(value, "id"))

		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		missiontype, _ := jsonparser.GetString(value, "mission", "type")
		missionfaction, _ := jsonparser.GetString(value, "mission", "faction")
		missionlocation, _ := jsonparser.GetString(value, "mission", "node")
		minEnemyLevel, _ := jsonparser.GetInt(value, "mission", "minEnemyLevel")
		maxEnemyLevel, _ := jsonparser.GetInt(value, "mission", "maxEnemyLevel")
		enemywaves, _ := jsonparser.GetInt(value, "mission", "maxWaveNum")
		rewardcredits, _ := jsonparser.GetInt(value, "mission", "reward", "credits")
		rewarditemsmany, _ := jsonparser.GetString(value, "mission", "reward", "countedItems", "[0]", "type")
		rewarditemsmanycount, _ := jsonparser.GetInt(value, "mission", "reward", "countedItems", "[0]", "count")
		rewarditem, _ := jsonparser.GetString(value, "mission", "reward", "items", "[0]")

		w := Alerts{id, started,
			ended, missiontype,
			missionfaction, missionlocation,
			minEnemyLevel, maxEnemyLevel, enemywaves,
			rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
		alerts = append(alerts, w)

	}, "alerts")

	topicf := "/wf/" + platform + "/alerts"
	messageJSON, _ := json.Marshal(alerts)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

	fmt.Println(len(alerts))
}

/*
func parseNews(platformno int, platform string, c mqtt.Client) {
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
	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
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
func parseSorties(platformno int, platform string, c mqtt.Client) {
	type Sortievariant struct {
		MissionType     string
		MissionMod      string
		MissionLocation string
	}
	type Sortie struct {
		ID      string
		Started int
		Ends    int
		Boss    string
		//Reward   string
		Variants []Sortievariant
		Twitter  bool
	}
	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
	var sortie []Sortie
	lang := string("en")
	Sortiearray := fsion.Get("Sorties").Array()
	for _, v := range Sortiearray {

		id := v.Get("_id").Get("$oid").ValueStr()

		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		boss := v.Get("Boss").ValueStr()
		//reward := sortierewards
		variantarray := v.Get("Variants").Array()
		var variants []Sortievariant
		for i := range variantarray {
			variants = append(variants, Sortievariant{
				MissionType:     variantarray[i].Get("missionType").ValueStr(),
				MissionMod:      variantarray[i].Get("modifierType").ValueStr(),
				MissionLocation: variantarray[i].Get("node").ValueStr(),
			})
		}

		twitter := v.Get("Twitter").ValueBool()

		w := Sortie{ID: id, Started: started,
			Ends: ended, Boss: boss, Variants: variants,
			Twitter: twitter}
		sortie = append(sortie, w)
	}
	topicf := "/wf/" + platform + "/sorties"
	messageJSON, _ := json.Marshal(sortie)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseSyndicateMissions(platformno int, platform string, c mqtt.Client) {
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
	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
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
	topicf := "/wf/" + platform + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseActiveMissions(platformno int, platform string, c mqtt.Client) {
	type ActiveMissions struct {
		ID          string
		Started     int
		Ends        int
		Region      int
		Node        string
		MissionType string
		Modifier    string
	}
	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
	var mission []ActiveMissions
	lang := string("en")
	ActiveMissionsarray := fsion.Get("ActiveMissions").Array()
	fmt.Println(len(ActiveMissionsarray))

	for _, v := range ActiveMissionsarray {
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		region := v.Get("Region").ValueInt()
		node := v.Get("Node").ValueStr()
		missiontype := v.Get("MissionType").ValueStr()
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
func parseInvasions(platformno int, platform string, c mqtt.Client) {
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
		Completion          float32
	}

	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
	var invasions []Invasion
	lang := string("en")
	Invasionarray := fsion.Get("Invasions").Array()
	for _, v := range Invasionarray {
		attackeritem := ""
		attackeritemcount := 0
		defenderitem := ""
		defenderitemcount := 0
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000 //k
		location := v.Get("Node").ValueStr()
		missiontype := v.Get("LocTag").ValueStr()
		completed := v.Get("Completed").ValueBool()
		attackerrewardarray := v.Get("AttackerReward").Get("countedItems").Array()
		if len(attackerrewardarray) != 0 {
			attackeritem = attackerrewardarray[0].Get("ItemType").ValueStr()
			attackeritemcount = attackerrewardarray[0].Get("ItemCount").ValueInt()
		}
		attackerfaction := v.Get("AttackerMissionInfo").Get("faction").ValueStr()
		defenderrewardarray := v.Get("DefenderReward").Get("countedItems").Array()
		if len(defenderrewardarray) != 0 {
			defenderitem = defenderrewardarray[0].Get("ItemType").ValueStr()
			defenderitemcount = defenderrewardarray[0].Get("ItemCount").ValueInt()
		}
		defenderfaction := v.Get("DefenderMissionInfo").Get("faction").ValueStr()
		completion := calcCompletion(v.Get("Count").ValueInt(), v.Get("Goal").ValueInt(), attackerfaction)
		if v.Get("Completed").ValueBool() != true {
			w := Invasion{id, location, missiontype, completed, started,
				attackeritem, attackeritemcount, attackerfaction,
				defenderitem, defenderitemcount, defenderfaction, completion}
			invasions = append(invasions, w)
		}
	}
	topicf := "/wf/" + platform + "/invasions"
	messageJSON, _ := json.Marshal(invasions)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}*/
func calcCompletion(count int, goal int, attacker string) (complete float32) {
	y := float32((1 + float32(count)/float32(goal)))
	x := float32(y * 50)
	if attacker == "Infested" {
		x = float32(y * 100)

	}
	//fmt.Println(y)
	return x
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

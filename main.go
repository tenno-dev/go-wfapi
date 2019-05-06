package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/profile"
	"github.com/robfig/cron"
)

//current supported lang
var langpool = [10]string{"en", "de", "es", "fr", "it", "ko", "pl", "pt", "ru", "zh"}

// lang end
// platforms start
var platforms = [4]string{"pc", "ps4", "xb1", "swi"}

// platforms end
// var translationtype = [10]string{"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"}
var bempty = "[{}]"
var langtest = "en"

// LangMap start
type LangMap map[string]interface{}

// LangMap2 d

// todo
var arcanesData map[string]interface{}
var conclaveData map[string]interface{}
var eventsData map[string]interface{}
var factionsData map[string]interface{}
var languages = map[string]string{}
var operationTypes = map[string]string{}
var persistentEnemyData map[string]interface{}

var syndicatesData = map[string]string{}
var synthTargets map[string]interface{}
var upgradeTypes map[string]interface{}
var warframes map[string]interface{} //
var weapons map[string]interface{}

// Apidata downloaded api data
var Apidata [][]byte
var sortierewards = ""

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	// mqtt client start
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:8884/mqtt").SetClientID("gotrivial2")
	//opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername("x")
	opts.SetPassword("x")
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
	for x1, v1 := range langpool {
		fmt.Println("x1:", x1)
		fmt.Println("v1:", v1)
		datasources.Loadlangdata(v1, x1)
	} /**/
	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		datasources.LoadApidata(v, x)
		for x1, v1 := range langpool {
			fmt.Println("x1:", x1)
			fmt.Println("v1:", v1)
			parser.ParseSorties(x, v, c, v1)
			parser.ParseNews(x, v, c, v1)
			parser.ParseAlerts(x, v, c, v1)
			parser.ParseFissures(x, v, c, v1)
		}
		/*
			parseAlerts(x, v, c)
			parseNews(x, v, c)
			parseSorties(x, v, c)
			parseSyndicateMissions(x, v, c)
			parseInvasions(x, v, c)
			parseCycles(x, v, c)
			parseFissures(x, v, c)
			parseDarvo(x, v, c)
			parseEvents(x, v, c)
			parseNightwave(x, v, c)
		*/
		PrintMemUsage()

	}
	PrintMemUsage()

	c1 := cron.New()
	c1.AddFunc("@every 1m1s", func() {

		fmt.Println("Tick")
		for x, v := range platforms {
			fmt.Println("x:", x)
			fmt.Println("v:", v)
			datasources.LoadApidata(v, x)
			for x1, v1 := range langpool {
				fmt.Println("x1:", x1)
				fmt.Println("v1:", v1)
				parser.ParseSorties(x, v, c, v1)
				parser.ParseNews(x, v, c, v1)
				parser.ParseAlerts(x, v, c, v1)
				parser.ParseFissures(x, v, c, v1)
			}
			/*
				parseSyndicateMissions(x, v, c)
				parseInvasions(x, v, c)
				parseCycles(x, v, c)
				parseDarvo(x, v, c)
				parseEvents(x, v, c)
				parseNightwave(x, v, c)
			*/
			PrintMemUsage()
		}
		/*
			parseActiveMissions(x, v, c)
			parseInvasions(x, v, c)
		*/
	})
	c1.Start()
	PrintMemUsage()
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}

}

func arseCycles(platformno int, platform string, c mqtt.Client, lang string) {
	type Cycles struct {
		EathID         string
		EarthEnds      string
		EarthIsDay     bool
		EarthTimeleft  string
		CetusID        string
		CetusEnds      string
		CetusIsDay     bool
		CetusIsCetus   bool
		CetusTimeleft  string
		VallisID       string
		VallisEnds     string
		VallisIsWarm   bool
		VallisTimeleft string
	}
	data := Apidata[platformno]
	var cycles []Cycles
	fmt.Println("Cycles reached")
	//  Earth
	earthid, _ := jsonparser.GetString(data, "earthCycle", "id")
	earthends, _ := jsonparser.GetString(data, "earthCycle", "expiry")
	earthisday, _ := jsonparser.GetBoolean(data, "earthCycle", "isDay")
	earthtimeleft, _ := jsonparser.GetString(data, "earthCycle", "timeLeft")
	// Cetus
	cetusid, _ := jsonparser.GetString(data, "cetusCycle", "id")
	cetusends, _ := jsonparser.GetString(data, "cetusCycle", "expiry")
	cetusisday, _ := jsonparser.GetBoolean(data, "cetusCycle", "isDay")
	cetusiscetus, _ := jsonparser.GetBoolean(data, "cetusCycle", "isCetus")
	cetustimeleft, _ := jsonparser.GetString(data, "cetusCycle", "timeLeft")
	// Vallis
	vallisid, _ := jsonparser.GetString(data, "vallisCycle", "id")
	vallisends, _ := jsonparser.GetString(data, "vallisCycle", "expiry")
	vallisiswarm, _ := jsonparser.GetBoolean(data, "vallisCycle", "isDay")
	vallistimeleft, _ := jsonparser.GetString(data, "vallisCycle", "timeLeft")

	w := Cycles{earthid, earthends, earthisday, earthtimeleft,
		cetusid, cetusends, cetusisday, cetusiscetus, cetustimeleft,
		vallisid, vallisends, vallisiswarm, vallistimeleft}
	cycles = append(cycles, w)

	topicf := "/wf/" + platform + "/" + langtest + "/cycles"
	messageJSON, _ := json.Marshal(cycles)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseDarvo(platformno int, platform string, c mqtt.Client) {
	type DarvoDeals struct {
		ID              string
		Ends            string
		Item            string
		Price           int64
		DealPrice       int64
		DiscountPercent int64
		Stock           int64
		Sold            int64
	}
	data := Apidata[platformno]
	var deals []DarvoDeals
	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "dailyDeals")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/" + langtest + "/darvodeals"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Darvo reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		ended, _ := jsonparser.GetString(value, "expiry")
		item, _ := jsonparser.GetString(value, "item")
		originalprice, _ := jsonparser.GetInt(value, "originalPrice")
		dealprice, _ := jsonparser.GetInt(value, "salePrice")
		stock, _ := jsonparser.GetInt(value, "total")
		sold, _ := jsonparser.GetInt(value, "sold")
		discount, _ := jsonparser.GetInt(value, "discount")

		w := DarvoDeals{id, ended, item, originalprice, dealprice,
			discount, stock, sold}
		deals = append(deals, w)
	}, "dailyDeals")

	topicf := "/wf/" + platform + "/" + langtest + "/darvodeals"
	messageJSON, _ := json.Marshal(deals)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseNightwave(platformno int, platform string, c mqtt.Client) {

	type DailyChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
	type WeeklyChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
	type WeeklyEliteChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
	type Nightwave struct {
		ID                    string
		Ends                  string
		Started               string
		Season                int64
		Tag                   string
		Phase                 int64
		params                string
		possibleChallenges    string
		DailyChallenges       []DailyChallenges
		WeeklyChallenges      []WeeklyChallenges
		WeeklyEliteChallenges []WeeklyEliteChallenges
	}
	data := Apidata[platformno]
	var nightwave []Nightwave
	var dchallenge []DailyChallenges
	var wchallenge []WeeklyChallenges
	var welitechallenge []WeeklyEliteChallenges

	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "nightwave")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/" + langtest + "/nightwave"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Nightwave reached")
		return
	}
	fmt.Println("nightwave reached")
	id, _ := jsonparser.GetString(data, "nightwave", "id")
	ended, _ := jsonparser.GetString(data, "nightwave", "expiry")
	acvtivation, _ := jsonparser.GetString(data, "nightwave", "activation")
	season, _ := jsonparser.GetInt(data, "nightwave", "season")
	tag1, _ := jsonparser.GetString(data, "nightwave", "tag")
	phase, _ := jsonparser.GetInt(data, "nightwave", "phase")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		idc, _ := jsonparser.GetString(value, "id")
		endedc, _ := jsonparser.GetString(value, "expiry")
		activationc, _ := jsonparser.GetString(value, "activation")
		cdesc, _ := jsonparser.GetString(value, "desc")
		ctitle, _ := jsonparser.GetString(value, "title")
		reputation, _ := jsonparser.GetInt(value, "reputation")
		active, _ := jsonparser.GetBoolean(value, "active")
		daily, _ := jsonparser.GetBoolean(value, "isDaily")
		elite, _ := jsonparser.GetBoolean(value, "isElite")
		if daily == true {
			dailyc := DailyChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			dchallenge = append(dchallenge, dailyc)
		}
		if daily == false && elite == false {
			weekc := WeeklyChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			wchallenge = append(wchallenge, weekc)
		}
		if daily == false && elite == true {
			weekelitec := WeeklyEliteChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			welitechallenge = append(welitechallenge, weekelitec)
		}
	}, "nightwave", "activeChallenges")
	w := Nightwave{id, ended, acvtivation, season, tag1,
		phase, "", "", dchallenge, wchallenge, welitechallenge}
	nightwave = append(nightwave, w)
	topicf := "/wf/" + platform + "/" + langtest + "/nightwave"
	messageJSON, _ := json.Marshal(nightwave)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseEvents(platformno int, platform string, c mqtt.Client) {
	type Events struct {
		ID              string
		Started         string
		Ends            string
		Active          bool
		MaxScore        int64
		CurrScore       int64
		Health          float64
		Health1         string
		Faction         string
		Description     string
		Tooltip         string
		ConcurrentNodes string //subject to change when api  has data for it
		Rewarditem      string
		Rewardcredits   int64
		interimSteps    string
		progressSteps   string //subject to change when api  has data for it
		PersonalEvent   bool
		CommunityEvent  bool
		Expired         bool
	}
	data := Apidata[platformno]
	var events []Events

	fmt.Println("Events  reached")
	errfissures, _ := jsonparser.GetString(data, "Eveventsents")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/" + langtest + "/events"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Events reached")
		return
	}
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		rewarditems := ""
		rewardcredits := int64(0)
		health := float64(0)
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		active, _ := jsonparser.GetBoolean(value, "active")
		maximumScore, _ := jsonparser.GetInt(value, "maximumScore")
		currentScore, _ := jsonparser.GetInt(value, "currentScore")
		health1, _ := jsonparser.GetString(value, "health")

		health, _ = strconv.ParseFloat(health1, 64)
		faction, _ := jsonparser.GetString(value, "faction")
		description, _ := jsonparser.GetString(value, "description")
		tooltip, _ := jsonparser.GetString(value, "tooltip")
		expired, _ := jsonparser.GetBoolean(value, "expired")
		rewarditems, _ = jsonparser.GetString(value, "rewards", "[0]", "items", "[0]")
		rewardcredits, _ = jsonparser.GetInt(value, "rewards", "[0]", "credits", "[0]")
		personal, _ := jsonparser.GetBoolean(value, "isPersonal")
		commu, _ := jsonparser.GetBoolean(value, "isCommunity")

		w := Events{id, started, ended, active, maximumScore, currentScore, health, health1, faction,
			description, tooltip, "", rewarditems, rewardcredits, "", "", personal, commu, expired}
		events = append(events, w)
	}, "events")

	topicf := "/wf/" + platform + "/" + langtest + "/events"
	messageJSON, _ := json.Marshal(events)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

func parseSyndicateMissions(platformno int, platform string, c mqtt.Client) {
	type SyndicateJobs struct {
		Jobtype        string
		Rewards        []string
		MinEnemyLevel  int64
		MaxEnemyLevel  int64
		StandingReward []int64
	}
	type SyndicateMissions struct {
		ID        string
		Started   string
		Ends      string
		Syndicate string
		Jobs      []SyndicateJobs
	}
	data := Apidata[platformno]
	var syndicates []SyndicateMissions
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		syndicatecheck, _ := jsonparser.GetString(value, "syndicate")
		if syndicatecheck != "Ostrons" && syndicatecheck != "Solaris United" {
			return
		}
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		syndicate, _ := jsonparser.GetString(value, "syndicate")
		var jobs []SyndicateJobs
		jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			jobtype, _ := jsonparser.GetString(value1, "type")
			rewards := make([]string, 0)
			jsonparser.ArrayEach(value1, func(reward []byte, dataType jsonparser.ValueType, offset int, err error) {
				rewards = append(rewards, string(reward))

			}, "rewardPool")

			minEnemyLevel, _ := jsonparser.GetInt(value1, "enemyLevels", "[0]")
			maxEnemyLevel, _ := jsonparser.GetInt(value1, "enemyLevels", "[1]")
			standing1, _ := jsonparser.GetInt(value1, "standingStages", "[0]")
			standing2, _ := jsonparser.GetInt(value1, "standingStages", "[1]")
			standing3, _ := jsonparser.GetInt(value1, "standingStages", "[2]")
			jobs = append(jobs, SyndicateJobs{
				Jobtype:        jobtype,
				Rewards:        rewards,
				MinEnemyLevel:  minEnemyLevel,
				MaxEnemyLevel:  maxEnemyLevel,
				StandingReward: []int64{standing1, standing2, standing3},
			})
		}, "jobs")

		w := SyndicateMissions{
			ID:        id,
			Started:   started,
			Ends:      ended,
			Syndicate: syndicate,
			Jobs:      jobs}
		syndicates = append(syndicates, w)
	}, "syndicateMissions")

	topicf := "/wf/" + platform + "/" + langtest + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseInvasions(platformno int, platform string, c mqtt.Client) {
	type Invasion struct {
		ID                  string
		Location            string
		MissionType         string
		Completed           bool
		Started             string
		VsInfested          bool
		AttackerRewardItem  string `json:",omitempty"`
		AttackerRewardCount int64  `json:",omitempty"`
		AttackerMissionInfo string `json:",omitempty"`
		DefenderRewardItem  string `json:",omitempty"`
		DefenderRewardCount int64  `json:",omitempty"`
		DefenderMissionInfo string `json:",omitempty"`
		Completion          float64
	}

	data := Apidata[platformno]
	invasioncheck, _, _, _ := jsonparser.Get(data, "invasions")
	if len(invasioncheck) == 0 {
		topicf := "/wf/" + platform + "/" + langtest + "/invasions"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		return
	}
	var invasions []Invasion
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		iscomplete, _ := jsonparser.GetBoolean(value, "completed")
		if iscomplete != true {
			attackeritem := ""
			attackeritemcount := int64(0)
			defenderitem := ""
			defenderitemcount := int64(0)
			id, _ := jsonparser.GetString(value, "id")
			started, _ := jsonparser.GetString(value, "activation")
			location, _ := jsonparser.GetString(value, "node")
			missiontype, _ := jsonparser.GetString(value, "desc")
			completed, _ := jsonparser.GetBoolean(value, "completed")
			vsinfested, _ := jsonparser.GetBoolean(value, "vsInfestation")
			_, _, _, ierror := jsonparser.Get(value, "attackerReward", "countedItems", "[0]", "type")
			if ierror == nil {
				attackeritem, _ = jsonparser.GetString(value, "attackerReward", "countedItems", "[0]", "type")
				attackeritemcount, _ = jsonparser.GetInt(value, "attackerReward", "countedItems", "[0]", "count")
			}
			attackerfaction, _ := jsonparser.GetString(value, "attackingFaction")
			_, _, _, ierror2 := jsonparser.Get(value, "defenderReward", "countedItems", "[0]", "type")
			if ierror2 == nil {
				defenderitem, _ = jsonparser.GetString(value, "defenderReward", "countedItems", "[0]", "type")
				defenderitemcount, _ = jsonparser.GetInt(value, "defenderReward", "countedItems", "[0]", "count")
			}
			defenderfaction, _ := jsonparser.GetString(value, "defendingFaction")
			completion, _ := jsonparser.GetFloat(value, "completion")
			w := Invasion{id, location, missiontype, completed, started, vsinfested,
				attackeritem, attackeritemcount, attackerfaction,
				defenderitem, defenderitemcount, defenderfaction, completion}
			invasions = append(invasions, w)
		}
	}, "invasions")

	topicf := "/wf/" + platform + "/" + langtest + "/invasions"
	messageJSON, _ := json.Marshal(invasions)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

/*
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
	data := &Apidata[platformno]
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
	topicf := "/wf/" + platform + "/"+ langtest + "/missions"
	messageJSON, _ := json.Marshal(mission)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
*/
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

// FloatToString convert
func FloatToString(inputnum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputnum, 'f', 6, 64)
}

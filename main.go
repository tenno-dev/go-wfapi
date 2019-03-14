package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/robfig/cron"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/profile"
)

//current supported lang
var langid = map[string]int{
	"en": 0,
}
var bempty = "[{}]"

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

	// mqtt client start
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
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

	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		apidata[x] = loadapidata(v)
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

		PrintMemUsage()

	}
	PrintMemUsage()

	c1 := cron.New()
	c1.AddFunc("@every 1m1s", func() {

		fmt.Println("Tick")
		for x, v := range platforms {
			fmt.Println("x:", x)
			fmt.Println("v:", v)
			apidata[x] = loadapidata(v)
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
			PrintMemUsage()

		}
		/*
				parseActiveMissions(x, v, c)
				parseInvasions(x, v, c)

		}*/
	})
	c1.Start()

	PrintMemUsage()

	// just for debuging - printing  full warframe api response
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/1", sayHello1)
	http.HandleFunc("/2", sayHello2)
	http.HandleFunc("/3", sayHello3)

	if err := http.ListenAndServe(":9090", nil); err != nil {
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
		topicf := "/wf/" + platform + "/alerts"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
}
func parseCycles(platformno int, platform string, c mqtt.Client) {
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
	data := apidata[platformno]
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

	topicf := "/wf/" + platform + "/cycles"
	messageJSON, _ := json.Marshal(cycles)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseFissures(platformno int, platform string, c mqtt.Client) {
	type Fissures struct {
		ID              string
		Started         string
		Ends            string
		Active          bool
		MissionType     string
		MissionFaction  string
		MissionLocation string
		Tier            string
		TierLevel       int64
		Expired         bool
	}
	data := apidata[platformno]
	var fissures []Fissures
	fmt.Println("Fissues  reached")
	_, _, _, errfissures := jsonparser.Get(data, "fissures")
	if errfissures != nil {
		topicf := "/wf/" + platform + "/fissures"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		active, _ := jsonparser.GetBoolean(value, "active")
		location, _ := jsonparser.GetString(value, "node")
		missiontype, _ := jsonparser.GetString(value, "missionType")
		faction, _ := jsonparser.GetString(value, "enemy")
		tier, _ := jsonparser.GetString(value, "tier")
		tiernum, _ := jsonparser.GetInt(value, "tierNum")
		expired, _ := jsonparser.GetBoolean(value, "expired")

		w := Fissures{id, started, ended, active,
			missiontype, faction, location, tier, tiernum,
			expired}
		fissures = append(fissures, w)
	}, "fissures")

	topicf := "/wf/" + platform + "/fissures"
	messageJSON, _ := json.Marshal(fissures)
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
	data := apidata[platformno]
	var deals []DarvoDeals
	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "dailyDeals")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/darvodeals"
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

	topicf := "/wf/" + platform + "/darvodeals"
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
	data := apidata[platformno]
	var nightwave []Nightwave
	var dchallenge []DailyChallenges
	var wchallenge []WeeklyChallenges
	var welitechallenge []WeeklyEliteChallenges

	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "nightwave")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/nightwave"
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
	topicf := "/wf/" + platform + "/nightwave"
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
	data := apidata[platformno]
	var events []Events

	fmt.Println("Events  reached")
	errfissures, _ := jsonparser.GetString(data, "Eveventsents")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/events"
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

	topicf := "/wf/" + platform + "/events"
	messageJSON, _ := json.Marshal(events)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseNews(platformno int, platform string, c mqtt.Client) {
	type Newsmessage struct {
		LanguageCode string
		Message      string
	}
	type News struct {
		ID       string
		Message  string
		URL      string
		Date     string
		priority bool
		Image    string
	}
	data := apidata[platformno]
	_, _, _, ernews := jsonparser.Get(data, "news")
	if ernews != nil {
		fmt.Println("error ernews reached")
		return
	}
	var news []News

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		_, _, _, translationerr := jsonparser.Get(value, "translations", "en")
		if translationerr != nil {
			return
		}
		image := "http://n9e5v4d8.ssl.hwcdn.net/uploads/e0b4d18d3330bb0e62dcdcb364d5f004.png"
		message := ""
		id, _ := jsonparser.GetString(value, "id")

		message, _ = jsonparser.GetString(value, "translations", "en")
		url, _ := jsonparser.GetString(value, "link")
		image, _ = jsonparser.GetString(value, "imageLink")
		date, _ := jsonparser.GetString(value, "date")
		/**/
		priority, _ := jsonparser.GetBoolean(value, "priority")
		w := News{ID: id, Message: message, URL: url, Date: date, Image: image, priority: priority}
		news = append(news, w)
		topicf := "/wf/" + platform + "/news"
		messageJSON, _ := json.Marshal(news)
		token := c.Publish(topicf, 0, true, messageJSON)
		token.Wait()

	}, "news")
}
func parseSorties(platformno int, platform string, c mqtt.Client) {
	type Sortievariant struct {
		MissionType     string
		MissionMod      string
		MissionModDesc  string
		MissionLocation string
	}
	type Sortie struct {
		ID       string
		Started  string
		Ends     string
		Boss     string
		Faction  string
		Reward   string
		Variants []Sortievariant
		Active   bool
	}
	fmt.Println("reached sortie start")
	data := apidata[platformno]
	sortieactive, sortieerr := jsonparser.GetBoolean(data, "sortie", "active")
	if sortieerr != nil || sortieactive != true {
		topicf := "/wf/" + platform + "/sorties"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("reached sortie error")

		return
	}
	fmt.Println("reached sortie start2")

	var sortie []Sortie
	id, _ := jsonparser.GetString(data, "sortie", "id")
	started, _ := jsonparser.GetString(data, "sortie", "activation")
	ended, _ := jsonparser.GetString(data, "sortie", "expiry")
	boss, _ := jsonparser.GetString(data, "sortie", "boss")
	faction, _ := jsonparser.GetString(data, "sortie", "faction")
	reward, _ := jsonparser.GetString(data, "sortie", "rewardPool")
	var variants []Sortievariant

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		mtype, _ := jsonparser.GetString(value, "missionType")
		mmod, _ := jsonparser.GetString(value, "modifier")
		mmoddesc, _ := jsonparser.GetString(value, "modifierDescription")
		mloc, _ := jsonparser.GetString(value, "node")

		variants = append(variants, Sortievariant{
			MissionType:     mtype,
			MissionMod:      mmod,
			MissionModDesc:  mmoddesc,
			MissionLocation: mloc,
		})
	}, "sortie", "variants")
	active, _ := jsonparser.GetBoolean(data, "sortie", "active")
	w := Sortie{ID: id, Started: started,
		Ends: ended, Boss: boss, Faction: faction, Reward: reward, Variants: variants,
		Active: active}
	sortie = append(sortie, w)

	topicf := "/wf/" + platform + "/sorties"
	messageJSON, _ := json.Marshal(sortie)
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
	data := apidata[platformno]
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

	topicf := "/wf/" + platform + "/syndicates"
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

	data := apidata[platformno]
	invasioncheck, _, _, _ := jsonparser.Get(data, "invasions")
	if len(invasioncheck) == 0 {
		topicf := "/wf/" + platform + "/invasions"
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

	topicf := "/wf/" + platform + "/invasions"
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

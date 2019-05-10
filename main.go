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
	"github.com/bitti09/go-wfapi/outputs"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kataras/muxie"
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
	mux := muxie.NewMux()
	mux.PathCorrection = true
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
			parser.ParseSyndicateMissions(x, v, c, v1)
			parser.ParseInvasions(x, v, c, v1)
			parser.ParseDarvoDeal(x, v, c, v1)
			parser.ParseNightwave(x, v, c, v1)
			/*
				parseInvasions(x, v, c)
				parseCycles(x, v, c)
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
					parser.ParseSyndicateMissions(x, v, c, v1)
					parser.ParseInvasions(x, v, c, v1)
					parser.ParseDarvoDeal(x, v, c, v1)
					parser.ParseNightwave(x, v, c, v1)
				}
				/*
					parseInvasions(x, v, c)
					parseCycles(x, v, c)
					parseDarvo(x, v, c)
					parseEvents(x, v, c)
					parseNightwave(x, v, c)
				*/
				PrintMemUsage()
			}
		})
		c1.Start()
		PrintMemUsage()

		// static root, matches http://localhost:8080
		// or http://localhost:8080/ (even if PathCorrection is false).
		mux.HandleFunc("/", outputs.IndexHandler)

		// named parameter, matches /profile/$something_here
		// but NOT /profile/anything/here neither /profile
		// and /profile/ (if PathCorrection is true).
		mux.HandleFunc("/:platform", outputs.ProfileHandler)
		mux.HandleFunc("/darvo/:platform/:lang", outputs.ProfileHandler2)

		fmt.Println("Server started at http://localhost:9090")

		if err := http.ListenAndServe(":9090", mux); err != nil {
			panic(err)
		}

	}
}
func parseCycles(platformno int, platform string, c mqtt.Client, lang string) {
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

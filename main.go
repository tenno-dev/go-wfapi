package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/outputs"
	"github.com/bitti09/go-wfapi/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
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

	datasources.InitLangDir()	
	r := mux.NewRouter()

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
	} 
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

		r.HandleFunc("/", outputs.IndexHandler)

		// routes for multilang http output
		r.HandleFunc("/{platform}", outputs.ProfileHandler)
		r.HandleFunc("/{platform}/darvo/", outputs.ProfileHandler2)
		r.HandleFunc("/{platform}/news/", outputs.ProfileHandler3)
		r.HandleFunc("/{platform}/alerts/", outputs.ProfileHandler4)
		r.HandleFunc("/{platform}/fissures/", outputs.ProfileHandler5)
		r.HandleFunc("/{platform}/nightwave/", outputs.ProfileHandler6)

		fmt.Println("Server started at http://localhost:9090")

		if err := http.ListenAndServe("127.0.0.1:9090", r); err != nil {
			panic(err)
		}

	}
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

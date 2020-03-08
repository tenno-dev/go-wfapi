package main

import (
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"github.com/bitti09/go-wfapi/datasources"
	_ "github.com/bitti09/go-wfapi/docs"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/bitti09/go-wfapi/outputs"
	"github.com/bitti09/go-wfapi/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber"
	"github.com/robfig/cron/v3"
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

// Data  vars
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
	datasources.LoadKuvadata()
	datasources.LoadAnomalydata()
	datasources.InitLangDir()
	app := fiber.New()
	app.Settings.Prefork = true // Prefork enabled
	// mqtt client start
	opts := mqtt.NewClientOptions().AddBroker("ws://127.0.0.1:8083/mqtt").SetClientID("wf-mqtt")
	//opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	//opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	//mqtt client end
	for x1, v1 := range langpool {
		fmt.Println("x1:", x1)
		fmt.Println("v1:", v1)
		datasources.Loadlangdata(v1, x1)
		datasources.LoadRegiondata(v1, x1)
		datasources.LoadResourcedata(v1, x1)
		datasources.LoadUpgradesdata(v1, x1)

	}
	test1 := helper.Sortietranslate("SolNode30", "sortieloc", "de")
	fmt.Println(test1)

	datasources.LoadTime()
	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		datasources.LoadApidata(v, x)

		fmt.Println("LoadApidata:", v)
		for _, v1 := range langpool {
			parser.ParseGoals(x, v, c, v1)
			parser.ParseAnomaly(x, v, c, v1)
			parser.ParseKuva(x, v, c, v1)
			parser.ParseSorties(x, v, c, v1)
			parser.ParseNews(x, v, c, v1)
			parser.ParseAlerts(x, v, c, v1)
			parser.ParseFissures(x, v, c, v1)
			parser.ParseSyndicateMissions(x, v, c, v1)
			parser.ParseInvasions(x, v, c, v1)
			parser.ParseDarvoDeal(x, v, c, v1)
			parser.ParseNightwave(x, v, c, v1)
			parser.ParseVoidTrader(x, v, c, v1)
			parser.ParseProgress1(x, v, c, v1)
			parser.ParseTime(x, v, c, v1)
			parser.ParsePenemy(x, v, c, v1)

			/*
				parseCycles(x, v, c)
				parseEvents(x, v, c)
			*/
		}
		PrintMemUsage()
	}
	c0 := cron.New()
	c0.AddFunc("@every 5m1s", func() {
		datasources.LoadKuvadata()
		datasources.LoadAnomalydata()
	})

	c1 := cron.New()
	c1.AddFunc("@every 1m1s", func() {
		datasources.LoadTime()

		fmt.Println("Tick")
		for x, v := range platforms {
			datasources.LoadApidata(v, x)
			for _, v1 := range langpool {
				parser.ParseGoals(x, v, c, v1)
				parser.ParseAnomaly(x, v, c, v1)
				parser.ParseKuva(x, v, c, v1)
				parser.ParseSorties(x, v, c, v1)
				parser.ParseNews(x, v, c, v1)
				parser.ParseAlerts(x, v, c, v1)
				parser.ParseFissures(x, v, c, v1)
				parser.ParseSyndicateMissions(x, v, c, v1)
				parser.ParseInvasions(x, v, c, v1)
				parser.ParseDarvoDeal(x, v, c, v1)
				parser.ParseNightwave(x, v, c, v1)
				parser.ParseVoidTrader(x, v, c, v1)
				parser.ParseProgress1(x, v, c, v1)
				parser.ParseTime(x, v, c, v1)
				parser.ParsePenemy(x, v, c, v1)

			}
			/*
				parseCycles(x, v, c)
				parseEvents(x, v, c)
			*/
			PrintMemUsage()
		}
	})
	c1.Start()
	c0.Start()
	PrintMemUsage()

	// app.Get("/", outputs.IndexHandler)
	app.Get("/:platform", outputs.Everything)
	app.Get("/:platform/darvo/", outputs.DarvoDeals)
	app.Get("/:platform/news/", outputs.News)
	app.Get("/:platform/alerts/", outputs.Alerts)
	app.Get("/:platform/fissures/", outputs.Fissures)
	app.Get("/:platform/nightwave/", outputs.Nightwave)
	app.Get("/:platform/penemy/", outputs.Penemy) // temp naming
	app.Get("/:platform/fissures/", outputs.Fissures)
	app.Get("/:platform/nightwave/", outputs.Nightwave)

	app.Listen(8080)

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

package main

import (
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	_ "github.com/bitti09/go-wfapi/docs"
	"github.com/bitti09/go-wfapi/parser"
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

func main() {
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	datasources.InitLangDir()

	for x1, v1 := range langpool {
		//fmt.Println("x1:", x1)
		//fmt.Println("v1:", v1)
		wg2.Add(4)
		go datasources.Loadlangdata(v1, x1, &wg2)
		go datasources.LoadRegiondata(v1, x1, &wg2)
		go datasources.LoadResourcedata(v1, x1, &wg2)
		go datasources.LoadUpgradesdata(v1, x1, &wg2)
		wg2.Wait()
		fmt.Println("langpool goroutine exit")

	}
	//test1 := helper.Sortietranslate("SolNode30", "sortieloc", "de")
	//fmt.Println(test1)
	c0 := cron.New()
	c0.AddFunc("@every 5m1s", func() {
		datasources.LoadKuvadata()
		datasources.LoadAnomalydata()
	})

	c1 := cron.New()
	c1.AddFunc("@every 1m1s", func() {
		datasources.LoadTime()

		fmt.Println("Tick1s")
		for x, v := range platforms {
			datasources.LoadApidata(v, x)
			for _, v1 := range langpool {
				wg.Add(15)
				go parser.ParseGoals(x, v, v1, &wg)
				go parser.ParseAnomaly(x, v, v1, &wg)
				go parser.ParseKuva(x, v, v1, &wg)
				go parser.ParseSorties(x, v, v1, &wg)
				go parser.ParseNews(x, v, v1, &wg)
				go parser.ParseAlerts(x, v, v1, &wg)
				go parser.ParseFissures(x, v, v1, &wg)
				go parser.ParseSyndicateMissions(x, v, v1, &wg)
				go parser.ParseInvasions(x, v, v1, &wg)
				go parser.ParseDarvoDeal(x, v, v1, &wg)
				go parser.ParseNightwave(x, v, v1, &wg)
				go parser.ParseVoidTrader(x, v, v1, &wg)
				go parser.ParseProgress1(x, v, v1, &wg)
				go parser.ParseTime(x, v, v1, &wg)
				go parser.ParsePenemy(x, v, v1, &wg)
				wg.Wait()
				fmt.Println("parsepool v=" + v + " goroutine exit")

			}
			/*
				parseCycles(x, v, c)
				parseEvents(x, v, c)
			*/
		}
	})
	c1.Run()
	c1.Start()
	c0.Run()
	c0.Start()
	PrintMemUsage()

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

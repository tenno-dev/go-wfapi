package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/bitti09/go-wfapi/datasources"
	_ "github.com/bitti09/go-wfapi/docs"
	"github.com/bitti09/go-wfapi/outputs"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-co-op/gocron"
	httpSwagger "github.com/swaggo/http-swagger"
)

//current supported lang
//var langpool = [10]string{"en", "de", "es", "fr", "it", "ko", "pl", "pt", "ru", "zh"}
var langpool = [2]string{"en", "de"} // debug & testing only

// lang end
// platforms start
var platforms = [4]string{"pc", "ps4", "xb1", "swi"}

// platforms end
// var translationtype = [10]string{"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"}

// LangMap start
type LangMap map[string]interface{}

// Apidata downloaded api data
var Apidata [][]byte

// @title Tenno.dev  APIs
// @version 0.1
// @description Tenno.dev  APIs
// @BasePath /
// @host      api.tenno.dev

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	cors := cors.AllowAll()
	r.Use(cors.Handler)
	r.Use(middleware.Heartbeat("/"))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	s := gocron.NewScheduler(time.UTC)

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
	s.Every("1m1s").Do(func() {
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
				//fmt.Println("parsepool v=" + v + " goroutine exit")

			}
			/*
				parseCycles(x, v, c)
				parseEvents(x, v, c)
			*/
		}
	})
	s.Every("5m1s").Do(func() {
		datasources.LoadKuvadata()
		datasources.LoadAnomalydata()
	})

	r.Get("/{platform}", outputs.Everything)                            // looks ok
	r.Get("/{platform}/test", outputs.Everything2)                      // debug
	r.Get("/{platform}/darvo", outputs.DarvoDeals)                      // looks ok
	r.Get("/{platform}/news", outputs.News)                             // looks ok
	r.Get("/{platform}/alerts", outputs.Alerts)                         // null response
	r.Get("/{platform}/fissures", outputs.Fissures)                     // MissionFaction & MissionLocation empty
	r.Get("/{platform}/nightwave", outputs.Nightwave)                   // looks ok
	r.Get("/{platform}/penemy", outputs.Penemy)                         // null response
	r.Get("/{platform}/invasion", outputs.Invasion)                     // empty location
	r.Get("/{platform}/time", outputs.Time)                             // empty response
	r.Get("/{platform}/sortie", outputs.Sortie)                         // looks ok
	r.Get("/{platform}/voidtrader", outputs.Voidtrader)                 // looks ok
	r.Get("/{platform}/syndicate", outputs.SyndicateMission)            // Rewards response is basic like "Narmer Table C Rewards" - needs  more work  in reward parser
	r.Get("/{platform}/anomaly", outputs.AnomalyData)                   // needs work
	r.Get("/{platform}/progress", outputs.Progress1)                    // looks ok
	r.Get("/{platform}/event", outputs.Event)                           // looks ok
	r.Get("/{platform}/arbitrationmission", outputs.ArbitrationMission) // null response - is intended as source is empty
	r.Get("/{platform}/kuvamission", outputs.KuvaMission)               // null response - is intended as source is empty
	/**/
	s.StartAsync()

	PrintMemUsage()
	fmt.Println("PrintMemUsage  created")
	http.ListenAndServe(":3000", r)

	//r.Run(":8080")
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

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

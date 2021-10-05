package parser

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// AnomalyData struct
type AnomalyData struct {
	ID               string
	Node             string
	Start            string // int to string
	Startstring      string // int to utc string
	Ends             string // int to string
	EndString        string // int to utc string
	Projection       string // int to string
	Projectionstring string // int to utc string
}

// AnomalyDataSet for http export
var AnomalyDataSet = make(map[int]map[string][]AnomalyData)

// ParseAnomaly Parse current Darvo Deal
func ParseAnomaly(platformno int, platform string, c mqtt.Client, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := AnomalyDataSet[platformno]; !ok {
		AnomalyDataSet[platformno] = make(map[string][]AnomalyData)
	}
	value := datasources.Anomalydata[:]
	var anoma []AnomalyData

	//id, _ := jsonparser.GetString(value, "id")
	id := "1"
	started, _ := jsonparser.GetInt(value, "start")
	started1 := strconv.FormatInt(started, 10)
	startedstring := time.Unix(started, 0).String()

	ended, _ := jsonparser.GetInt(value, "end")
	endedstring := time.Unix(ended, 0).String()
	ended1 := strconv.FormatInt(ended, 10)
	//start, _ := time.Parse(timeFormat, started)
	// end, _ := time.Parse(timeFormat, ended)
	node, _ := jsonparser.GetString(value, "name")
	projection, _ := jsonparser.GetInt(value, "projection")
	projection1 := strconv.FormatInt(projection, 10)
	projectionstring := time.Unix(projection, 0).String()

	a := AnomalyData{id, node, started1, startedstring, ended1, endedstring, projection1, projectionstring}
	anoma = append(anoma, a)

	topica := "wf/" + lang + "/" + platform + "/anomaly"
	messageJSONa, _ := json.Marshal(anoma)
	AnomalyDataSet[platformno][lang] = anoma
	tokena := c.Publish(topica, 0, true, messageJSONa)
	tokena.Wait()

}

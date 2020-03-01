package parser

import (
	"encoding/json"
	"strconv"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// AnomalyData struct
type AnomalyData struct {
	ID         string
	Start      string
	Ends       string
	Node       string
	projection string
}

// KuvaMission for http export
var AnomalyDataSet = make(map[int]map[string][]AnomalyData)

// ParseKuva Parse current Darvo Deal
func ParseAnomaly(platformno int, platform string, c mqtt.Client, lang string) {
	if _, ok := AnomalyDataSet[platformno]; !ok {
		AnomalyDataSet[platformno] = make(map[string][]AnomalyData)
	}
	value := datasources.Anomalydata[:]
	var anoma []AnomalyData

	//id, _ := jsonparser.GetString(value, "id")
	id := "1"
	started, _ := jsonparser.GetInt(value, "start")
	started1 := strconv.FormatInt(started, 10)
	ended, _ := jsonparser.GetInt(value, "end")
	ended1 := strconv.FormatInt(ended, 10)
	//start, _ := time.Parse(timeFormat, started)
	// end, _ := time.Parse(timeFormat, ended)
	node, _ := jsonparser.GetString(value, "name")
	projection, _ := jsonparser.GetInt(value, "projection")
	projection1 := strconv.FormatInt(projection, 10)

	a := AnomalyData{id, started1, ended1, node, projection1}
	anoma = append(anoma, a)

	topica := "wf/" + lang + "/" + platform + "/anomaly"
	messageJSONa, _ := json.Marshal(anoma)
	AnomalyDataSet[platformno][lang] = anoma
	tokena := c.Publish(topica, 0, true, messageJSONa)
	tokena.Wait()
}

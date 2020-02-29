package parser

import (
	"encoding/json"

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
	data := datasources.Anomalydata[:]
	var anoma []AnomalyData

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//id, _ := jsonparser.GetString(value, "id")
		id := "1"
		started, _ := jsonparser.GetString(value, "start")
		ended, _ := jsonparser.GetString(value, "end")
		//start, _ := time.Parse(timeFormat, started)
		// end, _ := time.Parse(timeFormat, ended)
		node, _ := jsonparser.GetString(value, "name")
		projection, _ := jsonparser.GetString(value, "projection")

		a := AnomalyData{id, started, ended, node, projection}
		anoma = append(anoma, a)
	})

	topica := "wf/" + lang + "/" + platform + "/anomaly"
	messageJSONa, _ := json.Marshal(anoma)
	AnomalyDataSet[platformno][lang] = anoma
	tokena := c.Publish(topica, 0, true, messageJSONa)
	tokena.Wait()
}

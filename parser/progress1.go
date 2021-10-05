package parser

import (
	"encoding/json"
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Progress1 - Progress1
type Progress1 struct {
	P1 float64
	P2 float64
	P3 float64
}

// ParseProgress1 Parse Void trader
func ParseProgress1(platformno int, platform string, c mqtt.Client, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	data := datasources.Apidata[platformno]
	var progress1 []Progress1

	_, _, _, pro1err := jsonparser.Get(data, "ProjectPct")
	if pro1err != nil {
		topicf := "wf/" + lang + "/" + platform + "/progress"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		// fmt.Println("reached progress error")

		return
	}
	// fmt.Println("reached progress start")

	p1, _ := jsonparser.GetFloat(data, "ProjectPct", "[0]")
	p2, _ := jsonparser.GetFloat(data, "ProjectPct", "[1]")
	p3, _ := jsonparser.GetFloat(data, "ProjectPct", "[2]")

	w := Progress1{P1: p1, P2: p2,
		P3: p3}
	progress1 = append(progress1, w)

	topicf := "wf/" + lang + "/" + platform + "/progress"
	messageJSON, _ := json.Marshal(progress1)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}

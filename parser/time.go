package parser

import (
	"encoding/json"
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Time1 - Time Base
type Time1 struct {
	Cetus  []Time2
	Vallis []Time2
	Earth  []Time2
}

// Time2 - Time Details
type Time2 struct {
	Start string `json:",omitempty"`
	End   string
	State string
}

// ParseTime Parse Void trader
func ParseTime(platformno int, platform string, c mqtt.Client, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	data1 := datasources.Cetustime
	data2 := datasources.Valistime
	data3 := datasources.Earthtime

	var time1 []Time1
	var cetus []Time2
	var valis []Time2
	var earth []Time2

	cetusbegin, _ := jsonparser.GetString(data1, "activation")
	cetusend, _ := jsonparser.GetString(data1, "expiry")
	cetusstate, _ := jsonparser.GetString(data1, "state")
	cetus = append(cetus, Time2{
		Start: cetusbegin,
		End:   cetusend,
		State: cetusstate,
	})

	valisbegin, _ := jsonparser.GetString(data2, "activation")
	valisend, _ := jsonparser.GetString(data2, "expiry")
	valisstate, _ := jsonparser.GetString(data2, "state")
	valis = append(valis, Time2{
		Start: valisbegin,
		End:   valisend,
		State: valisstate,
	})

	earthbegin, _ := jsonparser.GetString(data3, "activation")
	earthend, _ := jsonparser.GetString(data3, "expiry")
	earthstate, _ := jsonparser.GetString(data3, "state")
	earth = append(earth, Time2{
		Start: earthbegin,
		End:   earthend,
		State: earthstate,
	})

	w := Time1{Cetus: cetus, Vallis: valis,
		Earth: earth}
	time1 = append(time1, w)

	topicf := "wf/" + lang + "/" + platform + "/time"
	messageJSON, _ := json.Marshal(time1)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}

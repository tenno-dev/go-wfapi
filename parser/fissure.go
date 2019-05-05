package parser

import (
	"encoding/json"
	"fmt"

	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ParseFissures  parse Fissure data
func ParseFissures(platformno int, platform string, c mqtt.Client, lang string) {
	type Fissures struct {
		ID              string
		Started         string
		Ends            string
		Active          bool
		MissionType     string
		MissionFaction  string
		MissionLocation string
		Tier            string
		TierLevel       string
		Expired         bool
	}
	data := Apidata[platformno]
	var fissures []Fissures
	fmt.Println("Fissues  reached")
	_, _, _, errfissures := jsonparser.Get(data, "ActiveMissions")
	if errfissures != nil {
		topicf := "/wf/" + lang + "/" + platform + "/fissures"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("Fissues 2 reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		active := true
		location1, _ := jsonparser.GetString(value, "Node")
		location := helper.Sortietranslate(location1, "sortieloc", lang)
		missiontype1, _ := jsonparser.GetString(value, "MissionType")
		missiontype := helper.Missiontranslate(missiontype1, lang)
		tier1, _ := jsonparser.GetString(value, "Modifier")
		tier := helper.Voidtranslate(tier1, lang)
		expired, _ := jsonparser.GetBoolean(value, "expired")

		w := Fissures{id, started, ended, active,
			missiontype, location[1], location[0], tier[0], tier[1],
			expired}
		fissures = append(fissures, w)
	}, "ActiveMissions")

	topicf := "/wf/" + lang + "/" + platform + "/fissures"
	messageJSON, _ := json.Marshal(fissures)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
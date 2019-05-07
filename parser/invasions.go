package parser

import (
	"encoding/json"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ParseInvasions parse active Invasions
func ParseInvasions(platformno int, platform string, c mqtt.Client, lang string) {
	type Invasion struct {
		ID          string
		Location    string
		MissionType string
		Completed   bool
		Started     string
		//VsInfested          bool
		AttackerRewardItem  string `json:",omitempty"`
		AttackerRewardCount int64  `json:",omitempty"`
		AttackerMissionInfo string `json:",omitempty"`
		DefenderRewardItem  string `json:",omitempty"`
		DefenderRewardCount int64  `json:",omitempty"`
		DefenderMissionInfo string `json:",omitempty"`
		Completion          float32
	}

	data := datasources.Apidata[platformno]
	invasioncheck, _, _, _ := jsonparser.Get(data, "Invasions")
	if len(invasioncheck) == 0 {
		topicf := "/wf/" + lang + "/" + platform + "/invasions"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		return
	}
	var invasions []Invasion
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		iscomplete, _ := jsonparser.GetBoolean(value, "Completed")
		if iscomplete != true {
			attackeritem := ""
			attackeritemcount := int64(0)
			defenderitem := ""
			defenderitemcount := int64(0)
			id, _ := jsonparser.GetString(value, "_id", "$oid")
			started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
			location1, _ := jsonparser.GetString(value, "Node")
			location := helper.Sortietranslate(location1, "sortieloc", lang)
			missiontype, _ := jsonparser.GetString(value, "LocTag")
			completed, _ := jsonparser.GetBoolean(value, "Completed")
			//vsinfested, _ := jsonparser.GetBoolean(value, "vsInfestation")
			_, _, _, ierror := jsonparser.Get(value, "AttackerReward", "countedItems", "[0]", "ItemType")
			if ierror == nil {
				attackeritem, _ = jsonparser.GetString(value, "AttackerReward", "countedItems", "[0]", "ItemType")
				attackeritemcount, _ = jsonparser.GetInt(value, "AttackerReward", "countedItems", "[0]", "ItemCount")
			}
			attackerfaction1, _ := jsonparser.GetString(value, "AttackerMissionInfo", "faction")
			attackerfaction := helper.Sortietranslate(attackerfaction1, "sortiemodboss", lang)
			//attackerfaction := helper.Sortietranslate(attackerfaction1, "sortieloc", lang)
			_, _, _, ierror2 := jsonparser.Get(value, "DefenderReward", "countedItems", "[0]", "type")
			if ierror2 == nil {
				defenderitem, _ = jsonparser.GetString(value, "DefenderReward", "countedItems", "[0]", "ItemType")
				defenderitemcount, _ = jsonparser.GetInt(value, "DefenderReward", "countedItems", "[0]", "ItemCount")
			}
			defenderfaction1, _ := jsonparser.GetString(value, "DefenderMissionInfo", "faction")
			defenderfaction := helper.Sortietranslate(defenderfaction1, "sortiemodboss", lang)
			completion1, _ := jsonparser.GetInt(value, "Count")
			completion2, _ := jsonparser.GetInt(value, "Goal")
			completion := calcCompletion(completion1, completion2, attackerfaction[0])
			w := Invasion{id, location[0], missiontype, completed, started,
				attackeritem, attackeritemcount, attackerfaction[0],
				defenderitem, defenderitemcount, defenderfaction[0], completion}
			invasions = append(invasions, w)
		}
	}, "Invasions")

	topicf := "/wf/" + lang + "/" + platform + "/invasions"
	messageJSON, _ := json.Marshal(invasions)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func calcCompletion(count int64, goal int64, attacker string) (complete float32) {
	y := float32((1 + float32(count)/float32(goal)))
	x := float32(y * 50)
	if attacker == "Infested" || attacker == "FC_INFESTATION" {
		x = float32(y * 100)

	}
	//fmt.Println(y)
	return x
}

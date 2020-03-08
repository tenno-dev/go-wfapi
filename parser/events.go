package parser

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type InterimReward struct {
	Item    string
	Credits string
	XP      string
}

// EventsData struct
type EventsData struct {
	Debug             string
	ID                string
	Name              string
	Start             string          // int to string
	Ends              string          // int to string
	Location          string          `json:",omitempty"`
	Count             string          `json:",omitempty"`
	HealthPct         string          `json:",omitempty"`
	Goal              string          `json:",omitempty"`
	Mainreward        string          `json:",omitempty"`
	Mainrewardcredits string          `json:",omitempty"`
	Mainrewardxp      string          `json:",omitempty"`
	InterimGoalsteps  []string        `json:",omitempty"`
	InterimRewards    []InterimReward `json:",omitempty"`
}

// Eventdata - Event data
var Eventdata = make(map[int]map[string][]EventsData)
var interim []InterimReward
var interimsteps []string

// ParseGoals parsing Events data (Called Goals in warframe api)
func ParseGoals(platformno int, platform string, c mqtt.Client, lang string) {
	if _, ok := Eventdata[platformno]; !ok {
		Eventdata[platformno] = make(map[string][]EventsData)
	}
	data := datasources.Apidata[platformno]
	_, _, _, ernews := jsonparser.Get(data, "Goals")
	if ernews != nil {
		// fmt.Println("error ernews reached")
		return
	}
	var event []EventsData

	debug, _ := jsonparser.GetString(data, "Goals")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// id
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		name, _ := jsonparser.GetString(value, "Desc")
		name1 := helper.Langtranslate1(name, lang)

		started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")

		node, _ := jsonparser.GetString(value, "Node")

		//scorev, _ := jsonparser.GetString(value, "ScoreVar")
		count1, _ := jsonparser.GetInt(value, "Count")
		count2 := strconv.FormatInt(count1, 10)
		health, _ := jsonparser.GetFloat(value, "HealthPct")
		health1 := fmt.Sprintf("%.2f", health)
		goal, _ := jsonparser.GetInt(value, "Goal")
		goal1 := strconv.FormatInt(goal, 10)
		reward, _ := jsonparser.GetString(value, "Reward", "items", "[0]")
		rewards1 := helper.Langtranslate1(reward, lang)

		rewardcredits, _ := jsonparser.GetInt(value, "Reward", "credits")
		rewardcredits1 := strconv.FormatInt(rewardcredits, 10)
		rewardxp, _ := jsonparser.GetInt(value, "Reward", "xp")
		rewardxp1 := strconv.FormatInt(rewardxp, 10)

		var interim []InterimReward
		jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.GetString(value1, "countedItems", "[0]", "ItemType")
			item1 := helper.Langtranslate1(item, lang)

			if item == "" {
				item, _ = jsonparser.GetString(value1, "items", "[0]")
				item1 = helper.Langtranslate1(item, lang)

			}
			xp, _ := jsonparser.GetInt(value1, "xp")
			xp1 := strconv.FormatInt(xp, 10)
			credits, _ := jsonparser.GetInt(value1, "credits")
			credits1 := strconv.FormatInt(credits, 10)
			wt := InterimReward{item1, credits1, xp1}
			interim = append(interim, wt)
		}, "InterimRewards")
		var interimsteps []string

		jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			interimsteps = append(interimsteps, string(value1))
		}, "InterimGoals")

		w := EventsData{Debug: debug, ID: id, Name: name1, Start: started, Ends: ended, Location: node, Count: count2, HealthPct: health1, Goal: goal1, Mainreward: rewards1, Mainrewardxp: rewardxp1, Mainrewardcredits: rewardcredits1, InterimGoalsteps: interimsteps, InterimRewards: interim}
		event = append(event, w)
	}, "Goals")
	topicf := "wf/" + lang + "/" + platform + "/goals"
	Eventdata[platformno][lang] = event
	messageJSON, _ := json.Marshal(event)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

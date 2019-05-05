package parser

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var apidata [][]byte
// ParseAlerts parsing Alerts data
func ParseAlerts(platformno int, platform string, c mqtt.Client, lang string) {
	type Alerts struct {
		ID                  string
		Started             string
		Ends                string
		MissionType         string
		MissionFaction      string
		MissionLocation     string
		MinEnemyLevel       int64
		MaxEnemyLevel       int64
		EnemyWaves          int64 `json:",omitempty"`
		RewardCredits       int64
		RewardItemMany      string `json:",omitempty"`
		RewardItemManyCount int64  `json:",omitempty"`
		RewardItem          string `json:",omitempty"`
	}
	data := apidata[platformno]
	var alerts []Alerts
	_, _, _, erralert := jsonparser.Get(data, "Alerts")
	fmt.Println(erralert)
	if erralert != nil || erralert == nil { // disable  parsing until api returns data
		topicf := "/wf/" + lang + "/" + platform + "/alerts"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		missiontype, _ := jsonparser.GetString(value, "mission", "type")
		missionfaction, _ := jsonparser.GetString(value, "mission", "faction")
		missionlocation, _ := jsonparser.GetString(value, "mission", "node")
		minEnemyLevel, _ := jsonparser.GetInt(value, "mission", "minEnemyLevel")
		maxEnemyLevel, _ := jsonparser.GetInt(value, "mission", "maxEnemyLevel")
		enemywaves, _ := jsonparser.GetInt(value, "mission", "maxWaveNum")
		rewardcredits, _ := jsonparser.GetInt(value, "mission", "reward", "credits")
		rewarditemsmany, _ := jsonparser.GetString(value, "mission", "reward", "countedItems", "[0]", "type")
		rewarditemsmanycount, _ := jsonparser.GetInt(value, "mission", "reward", "countedItems", "[0]", "count")
		rewarditem, _ := jsonparser.GetString(value, "mission", "reward", "items", "[0]")

		w := Alerts{id, started,
			ended, missiontype,
			missionfaction, missionlocation,
			minEnemyLevel, maxEnemyLevel, enemywaves,
			rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
		alerts = append(alerts, w)

	}, "Alerts")

	topicf := "/wf/" + lang + "/" + platform + "/alerts"
	messageJSON, _ := json.Marshal(alerts)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

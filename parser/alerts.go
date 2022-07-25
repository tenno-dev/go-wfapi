package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// Alerts struct
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

// Alertsdata export Alertsdata
var Alertsdata = make(map[int]map[string][]Alerts)

// ParseAlerts parsing Alerts data
func ParseAlerts(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Alertsdata[platformno]; !ok {
		Alertsdata[platformno] = make(map[string][]Alerts)
	}
	data := datasources.Apidata[platformno]
	var alerts []Alerts
	_, _, _, erralert := jsonparser.Get(data, "Alerts")
	// fmt.Println(erralert)
	if erralert != nil {
		return
	}
	// fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		missiontype, _ := jsonparser.GetString(value, "MissionInfo", "missionType")
		missiontype = helper.Missiontranslate(missiontype, lang)
		missionfaction, _ := jsonparser.GetString(value, "MissionInfo", "faction")
		missionfaction = helper.Factionstranslate(missionfaction, lang)
		missionlocation, _ := jsonparser.GetString(value, "MissionInfo", "location")
		missionlocation = helper.Regiontranslate(missionlocation, lang)
		minEnemyLevel, _ := jsonparser.GetInt(value, "MissionInfo", "minEnemyLevel")
		maxEnemyLevel, _ := jsonparser.GetInt(value, "MissionInfo", "maxEnemyLevel")
		enemywaves, _ := jsonparser.GetInt(value, "MissionInfo", "maxWaveNum")
		rewardcredits, _ := jsonparser.GetInt(value, "MissionInfo", "missionReward", "credits")
		rewarditemsmany, _ := jsonparser.GetString(value, "MissionInfo", "missionReward", "countedItems", "[0]", "ItemType")
		rewarditemsmany = helper.Langtranslate1(rewarditemsmany, lang)
		rewarditemsmanycount, _ := jsonparser.GetInt(value, "MissionInfo", "missionReward", "countedItems", "[0]", "ItemCount")
		rewarditem, _ := jsonparser.GetString(value, "MissionInfo", "missionReward", "items", "[0]")
		rewarditem = helper.Langtranslate1(rewarditem, lang)

		w := Alerts{id, started,
			ended, missiontype,
			missionfaction, missionlocation,
			minEnemyLevel, maxEnemyLevel, enemywaves,
			rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
		alerts = append(alerts, w)

	}, "Alerts")
	Alertsdata[platformno][lang] = alerts
}

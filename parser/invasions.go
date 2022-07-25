package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// Invasion struct
type Invasion struct {
	ID                  string
	Location            string
	MissionType         string
	Completed           bool
	Started             string
	AttackerRewardItem  string `json:",omitempty"`
	AttackerRewardCount int64  `json:",omitempty"`
	AttackerMissionInfo string `json:",omitempty"`
	DefenderRewardItem  string `json:",omitempty"`
	DefenderRewardCount int64  `json:",omitempty"`
	DefenderMissionInfo string `json:",omitempty"`
	Completion          float32
}

// Invasiondata export Invasiondata
var Invasiondata = make(map[int]map[string][]Invasion)

// ParseInvasions parse active Invasions
func ParseInvasions(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Invasiondata[platformno]; !ok {
		Invasiondata[platformno] = make(map[string][]Invasion)
	}

	data := datasources.Apidata[platformno]
	invasioncheck, _, _, _ := jsonparser.Get(data, "Invasions")
	if len(invasioncheck) == 0 {
		return
	}
	var invasions []Invasion
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		iscomplete, _ := jsonparser.GetBoolean(value, "Completed")
		if !iscomplete {
			attackeritem := ""
			attackeritemcount := int64(0)
			defenderitem := ""
			defenderitemcount := int64(0)
			id, _ := jsonparser.GetString(value, "_id", "$oid")
			started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
			location1, _ := jsonparser.GetString(value, "Node")
			location := helper.Regiontranslate(location1, lang)
			missiontype, _ := jsonparser.GetString(value, "LocTag")
			missiontype = helper.Langtranslate1(missiontype, lang)
			completed, _ := jsonparser.GetBoolean(value, "Completed")
			_, _, _, ierror := jsonparser.Get(value, "AttackerReward", "countedItems", "[0]", "ItemType")
			if ierror == nil {
				attackeritem, _ = jsonparser.GetString(value, "AttackerReward", "countedItems", "[0]", "ItemType")
				attackeritem = helper.Langtranslate1(attackeritem, lang)
				attackeritemcount, _ = jsonparser.GetInt(value, "AttackerReward", "countedItems", "[0]", "ItemCount")
			}
			attackerfaction, _ := jsonparser.GetString(value, "AttackerMissionInfo", "faction")
			attackerfaction = helper.Factionstranslate(attackerfaction, lang)
			_, _, _, ierror2 := jsonparser.Get(value, "DefenderReward", "countedItems", "[0]", "ItemType")
			if ierror2 == nil {
				defenderitem, _ = jsonparser.GetString(value, "DefenderReward", "countedItems", "[0]", "ItemType")
				defenderitem = helper.Langtranslate1(defenderitem, lang)
				defenderitemcount, _ = jsonparser.GetInt(value, "DefenderReward", "countedItems", "[0]", "ItemCount")
			}

			defenderfaction, _ := jsonparser.GetString(value, "DefenderMissionInfo", "faction")
			defenderfaction = helper.Factionstranslate(defenderfaction, lang)
			completion1, _ := jsonparser.GetInt(value, "Count")
			completion2, _ := jsonparser.GetInt(value, "Goal")
			completion := calcCompletion(completion1, completion2, attackerfaction)
			w := Invasion{id, location, missiontype, completed, started,
				attackeritem, attackeritemcount, attackerfaction,
				defenderitem, defenderitemcount, defenderfaction, completion}
			invasions = append(invasions, w)
		}
	}, "Invasions")
	Invasiondata[platformno][lang] = invasions
}
func calcCompletion(count int64, goal int64, attacker string) (complete float32) {
	y := float32((1 + float32(count)/float32(goal)))
	x := float32(y * 50)
	if attacker == "Infested" || attacker == "FC_INFESTATION" {
		x = float32(y * 100)
	}
	return x
}

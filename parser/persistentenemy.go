package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// Penemy struct
type Penemy struct {
	ID              string
	Health          float64
	FleeDamage      int64
	Rank            int64
	Region          int64
	MissionLocation string
	MissionType     string
	MissionFaction  string
	Lasttime        string
	Enemy           string
	Discovered      bool
	UseTicketing    bool
}

// Penemydata export Penemydata
var Penemydata = make(map[int]map[string][]Penemy)

// ParsePenemy parsing  persistent enemy data
func ParsePenemy(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Penemydata[platformno]; !ok {
		Penemydata[platformno] = make(map[string][]Penemy)
	}
	data := datasources.Apidata[platformno]
	var penemy []Penemy
	_, _, _, erralert := jsonparser.Get(data, "PersistentEnemies")
	// fmt.Println(erralert)
	if erralert != nil {
		return
	}
	// fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		health, _ := jsonparser.GetFloat(value, "HealthPercent")
		fleedamage, _ := jsonparser.GetInt(value, "FleeDamage")
		rank, _ := jsonparser.GetInt(value, "Rank")
		region, _ := jsonparser.GetInt(value, "Region")
		lastlocation1, _ := jsonparser.GetString(value, "LastDiscoveredLocation")
		lastlocation2 := helper.Regiontranslate(lastlocation1, lang)
		lastlocation := lastlocation2[0]
		lasttime, _ := jsonparser.GetString(value, "LastDiscoveredTime", "$date", "$numberLong")

		missiontype, _ := jsonparser.GetString(value, "MissionInfo", "missionType")
		missiontype = helper.Missiontranslate(missiontype, lang)
		missionfaction, _ := jsonparser.GetString(value, "MissionInfo", "faction")
		missionfaction = helper.Factionstranslate(missionfaction, lang)
		enemy1, _ := jsonparser.GetString(value, "AgentType")
		enemy := helper.Langtranslate1(enemy1, lang)
		discovered, _ := jsonparser.GetBoolean(value, "Discovered")
		ticketing, _ := jsonparser.GetBoolean(value, "UseTicketing")
		w := Penemy{id, health,
			fleedamage, rank,
			region, lastlocation, missiontype, missionfaction, lasttime,
			enemy, discovered, ticketing}
		penemy = append(penemy, w)

	}, "PersistentEnemies")
	Penemydata[platformno][lang] = penemy
}

package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// Fissures struct
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

// Fissuresdata export Fissuresdata
var Fissuresdata = make(map[int]map[string][]Fissures)

// ParseFissures  parse Fissure data
func ParseFissures(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Fissuresdata[platformno]; !ok {
		Fissuresdata[platformno] = make(map[string][]Fissures)
	}
	data := datasources.Apidata[platformno]
	var fissures []Fissures
	// fmt.Println("Fissues  reached")
	_, _, _, errfissures := jsonparser.Get(data, "ActiveMissions")
	if errfissures != nil {
		return
	}
	// fmt.Println("Fissues 2 reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		active := true
		location1, _ := jsonparser.GetString(value, "Node")
		location := helper.Regiontranslate(location1, lang)
		missiontype1, _ := jsonparser.GetString(value, "MissionType")
		missiontype := helper.Missiontranslate(missiontype1, lang)
		missionenemy := helper.Enemytranslate(location1, lang)
		tier1, _ := jsonparser.GetString(value, "Modifier")
		tier := helper.Voidtranslate(tier1, lang)
		expired, _ := jsonparser.GetBoolean(value, "expired")

		w := Fissures{id, started, ended, active,
			missiontype, missionenemy, location, tier[0], tier[1],
			expired}
		fissures = append(fissures, w)
	}, "ActiveMissions")
	Fissuresdata[platformno][lang] = fissures
}

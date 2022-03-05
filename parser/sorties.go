package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

type Sortievariant struct {
	MissionType     string
	MissionMod      string
	MissionModDesc  string
	MissionLocation string
}
type Sortie struct {
	ID       string
	Started  string
	Ends     string
	Boss     string
	Faction  string
	Reward   string
	Variants []Sortievariant
	Active   bool
}

var Sortiedata = make(map[int]map[string][]Sortie)

// ParseSorties parsing Sorties data
func ParseSorties(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := Sortiedata[platformno]; !ok {
		Sortiedata[platformno] = make(map[string][]Sortie)
	}
	// fmt.Println("reached sortie start")
	data := datasources.Apidata[platformno]
	_, _, _, sortieerr := jsonparser.Get(data, "Sorties")
	if sortieerr != nil {
		return
	}
	// fmt.Println("reached sortie start")

	var sortie []Sortie
	id, _ := jsonparser.GetString(data, "Sorties", "[0]", "_id", "$oid")
	started, _ := jsonparser.GetString(data, "Sorties", "[0]", "Activation", "$date", "$numberLong")
	ended, _ := jsonparser.GetString(data, "Sorties", "[0]", "Expiry", "$date", "$numberLong")
	boss1, _ := jsonparser.GetString(data, "Sorties", "[0]", "Boss")
	boss := helper.Sortietranslate(boss1, "sortiemodboss", lang)
	reward, _ := jsonparser.GetString(data, "Sorties", "[0]", "Reward")
	reward1 := helper.Sortietranslate2(reward, lang)
	reward = string(reward1[:]) // temp
	var variants []Sortievariant

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		mtype1, _ := jsonparser.GetString(value, "missionType")
		mtype := helper.Missiontranslate(mtype1, lang)
		mmod1, _ := jsonparser.GetString(value, "modifierType")
		mmod := helper.Sortietranslate(mmod1, "sortiemod", lang)
		mloc1, _ := jsonparser.GetString(value, "node")
		mloc2 := helper.Regiontranslate(mloc1, lang)
		mloc := mloc2[0] + " (" + mloc2[1] + ")"
		variants = append(variants, Sortievariant{
			MissionType:     mtype,
			MissionMod:      mmod[0],
			MissionModDesc:  mmod[1],
			MissionLocation: mloc,
		})
	}, "Sorties", "[0]", "Variants")
	active := true
	w := Sortie{ID: id, Started: started,
		Ends: ended, Boss: boss[1], Faction: boss[0], Reward: reward, Variants: variants,
		Active: active}
	sortie = append(sortie, w)
	Sortiedata[platformno][lang] = sortie
}

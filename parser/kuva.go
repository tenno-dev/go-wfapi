package parser

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
)

// KuvaData struct
type KuvaData struct {
	ID          string
	Start       string
	Ends        string
	Node        string
	Node2       string
	Planet      string
	Missiontype string
	Enemy       string
	Archwing    bool
	Sharkwing   bool
}

// ArbitrationData struct
type ArbitrationData struct {
	ID          string
	Start       string
	Ends        string
	Node        string
	Node2       string
	Planet      string
	Missiontype string
	Enemy       string
	Archwing    bool
	Sharkwing   bool
}

const timeFormat = "2006-01-02T15:04:05.000Z"

// KuvaMission for http export
var KuvaMission = make(map[int]map[string][]KuvaData)

// ArbitrationMission for http export
var ArbitrationMission = make(map[int]map[string][]ArbitrationData)

// ParseKuva Parse current Darvo Deal
func ParseKuva(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := KuvaMission[platformno]; !ok {
		KuvaMission[platformno] = make(map[string][]KuvaData)
	}
	if _, ok := ArbitrationMission[platformno]; !ok {
		ArbitrationMission[platformno] = make(map[string][]ArbitrationData)
	}
	data := datasources.Kuvadata[:]
	fmt.Println(string(data))
	var kuva []KuvaData
	var arbi []ArbitrationData

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//id, _ := jsonparser.GetString(value, "id")
		id := "1"
		now := time.Now().UTC()
		started, _ := jsonparser.GetString(value, "start")
		ended, _ := jsonparser.GetString(value, "end")
		start, _ := time.Parse(timeFormat, started)
		end, _ := time.Parse(timeFormat, ended)
		if start.Before(now) && end.After(now) {
			sdate := started
			edate := ended
			mtyperaw, _ := jsonparser.GetString(value, "missiontype")
			arch, _ := jsonparser.GetBoolean(value, "archwing")
			shark, _ := jsonparser.GetBoolean(value, "sharkwing")

			if mtyperaw == "EliteAlertMission" {
				a := ArbitrationData{id, sdate, edate, "node", "node2", "planet", "mtype", "enemy",
					arch, shark}
				arbi = append(arbi, a)
			} else {
				k := KuvaData{id, sdate, edate, "node", "node2", "planet", "mtype", "enemy",
					arch, shark}
				kuva = append(kuva, k)
			}
		}

	})
	ArbitrationMission[platformno][lang] = arbi
	KuvaMission[platformno][lang] = kuva
}

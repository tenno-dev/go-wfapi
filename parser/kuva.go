package parser

import (
	"encoding/json"
	"time"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// KuvaMission struct
type KuvaData struct {
	ID          string
	Start       string
	Ends        string
	Node        string
	Planet      string
	Missiontype string
	Enemy       string
	Archwing    bool
	Sharkwing   bool
}

// ArbitrationMission struct
type ArbitrationData struct {
	ID          string
	Start       string
	Ends        string
	Node        string
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

// ParseDarvoDeal Parse current Darvo Deal
func ParseKuva(platformno int, platform string, c mqtt.Client, lang string) {
	if _, ok := KuvaMission[platformno]; !ok {
		KuvaMission[platformno] = make(map[string][]KuvaData)
	}
	if _, ok := ArbitrationMission[platformno]; !ok {
		ArbitrationMission[platformno] = make(map[string][]ArbitrationData)
	}
	data := datasources.Kuvadata[:]
	var kuva []KuvaData
	var arbi []ArbitrationData

	// fmt.Println("Darvo  reached")
	/*errfissures, _ := jsonparser.GetString(data, "DailyDeals")
	if errfissures != "" {
		topicf := "wf/" + lang + "/" + platform + "/darvodeals"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		// fmt.Println("error Darvo reached")
		return
	}*/
	// fmt.Println("Darvo2 reached")
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

			mloc1, _ := jsonparser.GetString(value, "solnode")
			mloc2 := helper.Regiontranslate(mloc1, lang)
			mtype := mloc2[3]
			node := mloc2[0]
			planet := mloc2[1]
			enemy := mloc2[2]
			arch, _ := jsonparser.GetBoolean(value, "archwing")
			shark, _ := jsonparser.GetBoolean(value, "sharkwing")

			if mtyperaw == "EliteAlertMission" {
				a := ArbitrationData{id, sdate, edate, node, planet, mtype, enemy,
					arch, shark}
				arbi = append(arbi, a)
			} else {
				k := KuvaData{id, sdate, edate, node, planet, mtype, enemy,
					arch, shark}
				kuva = append(kuva, k)
			}
		}

	})

	topica := "wf/" + lang + "/" + platform + "/arbitration"
	messageJSONa, _ := json.Marshal(arbi)
	ArbitrationMission[platformno][lang] = arbi
	tokena := c.Publish(topica, 0, true, messageJSONa)
	topick := "wf/" + lang + "/" + platform + "/kuva"
	tokena.Wait()
	messageJSONk, _ := json.Marshal(kuva)
	KuvaMission[platformno][lang] = kuva
	tokenk := c.Publish(topick, 0, true, messageJSONk)
	tokenk.Wait()
}

package parser

import (
	"encoding/json"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DailyChallenges - DailyChallenges
type DailyChallenges struct {
	ID          string
	Ends        int64
	Started     int64
	Active      bool
	Reputation  int64
	Description string
	Title       string
}

// WeeklyChallenges - WeeklyChallenges
type WeeklyChallenges struct {
	ID          string
	Ends        int64
	Started     int64
	Active      bool
	Reputation  int64
	Description string
	Title       string
}

// WeeklyEliteChallenges - WeeklyEliteChallenges
type WeeklyEliteChallenges struct {
	ID          string
	Ends        int64
	Started     int64
	Active      bool
	Reputation  int64
	Description string
	Title       string
}

// Nightwave - Nightwave
type Nightwave struct {
	ID                    string
	Ends                  string
	Started               string
	Season                int64
	Tag                   string
	Phase                 int64
	params                string
	possibleChallenges    string
	DailyChallenges       []DailyChallenges
	WeeklyChallenges      []WeeklyChallenges
	WeeklyEliteChallenges []WeeklyEliteChallenges
}

// Nightwavedata export Alertsdata
var Nightwavedata = make(map[int]map[string][]Nightwave)

// ParseNightwave Parse Nightwave Season Info
func ParseNightwave(platformno int, platform string, c mqtt.Client, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Nightwavedata[platformno]; !ok {
		Nightwavedata[platformno] = make(map[string][]Nightwave)
	}
	data := datasources.Apidata[platformno]
	var nightwave []Nightwave
	var dchallenge []DailyChallenges
	var wchallenge []WeeklyChallenges
	var welitechallenge []WeeklyEliteChallenges
	errfissures, _ := jsonparser.GetString(data, "SeasonInfo")
	if errfissures != "" {
		topicf := "wf/" + lang + "/" + platform + "/nightwave"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		// fmt.Println("error Nightwave reached")
		return
	}
	timenow := time.Now().Unix()
	// fmt.Println("nightwave reached")
	id := "1"
	ended, _ := jsonparser.GetString(data, "SeasonInfo", "Expiry", "$date", "$numberLong")
	activation, _ := jsonparser.GetString(data, "SeasonInfo", "Activation", "$date", "$numberLong")
	season, _ := jsonparser.GetInt(data, "SeasonInfo", "Season")
	tag1, _ := jsonparser.GetString(data, "SeasonInfo", "AffiliationTag")
	phase, _ := jsonparser.GetInt(data, "SeasonInfo", "Phase")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		idc, _ := jsonparser.GetString(value, "_id", "$oid")
		endedc, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		endedc1, _ := strconv.ParseInt(endedc, 10, 64)
		endedc2 := endedc1 / 1000
		// fmt.Println(idc)
		// fmt.Println(endedc2)
		activationc, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		activationc1, _ := strconv.ParseInt(activationc, 10, 64)
		activationc2 := activationc1 / 1000

		mission, _ := jsonparser.GetString(value, "Challenge")
		cdesc := helper.Langtranslate2(mission, lang)
		reputation, _ := jsonparser.GetInt(value, "reputation")
		active := true // temp
		if endedc2 > timenow {
			active = true // temp
		}
		daily, _ := jsonparser.GetBoolean(value, "Daily")
		elite, _ := regexp.MatchString(`WeeklyHard.*`, mission)
		if daily == true {
			dailyc := DailyChallenges{idc, endedc2, activationc2, active, reputation, cdesc[1], cdesc[0]}
			dchallenge = append(dchallenge, dailyc)
		}
		if daily == false && elite == false {
			weekc := WeeklyChallenges{idc, endedc2, activationc2, active, reputation, cdesc[1], cdesc[0]}
			wchallenge = append(wchallenge, weekc)
		}
		if daily == false && elite == true {
			weekelitec := WeeklyEliteChallenges{idc, endedc2, activationc2, active, reputation, cdesc[1], cdesc[0]}
			welitechallenge = append(welitechallenge, weekelitec)
		}
	}, "SeasonInfo", "ActiveChallenges")
	w := Nightwave{id, ended, activation, season + 1, tag1,
		phase, "", "", dchallenge, wchallenge, welitechallenge}
	nightwave = append(nightwave, w)
	topicf := "wf/" + lang + "/" + platform + "/nightwave"
	Nightwavedata[platformno][lang] = nightwave
	messageJSON, _ := json.Marshal(nightwave)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

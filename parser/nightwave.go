package parser

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ParseNightwave Parse Nightwave Season Info
func ParseNightwave(platformno int, platform string, c mqtt.Client, lang string) {

	type DailyChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
	type WeeklyChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
	type WeeklyEliteChallenges struct {
		ID          string
		Ends        string
		Started     string
		Active      bool
		Reputation  int64
		Description string
		Title       string
	}
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
	data := datasources.Apidata[platformno]
	var nightwave []Nightwave
	var dchallenge []DailyChallenges
	var wchallenge []WeeklyChallenges
	var welitechallenge []WeeklyEliteChallenges

	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "SeasonInfo")
	if errfissures != "" {
		topicf := "/wf/" + lang + "/" + platform + "/nightwave"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Nightwave reached")
		return
	}
	fmt.Println("nightwave reached")
	id := "1"
	ended, _ := jsonparser.GetString(data, "SeasonInfo", "Expiry", "$date", "$numberLong")
	acvtivation, _ := jsonparser.GetString(data, "SeasonInfo", "Activation", "$date", "$numberLong")
	season, _ := jsonparser.GetInt(data, "SeasonInfo", "Season")
	tag1, _ := jsonparser.GetString(data, "SeasonInfo", "AffiliationTag")
	phase, _ := jsonparser.GetInt(data, "SeasonInfo", "Phase")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		idc, _ := jsonparser.GetString(value, "_id", "$oid")
		endedc, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		activationc, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		mission, _ := jsonparser.GetString(value, "Challenge")
		cdesc := helper.Langtranslate2(mission, lang)
		ctitle, _ := jsonparser.GetString(value, "title")
		reputation, _ := jsonparser.GetInt(value, "reputation")
		active, _ := jsonparser.GetBoolean(value, "active")
		daily, _ := jsonparser.GetBoolean(value, "Daily")
		elite, _ := regexp.MatchString(`WeeklyHard.*`, mission)
		if daily == true {
			dailyc := DailyChallenges{idc, endedc, activationc, active, reputation, cdesc[1], ctitle}
			dchallenge = append(dchallenge, dailyc)
		}
		if daily == false && elite == false {
			weekc := WeeklyChallenges{idc, endedc, activationc, active, reputation, cdesc[1], ctitle}
			wchallenge = append(wchallenge, weekc)
		}
		if daily == false && elite == true {
			weekelitec := WeeklyEliteChallenges{idc, endedc, activationc, active, reputation, cdesc[1], ctitle}
			welitechallenge = append(welitechallenge, weekelitec)
		}
	}, "SeasonInfo", "ActiveChallenges")
	w := Nightwave{id, ended, acvtivation, season + 1, tag1,
		phase, "", "", dchallenge, wchallenge, welitechallenge}
	nightwave = append(nightwave, w)
	topicf := "/wf/" + lang + "/" + platform + "/nightwave"
	messageJSON, _ := json.Marshal(nightwave)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

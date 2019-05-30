package parser

import (
	"encoding/json"
	"strings"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ParseSyndicateMissions Parse Ostrons & Solaris United Missions
func ParseSyndicateMissions(platformno int, platform string, c mqtt.Client, lang string) {
	type SyndicateJobs struct {
		Jobtype        string
		Rewards        []string
		MinEnemyLevel  int64
		MaxEnemyLevel  int64
		StandingReward []string
	}
	type SyndicateMissions struct {
		ID        string
		Started   string
		End       string
		Syndicate string
		Jobs      []SyndicateJobs
	}
	data := datasources.Apidata[platformno]
	var syndicates []SyndicateMissions
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		syndicatecheck, _ := jsonparser.GetString(value, "Tag")
		if syndicatecheck == "CetusSyndicate" || syndicatecheck == "SolarisSyndicate" {
			id, _ := jsonparser.GetString(value, "_id", "$oid")
			started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
			ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
			syndicate, _ := jsonparser.GetString(value, "Tag")
			var jobs []SyndicateJobs
			jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
				jobtype, _ := jsonparser.GetString(value1, "jobType")
				rewards0, _ := jsonparser.GetString(value1, "rewards")
				rewards1 := helper.Langtranslate1(rewards0, lang)
				rewards := strings.Split(rewards1, ",")
				minEnemyLevel, _ := jsonparser.GetInt(value1, "minEnemyLevel")
				maxEnemyLevel, _ := jsonparser.GetInt(value1, "maxEnemyLevel")
				standing := make([]string, 0)
				jsonparser.ArrayEach(value1, func(xpam []byte, dataType jsonparser.ValueType, offset int, err error) {
					standing = append(standing, string(xpam))

				}, "xpAmounts")
				jobs = append(jobs, SyndicateJobs{
					Jobtype:        jobtype,
					Rewards:        rewards,
					MinEnemyLevel:  minEnemyLevel,
					MaxEnemyLevel:  maxEnemyLevel,
					StandingReward: standing,
				})
			}, "Jobs")

			w := SyndicateMissions{
				ID:        id,
				Started:   started,
				End:       ended,
				Syndicate: syndicate,
				Jobs:      jobs}
			syndicates = append(syndicates, w)
		}
	}, "SyndicateMissions")

	topicf := "wf/" + lang + "/" + platform + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

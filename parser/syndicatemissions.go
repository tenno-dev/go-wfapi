package parser

import (
	"encoding/json"

	"github.com/bitti09/go-wfapi/datasources"
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
		StandingReward []int64
	}
	type SyndicateMissions struct {
		ID        string
		Started   string
		Ends      string
		Syndicate string
		Jobs      []SyndicateJobs
	}
	data := datasources.Apidata[platformno]
	var syndicates []SyndicateMissions
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		syndicatecheck, _ := jsonparser.GetString(value, "Tag")
		if syndicatecheck != "CetusSyndicate" && syndicatecheck != "SolarisSyndicate" {
			return
		}
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		syndicate, _ := jsonparser.GetString(value, "syndicate")
		var jobs []SyndicateJobs
		jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			jobtype, _ := jsonparser.GetString(value1, "type")
			rewards := make([]string, 0)
			jsonparser.ArrayEach(value1, func(reward []byte, dataType jsonparser.ValueType, offset int, err error) {
				rewards = append(rewards, string(reward))

			}, "rewardPool")

			minEnemyLevel, _ := jsonparser.GetInt(value1, "enemyLevels", "[0]")
			maxEnemyLevel, _ := jsonparser.GetInt(value1, "enemyLevels", "[1]")
			standing1, _ := jsonparser.GetInt(value1, "standingStages", "[0]")
			standing2, _ := jsonparser.GetInt(value1, "standingStages", "[1]")
			standing3, _ := jsonparser.GetInt(value1, "standingStages", "[2]")
			jobs = append(jobs, SyndicateJobs{
				Jobtype:        jobtype,
				Rewards:        rewards,
				MinEnemyLevel:  minEnemyLevel,
				MaxEnemyLevel:  maxEnemyLevel,
				StandingReward: []int64{standing1, standing2, standing3},
			})
		}, "jobs")

		w := SyndicateMissions{
			ID:        id,
			Started:   started,
			Ends:      ended,
			Syndicate: syndicate,
			Jobs:      jobs}
		syndicates = append(syndicates, w)
	}, "SyndicateMissions")

	topicf := "/wf/" + lang + "/" + platform + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}

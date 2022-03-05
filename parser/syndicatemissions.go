package parser

import (
	"strings"
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// SyndicateJobs struct
type SyndicateJobs struct {
	Jobtype        string
	Rewards        []string
	MinEnemyLevel  int64
	MaxEnemyLevel  int64
	StandingReward []string
}

// SyndicateMissions struct
type SyndicateMissions struct {
	ID        string
	Started   string
	End       string
	Syndicate string
	Jobs      []SyndicateJobs
}

// SyndicateMissionsdata export SyndicateMissions
var SyndicateMissionsdata = make(map[int]map[string][]SyndicateMissions)

// ParseSyndicateMissions Parse Ostrons & Solaris United Missions
func ParseSyndicateMissions(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := SyndicateMissionsdata[platformno]; !ok {
		SyndicateMissionsdata[platformno] = make(map[string][]SyndicateMissions)
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
				jobtype = helper.Langtranslate1(jobtype, lang)
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
	SyndicateMissionsdata[platformno][lang] = syndicates

}

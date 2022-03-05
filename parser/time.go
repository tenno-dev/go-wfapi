package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
)

// Time1 - Time Base
type Time1 struct {
	Cetus  []Time2
	Vallis []Time2
	Earth  []Time2
}

// Time2 - Time Details
type Time2 struct {
	Start string `json:",omitempty"`
	End   string
	State string
}

var Time1sdata = make(map[int]map[string][]Time1)

// ParseTime Parse Void trader
func ParseTime(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := Time1sdata[platformno]; !ok {
		Time1sdata[platformno] = make(map[string][]Time1)
	}
	data1 := datasources.Cetustime
	data2 := datasources.Valistime
	data3 := datasources.Earthtime

	var time1 []Time1
	var cetus []Time2
	var valis []Time2
	var earth []Time2

	cetusbegin, _ := jsonparser.GetString(data1, "activation")
	cetusend, _ := jsonparser.GetString(data1, "expiry")
	cetusstate, _ := jsonparser.GetString(data1, "state")
	cetus = append(cetus, Time2{
		Start: cetusbegin,
		End:   cetusend,
		State: cetusstate,
	})

	valisbegin, _ := jsonparser.GetString(data2, "activation")
	valisend, _ := jsonparser.GetString(data2, "expiry")
	valisstate, _ := jsonparser.GetString(data2, "state")
	valis = append(valis, Time2{
		Start: valisbegin,
		End:   valisend,
		State: valisstate,
	})

	earthbegin, _ := jsonparser.GetString(data3, "activation")
	earthend, _ := jsonparser.GetString(data3, "expiry")
	earthstate, _ := jsonparser.GetString(data3, "state")
	earth = append(earth, Time2{
		Start: earthbegin,
		End:   earthend,
		State: earthstate,
	})

	w := Time1{Cetus: cetus, Vallis: valis,
		Earth: earth}
	time1 = append(time1, w)
	Time1sdata[platformno][lang] = time1
}

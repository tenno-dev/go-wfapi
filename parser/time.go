package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/tidwall/gjson"
)

const nightTime = 3000

// Time1 - Time Base
type Time1 struct {
	Cetus   Time2
	Vallis  Time2
	Cambion Time2
}

// Time2 - Time Details
type Time2 struct {
	Start    string `json:",omitempty"`
	End      string
	State    string
	TimeLeft string
}

// Time1sdata export Time1
var Time1sdata = make(map[int]map[string]Time1)

// ParseTime Parse Void trader
func ParseTime(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := Time1sdata[platformno]; !ok {
		Time1sdata[platformno] = make(map[string]Time1)
	}
	data := datasources.CetusTimedata[platformno]
	data1 := datasources.VallisTimedata[platformno]
	data2 := datasources.CambionTimedata[platformno]

	// Cetus
	cetus := Time2{Start: gjson.Get(string(data), "activation").String(), End: gjson.Get(string(data), "expiry").String(), State: gjson.Get(string(data), "state").String(), TimeLeft: gjson.Get(string(data), "timeLeft").String()}

	// data1
	vallis := Time2{Start: gjson.Get(string(data1), "activation").String(), End: gjson.Get(string(data1), "expiry").String(), State: gjson.Get(string(data1), "state").String(), TimeLeft: gjson.Get(string(data1), "timeLeft").String()}

	// cambion
	cambion := Time2{Start: gjson.Get(string(data2), "activation").String(), End: gjson.Get(string(data2), "expiry").String(), State: gjson.Get(string(data2), "active").String(), TimeLeft: gjson.Get(string(data2), "timeLeft").String()}
	w := Time1{Cetus: cetus, Vallis: vallis, Cambion: cambion}
	Time1sdata[platformno][lang] = w
}

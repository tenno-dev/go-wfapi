package parser

import (
	"fmt"
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/nleeper/goment"
	"github.com/tidwall/gjson"
)

type Zarimancycle struct {
	ID              string
	bountiesEndDate string
	expiry          string
	activation      string
	isCorpus        bool
	state           string
	timeLeft        string
}

var Zariman = make(map[int]map[string][]Zarimancycle)

func ParseZariman(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	//corpusTimeMillis := 1655182800000
	//fullCycle := 18000000
	//stateMaximum := 9000000
	data := datasources.Apidata[platformno]
	//i := 0
	ZarimanBountyEnd := gjson.Get(string(data), "SyndicateMissions.#(Tag == \"ZarimanSyndicate\").Expiry.$date.$numberLong").Int()
	fmt.Println(platformno, platform, lang)
	fmt.Println(ZarimanBountyEnd)
	g, _ := goment.Unix(ZarimanBountyEnd / 1000)
	fmt.Println(g.ToNow())

}

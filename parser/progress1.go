package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
)

// Progress1 - Progress1
type Progress1 struct {
	P1 float64
	P2 float64
	P3 float64
}

var Progress1data = make(map[int]map[string][]Progress1)

// ParseProgress1 Parse Void trader
func ParseProgress1(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := Progress1data[platformno]; !ok {
		Progress1data[platformno] = make(map[string][]Progress1)
	}
	data := datasources.Apidata[platformno]
	var progress1 []Progress1

	_, _, _, pro1err := jsonparser.Get(data, "ProjectPct")
	if pro1err != nil {
		return
	}
	// fmt.Println("reached progress start")

	p1, _ := jsonparser.GetFloat(data, "ProjectPct", "[0]")
	p2, _ := jsonparser.GetFloat(data, "ProjectPct", "[1]")
	p3, _ := jsonparser.GetFloat(data, "ProjectPct", "[2]")

	w := Progress1{P1: p1, P2: p2,
		P3: p3}
	progress1 = append(progress1, w)
	Progress1data[platformno][lang] = progress1

}

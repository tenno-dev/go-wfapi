package outputs

import (
	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/gin-gonic/gin"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

// Everything test 2
func Everything(c *gin.Context) {
	v, _, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	header.Set("Content-Type", "application/json; charset=utf-8")
	test := datasources.Apidata[v]
	c.String(200, string(test[:]))
}

// DarvoDeals DarvoDeals
func DarvoDeals(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Darvodata[v][t1])
}

// News godoc
func News(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Newsdata[v][t1])
}

// Alerts Alertsdata
func Alerts(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Alertsdata[v][t1])
}

// Fissures Fissuresdata
func Fissures(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Fissuresdata[v][t1])
}

// Nightwave data
func Nightwave(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Nightwavedata[v][t1])
}

// Penemy sdata
func Penemy(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Penemydata[v][t1])
}

// Invasion sdata
func Invasion(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Invasiondata[v][t1])
}

// Time sdata
func Time(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Time1sdata[v][t1])
}

// Sortie sdata
func Sortie(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Sortiedata[v][t1])
}

// Voidtrader sdata
func Voidtrader(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Voidtraderdata[v][t1])
}

// SyndicateMission sdata
func SyndicateMission(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.SyndicateMissionsdata[v][t1])
}

// AnomalyData sdata
func AnomalyData(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.AnomalyDataSet[v][t1])
}

// Progress1 sdata
func Progress1(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Progress1data[v][t1])
}

// Event sdata
func Event(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.Eventdata[v][t1])
}

// ArbitrationMission sdata
func ArbitrationMission(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.ArbitrationMission[v][t1])
}

// KuvaMission sdata
func KuvaMission(c *gin.Context) {
	v, t1, t2 := getPlatformValueAndTokens(c)
	w := c.Writer
	header := w.Header()
	header.Set("Accept-Language", t2)
	c.JSON(200, parser.KuvaMission[v][t1])
}
func getPlatformValueAndTokens(c *gin.Context) (int, string, string) {
	platform := c.Params.ByName("platform")
	token := c.GetHeader("Accept-Language")
	value := intMap[platform]
	token1 := token[0:2]
	return value, token1, token
}

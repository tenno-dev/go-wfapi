package outputs

import (
	"net/http"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

type Test struct {
	Timestamp  string                     `json:"timestamp"`
	Darvo      []parser.DarvoDeals        `json:"darvo"`
	News       []parser.News              `json:"news"`
	Nightwave  []parser.Nightwave         `json:"nightwave"`
	Alerts     []parser.Alerts            `json:"alerts"`
	Penemydata parser.Progress1           `json:"progress"`
	Fissues    []parser.Fissures          `json:"fissures"`
	Time       parser.Time1               `json:"time"`
	Sortie     []parser.Sortie            `json:"sortie"`
	Voidtrader []parser.Voidtrader        `json:"voidtrader"`
	Syndicate  []parser.SyndicateMissions `json:"syndicate"`
	Invasion   []parser.Invasion          `json:"invasion"`
}

func Everything(w http.ResponseWriter, r *http.Request) {
	v, _, _ := getPlatformValueAndTokens(r)
	test := datasources.Apidata[v]
	w.WriteHeader(http.StatusOK)
	w.Write(test)
}

func Everything2(w http.ResponseWriter, r *http.Request) {
	v, t1, _ := getPlatformValueAndTokens(r)
	test := Test{datasources.Timestamp[v], parser.Darvodata[v][t1], parser.Newsdata[v][t1], parser.Nightwavedata[v][t1], parser.Alertsdata[v][t1], parser.Progress1data[v][t1], parser.Fissuresdata[v][t1], parser.Time1sdata[v][t1],
		parser.Sortiedata[v][t1], parser.Voidtraderdata[v][t1], parser.SyndicateMissionsdata[v][t1], parser.Invasiondata[v][t1]}
	render.JSON(w, r, test)
}

// DarvoDeals godoc
// @Summary      Show active  Darvo Deals
// @Description  get string by ID
// @Tags         Show DarvoDeals
// @Accept       json
// @Produce      json
// @Param        platform   path string  true  "Platform"
// @Param        lang    query     string  false  "lang select"
// @Success      200  {object}  parser.DarvoDeals
// @Router       /{platform}/darvo [get]
// DarvoDeals DarvoDeals
func DarvoDeals(w http.ResponseWriter, r *http.Request) {
	v, t1, _ := getPlatformValueAndTokens(r)
	render.JSON(w, r, parser.Darvodata[v][t1])
}

// Newsdata godoc
// @Summary      Show current News
// @Description  get string by ID
// @Tags         Show Newsdata
// @Accept       json
// @Produce      json
// @Param        platform   path string  true  "Platform"
// @Param        lang    query     string  false  "lang select"
// @Success      200  {object}  parser.News
// @Router       /{platform}/news [get]
// News godoc
func News(w http.ResponseWriter, r *http.Request) {
	v, t1, _ := getPlatformValueAndTokens(r)
	render.JSON(w, r, parser.Newsdata[v][t1])
}

// Alertsdata godoc
// @Summary      Show current Alerts
// @Description  get string by ID
// @Tags         Show Alertsdata
// @Accept       json
// @Produce      json
// @Param        platform   path string  true  "Platform"
// @Param        lang    query     string  false  "lang select"
// @Success      200  {object}  parser.Alerts
// @Router       /{platform}/alerts [get]
// Alerts Alertsdata
func Alerts(w http.ResponseWriter, r *http.Request) {
	v, t1, _ := getPlatformValueAndTokens(r)
	render.JSON(w, r, parser.Alertsdata[v][t1])
}

/*
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
*/
func getPlatformValueAndTokens(r *http.Request) (int, string, string) {
	platform := chi.URLParam(r, "platform")
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	value := intMap[platform]
	return value, lang, lang
}

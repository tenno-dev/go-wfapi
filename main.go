package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"github.com/buger/jsonparser"
//	"github.com/robfig/cron"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/profile"
	/**/
)

//current supported lang
var langid = map[string]int{
	"en": 0,
}
var langpool =  [10]string{"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"}
// lang end
// platforms start
var platforms = [4]string{"pc", "ps4", "xb1", "swi"}
// platforms end
var translationtype = [10]string{"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"} 
var bempty = "[{}]"
var langtest = "en"
// LangMap start
type LangMap map[string]interface{}
// LangMap2 d
type LangMap2  interface{}

var sortieloc = make(map[string]map[string]interface{})
var sortiemodtypes= make(map[string]map[string]interface{})
var sortiemoddesc= make(map[string]map[string]interface{})
var sortiemodbosses= make(map[string]map[string]interface{})
var sortielang = make(map[string]map[string]interface{})// temp
var sortielang1  map[string]interface{}// temp
var sortielang2  map[string]string// temp
var fissureModifiers= make(map[string]map[string]interface{})
var missionTypes  =make(map[string]map[string]interface{})

// todo
var arcanesData  map[string]interface{}
var conclaveData   map[string]interface{}
var eventsData  map[string]interface{}
var factionsData   map[string]interface{}
var languages =  map[string]string{}
var operationTypes =  map[string]string{}
var persistentEnemyData   map[string]interface{}

var syndicatesData =  map[string]string{}
var synthTargets   map[string]interface{}
var upgradeTypes   map[string]interface{}
var warframes   map[string]interface{}//
var weapons  map[string]interface{}
// langdata end
var apidata = make([][]byte, 4)
var sortierewards = ""

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func loadapidata(id1 string) (ret []byte) {
	// WF API Source
	client := &http.Client{}

	url := "http://content.warframe.com/dynamic/worldState.php"
	if id1 != "pc" {
		url = "http://content." + id1 + ".warframe.com/dynamic/worldState.php"
	}
		fmt.Println("url:", url)

	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	return body[:]	 
}
func loadlangdata(id1 string,id2 int) {
	
	client := &http.Client{}
	/*
	// arcanesData

	url := "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/arcanesData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/arcanesData.json"
	}
	// fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	arcanesData[id1] = string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// conclaveData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/conclaveData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/conclaveData.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	conclaveData[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// eventsData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/eventsData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/eventsData.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	eventsData[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// factionsData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/factionsData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/factionsData.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	factionsData[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)


	// languages
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/languages.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/languages.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	languages[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// missionTypes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/missionTypes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/missionTypes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	missionTypes[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// operationTypes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/operationTypes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/operationTypes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	operationTypes[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// persistentEnemyData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/persistentEnemyData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/persistentEnemyData.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	persistentEnemyData[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)
	*/
	// solNodes
	url := "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/solNodes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/solNodes.json"
	}
	fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(id1)
		var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
		sortieloc[id1] = result

	_, _ = io.Copy(ioutil.Discard, res.Body)
	
	// sortieData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/sortieData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/sortieData.json"
	}
	//fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &sortielang)
		fmt.Println("test2")
	//		sortielang = body.(map[string]interface{})
	sortiemodtypes[id1] = make(LangMap)
			fmt.Println("test3")

	 sortiemodtypes[id1] = sortielang["modifierTypes"]
	//	 fmt.Println("test2",sortiemodtypes[id1]["SORTIE_MODIFIER_ARMOR"])
	sortiemoddesc[id1] = sortielang["modifierDescriptions"]
	sortiemodbosses[id1] = sortielang["bosses"]
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// fissureModifiers
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/fissureModifiers.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/fissureModifiers.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &result)
	fissureModifiers[id1]= result
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// missionTypes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/missionTypes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/missionTypes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &result)
	missionTypes[id1]= result
	_, _ = io.Copy(ioutil.Discard, res.Body)


	/*
	// syndicatesData
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/syndicatesData.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/syndicatesData.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	syndicatesData[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// synthTargets
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/synthTargets.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/synthTargets.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	synthTargets[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)

	// upgradeTypes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/upgradeTypes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/upgradeTypes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	upgradeTypes[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)
	
	// warframes
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/warframes.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/warframes.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	warframes[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)
	
	// weapons
	url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/"+id1+"/weapons.json"
	if (id1 =="en"){
			url = "https://raw.githubusercontent.com/WFCD/warframe-worldstate-data/master/data/weapons.json"
	}
	// fmt.Println("url:", url)
	req, _ = http.NewRequest("GET", url, nil)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	weapons[id1]= string(body[:])
	_, _ = io.Copy(ioutil.Discard, res.Body)		
	*/
	return
}
func main() {
	defer profile.Start(profile.MemProfile).Stop()

	// mqtt client start
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:8884/mqtt").SetClientID("gotrivial2")
	//opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername("x")
	opts.SetPassword("x")
	//opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe("test/topic", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	//mqtt client end
			for x1, v1 := range langpool {
		fmt.Println("x1:", x1)
		fmt.Println("v1:", v1)
		loadlangdata(v1,x1)
			}/**/
	for x, v := range platforms {
		fmt.Println("x:", x)
		fmt.Println("v:", v)
		apidata[x] = loadapidata(v)
			for x1, v1 := range langpool {
		fmt.Println("x1:", x1)
		fmt.Println("v1:", v1)
		parseSorties(x, v, c,v1)
		parseNews(x, v, c,v1)
		parseAlerts(x, v, c,v1)
		parseFissures(x, v, c,v1)

			}
			//fmt.Println("test:", languageslang)
		/*
		parseAlerts(x, v, c)
		parseNews(x, v, c)
		parseSorties(x, v, c)
		parseSyndicateMissions(x, v, c)
		parseInvasions(x, v, c)
		parseCycles(x, v, c)
		parseFissures(x, v, c)
		parseDarvo(x, v, c)
		parseEvents(x, v, c)
		parseNightwave(x, v, c)
	*/
		PrintMemUsage()

	}
	PrintMemUsage()
	/*
	c1 := cron.New()
	c1.AddFunc("@every 1m1s", func() {

		fmt.Println("Tick")
		for x, v := range platforms {
			fmt.Println("x:", x)
			fmt.Println("v:", v)
			apidata[x] = loadapidata(v)
			/*
			parseAlerts(x, v, c)
			parseNews(x, v, c)
			parseSorties(x, v, c)
			parseSyndicateMissions(x, v, c)
			parseInvasions(x, v, c)
			parseCycles(x, v, c)
			parseFissures(x, v, c)
			parseDarvo(x, v, c)
			parseEvents(x, v, c)
			parseNightwave(x, v, c)
			 
			PrintMemUsage()

		}
		
		/*
				parseActiveMissions(x, v, c)
				parseInvasions(x, v, c)

		
	})
	c1.Start()
	*/
	PrintMemUsage()

	// just for debuging - printing  full warframe api response
	 

	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}

}

func parseAlerts(platformno int, platform string, c mqtt.Client , lang string) {
	type Alerts struct {
		ID                  string
		Started             string
		Ends                string
		MissionType         string
		MissionFaction      string
		MissionLocation     string
		MinEnemyLevel       int64
		MaxEnemyLevel       int64
		EnemyWaves          int64 `json:",omitempty"`
		RewardCredits       int64
		RewardItemMany      string `json:",omitempty"`
		RewardItemManyCount int64  `json:",omitempty"`
		RewardItem          string `json:",omitempty"`
	}
	data := apidata[platformno]
	var alerts []Alerts
	_, _, _, erralert := jsonparser.Get(data, "Alerts")
	fmt.Println(erralert)
	if erralert != nil  || erralert == nil {  // disable  parsing until api returns data
		topicf := "/wf/" +  lang +"/"+ platform + "/alerts"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		missiontype, _ := jsonparser.GetString(value, "mission", "type")
		missionfaction, _ := jsonparser.GetString(value, "mission", "faction")
		missionlocation, _ := jsonparser.GetString(value, "mission", "node")
		minEnemyLevel, _ := jsonparser.GetInt(value, "mission", "minEnemyLevel")
		maxEnemyLevel, _ := jsonparser.GetInt(value, "mission", "maxEnemyLevel")
		enemywaves, _ := jsonparser.GetInt(value, "mission", "maxWaveNum")
		rewardcredits, _ := jsonparser.GetInt(value, "mission", "reward", "credits")
		rewarditemsmany, _ := jsonparser.GetString(value, "mission", "reward", "countedItems", "[0]", "type")
		rewarditemsmanycount, _ := jsonparser.GetInt(value, "mission", "reward", "countedItems", "[0]", "count")
		rewarditem, _ := jsonparser.GetString(value, "mission", "reward", "items", "[0]")

		w := Alerts{id, started,
			ended, missiontype,
			missionfaction, missionlocation,
			minEnemyLevel, maxEnemyLevel, enemywaves,
			rewardcredits, rewarditemsmany, rewarditemsmanycount, rewarditem}
		alerts = append(alerts, w)

	}, "Alerts")

	topicf := "/wf/" +  lang +"/"+ platform + "/alerts"
	messageJSON, _ := json.Marshal(alerts)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

func parseNews(platformno int, platform string, c mqtt.Client , lang string) {
	type News struct {
		ID       string
		Message  string
		URL      string
		Date     string
		priority bool
		Image    string
	}
	data := apidata[platformno]
	_, _, _, ernews := jsonparser.Get(data, "Events")
	if ernews != nil {
		fmt.Println("error ernews reached")
		return
	}
	var errnews2 bool
	var message string

	var news []News
		fmt.Println("news reached")


	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		message =""
					jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
				newstemp1,_ :=jsonparser.GetString(value1, "LanguageCode")

		if (newstemp1 == lang) {
			message, _ = jsonparser.GetString(value1, "Message")
			 	fmt.Println("news lang", newstemp1)
 	fmt.Println("news lang1", lang)
 	fmt.Println("news ", message)

		}
			}, "Messages")	

	if (message !=""){
	errnews2 = false
	  	fmt.Println("errnews2 ", message)

	} else {
	errnews2 = true
	}


	if errnews2 == false {
				image := "http://n9e5v4d8.ssl.hwcdn.net/uploads/e0b4d18d3330bb0e62dcdcb364d5f004.png"
		id, _ := jsonparser.GetString(value, "_id", "$oid")

	
		url, _ := jsonparser.GetString(value, "Prop")
		image, _ = jsonparser.GetString(value, "ImageUrl")

		if (strings.HasPrefix(image, "https://forums.warframe.com")) {
			image =strings.Split(image, "=")[1]
					image = strings.Split(image, "&key")[0]

		}
		date, _ := jsonparser.GetString(value, "Date","$date","$numberLong")
		/**/
		priority, _ := jsonparser.GetBoolean(value, "priority")
		w := News{ID: id, Message: message, URL: url, Date: date, Image: image, priority: priority}
		news = append(news, w)
		topicf := "/wf/" + lang + "/"+ platform + "/news"
		messageJSON, _ := json.Marshal(news)
		token := c.Publish(topicf, 0, true, messageJSON)
		token.Wait()
	}
	}, "Events")
}
func parseCycles(platformno int, platform string, c mqtt.Client , lang string) {
	type Cycles struct {
		EathID         string
		EarthEnds      string
		EarthIsDay     bool
		EarthTimeleft  string
		CetusID        string
		CetusEnds      string
		CetusIsDay     bool
		CetusIsCetus   bool
		CetusTimeleft  string
		VallisID       string
		VallisEnds     string
		VallisIsWarm   bool
		VallisTimeleft string
	}
	data := apidata[platformno]
	var cycles []Cycles
	fmt.Println("Cycles reached")
	//  Earth
	earthid, _ := jsonparser.GetString(data, "earthCycle", "id")
	earthends, _ := jsonparser.GetString(data, "earthCycle", "expiry")
	earthisday, _ := jsonparser.GetBoolean(data, "earthCycle", "isDay")
	earthtimeleft, _ := jsonparser.GetString(data, "earthCycle", "timeLeft")
	// Cetus
	cetusid, _ := jsonparser.GetString(data, "cetusCycle", "id")
	cetusends, _ := jsonparser.GetString(data, "cetusCycle", "expiry")
	cetusisday, _ := jsonparser.GetBoolean(data, "cetusCycle", "isDay")
	cetusiscetus, _ := jsonparser.GetBoolean(data, "cetusCycle", "isCetus")
	cetustimeleft, _ := jsonparser.GetString(data, "cetusCycle", "timeLeft")
	// Vallis
	vallisid, _ := jsonparser.GetString(data, "vallisCycle", "id")
	vallisends, _ := jsonparser.GetString(data, "vallisCycle", "expiry")
	vallisiswarm, _ := jsonparser.GetBoolean(data, "vallisCycle", "isDay")
	vallistimeleft, _ := jsonparser.GetString(data, "vallisCycle", "timeLeft")

	w := Cycles{earthid, earthends, earthisday, earthtimeleft,
		cetusid, cetusends, cetusisday, cetusiscetus, cetustimeleft,
		vallisid, vallisends, vallisiswarm, vallistimeleft}
	cycles = append(cycles, w)

	topicf := "/wf/" + platform + "/"+ langtest + "/cycles"
	messageJSON, _ := json.Marshal(cycles)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseFissures(platformno int, platform string, c mqtt.Client , lang string) {
	type Fissures struct {
		ID              string
		Started         string
		Ends            string
		Active          bool
		MissionType     string
		MissionFaction  string
		MissionLocation string
		Tier            string
		TierLevel       string
		Expired         bool
	}
	data := apidata[platformno]
	var fissures []Fissures
	fmt.Println("Fissues  reached")
	_, _, _, errfissures := jsonparser.Get(data, "ActiveMissions")
	if errfissures != nil  {
		topicf := "/wf/" + lang + "/"+ platform + "/fissures"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error alert reached")
		return
	}
	fmt.Println("Fissues 2 reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "_id","$oid")
		started, _ := jsonparser.GetString(value, "Activation","$date","$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry","$date","$numberLong")
		active := true
		location1, _ := jsonparser.GetString(value, "Node")
		location := sortietranslate(location1,"sortieloc",lang)
		missiontype1, _ := jsonparser.GetString(value, "MissionType")
		missiontype := missionranslate(missiontype1,lang)
		tier1, _ := jsonparser.GetString(value, "Modifier")
		tier := voidranslate(tier1,lang)
		expired, _ := jsonparser.GetBoolean(value, "expired")

		w := Fissures{id, started, ended, active,
			missiontype, location[1], location[0], tier[0],tier[1],
			expired}
		fissures = append(fissures, w)
	}, "ActiveMissions")

	topicf := "/wf/" + lang + "/"+ platform + "/fissures"
	messageJSON, _ := json.Marshal(fissures)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseDarvo(platformno int, platform string, c mqtt.Client) {
	type DarvoDeals struct {
		ID              string
		Ends            string
		Item            string
		Price           int64
		DealPrice       int64
		DiscountPercent int64
		Stock           int64
		Sold            int64
	}
	data := apidata[platformno]
	var deals []DarvoDeals
	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "dailyDeals")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/"+ langtest + "/darvodeals"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Darvo reached")
		return
	}
	fmt.Println("alert reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		ended, _ := jsonparser.GetString(value, "expiry")
		item, _ := jsonparser.GetString(value, "item")
		originalprice, _ := jsonparser.GetInt(value, "originalPrice")
		dealprice, _ := jsonparser.GetInt(value, "salePrice")
		stock, _ := jsonparser.GetInt(value, "total")
		sold, _ := jsonparser.GetInt(value, "sold")
		discount, _ := jsonparser.GetInt(value, "discount")

		w := DarvoDeals{id, ended, item, originalprice, dealprice,
			discount, stock, sold}
		deals = append(deals, w)
	}, "dailyDeals")

	topicf := "/wf/" + platform + "/"+ langtest + "/darvodeals"
	messageJSON, _ := json.Marshal(deals)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseNightwave(platformno int, platform string, c mqtt.Client) {

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
	data := apidata[platformno]
	var nightwave []Nightwave
	var dchallenge []DailyChallenges
	var wchallenge []WeeklyChallenges
	var welitechallenge []WeeklyEliteChallenges

	fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "nightwave")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/"+ langtest + "/nightwave"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Nightwave reached")
		return
	}
	fmt.Println("nightwave reached")
	id, _ := jsonparser.GetString(data, "nightwave", "id")
	ended, _ := jsonparser.GetString(data, "nightwave", "expiry")
	acvtivation, _ := jsonparser.GetString(data, "nightwave", "activation")
	season, _ := jsonparser.GetInt(data, "nightwave", "season")
	tag1, _ := jsonparser.GetString(data, "nightwave", "tag")
	phase, _ := jsonparser.GetInt(data, "nightwave", "phase")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		idc, _ := jsonparser.GetString(value, "id")
		endedc, _ := jsonparser.GetString(value, "expiry")
		activationc, _ := jsonparser.GetString(value, "activation")
		cdesc, _ := jsonparser.GetString(value, "desc")
		ctitle, _ := jsonparser.GetString(value, "title")
		reputation, _ := jsonparser.GetInt(value, "reputation")
		active, _ := jsonparser.GetBoolean(value, "active")
		daily, _ := jsonparser.GetBoolean(value, "isDaily")
		elite, _ := jsonparser.GetBoolean(value, "isElite")
		if daily == true {
			dailyc := DailyChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			dchallenge = append(dchallenge, dailyc)
		}
		if daily == false && elite == false {
			weekc := WeeklyChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			wchallenge = append(wchallenge, weekc)
		}
		if daily == false && elite == true {
			weekelitec := WeeklyEliteChallenges{idc, endedc, activationc, active, reputation, cdesc, ctitle}
			welitechallenge = append(welitechallenge, weekelitec)
		}
	}, "nightwave", "activeChallenges")
	w := Nightwave{id, ended, acvtivation, season, tag1,
		phase, "", "", dchallenge, wchallenge, welitechallenge}
	nightwave = append(nightwave, w)
	topicf := "/wf/" + platform + "/"+ langtest + "/nightwave"
	messageJSON, _ := json.Marshal(nightwave)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseEvents(platformno int, platform string, c mqtt.Client) {
	type Events struct {
		ID              string
		Started         string
		Ends            string
		Active          bool
		MaxScore        int64
		CurrScore       int64
		Health          float64
		Health1         string
		Faction         string
		Description     string
		Tooltip         string
		ConcurrentNodes string //subject to change when api  has data for it
		Rewarditem      string
		Rewardcredits   int64
		interimSteps    string
		progressSteps   string //subject to change when api  has data for it
		PersonalEvent   bool
		CommunityEvent  bool
		Expired         bool
	}
	data := apidata[platformno]
	var events []Events

	fmt.Println("Events  reached")
	errfissures, _ := jsonparser.GetString(data, "Eveventsents")
	if errfissures != "" {
		topicf := "/wf/" + platform + "/"+ langtest + "/events"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("error Events reached")
		return
	}
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		rewarditems := ""
		rewardcredits := int64(0)
		health := float64(0)
		id, _ := jsonparser.GetString(value, "id")
		started, _ := jsonparser.GetString(value, "activation")
		ended, _ := jsonparser.GetString(value, "expiry")
		active, _ := jsonparser.GetBoolean(value, "active")
		maximumScore, _ := jsonparser.GetInt(value, "maximumScore")
		currentScore, _ := jsonparser.GetInt(value, "currentScore")
		health1, _ := jsonparser.GetString(value, "health")

		health, _ = strconv.ParseFloat(health1, 64)
		faction, _ := jsonparser.GetString(value, "faction")
		description, _ := jsonparser.GetString(value, "description")
		tooltip, _ := jsonparser.GetString(value, "tooltip")
		expired, _ := jsonparser.GetBoolean(value, "expired")
		rewarditems, _ = jsonparser.GetString(value, "rewards", "[0]", "items", "[0]")
		rewardcredits, _ = jsonparser.GetInt(value, "rewards", "[0]", "credits", "[0]")
		personal, _ := jsonparser.GetBoolean(value, "isPersonal")
		commu, _ := jsonparser.GetBoolean(value, "isCommunity")

		w := Events{id, started, ended, active, maximumScore, currentScore, health, health1, faction,
			description, tooltip, "", rewarditems, rewardcredits, "", "", personal, commu, expired}
		events = append(events, w)
	}, "events")

	topicf := "/wf/" + platform + "/"+ langtest + "/events"
	messageJSON, _ := json.Marshal(events)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
func parseSorties(platformno int, platform string, c mqtt.Client , lang string) {
	type Sortievariant struct {
		MissionType     string
		MissionMod      string
		MissionModDesc  string
		MissionLocation string
	}
	type Sortie struct {
		ID       string
		Started  string
		Ends     string
		Boss     string
		Faction  string
		Reward   string
		Variants []Sortievariant
		Active   bool
	}
	fmt.Println("reached sortie start")
	data := apidata[platformno]
	_,_,_, sortieerr := jsonparser.Get(data, "Sorties")
	if sortieerr != nil  {
		topicf := "/wf/" + lang + "/"+ platform + "/sorties"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("reached sortie error")

		return
	}
	fmt.Println("reached sortie start")

	var sortie []Sortie
	id, _ := jsonparser.GetString(data, "Sorties", "[0]","_id","$oid")
	started, _ := jsonparser.GetString(data, "Sorties", "[0]", "Activation","$date","$numberLong")
	ended, _ := jsonparser.GetString(data, "Sorties", "[0]", "Expiry","$date","$numberLong")
	boss1, _ := jsonparser.GetString(data, "Sorties", "[0]", "Boss")
	boss := sortietranslate(boss1,"sortiemodboss",lang)
	reward, _ := jsonparser.GetString(data, "Sorties", "[0]", "Reward")
	var variants []Sortievariant

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		mtype, _ := jsonparser.GetString(value, "missionType")
		mmod1, _ := jsonparser.GetString(value, "modifierType")
		mmod  := sortietranslate(mmod1,"sortiemod",lang)
		mloc1, _ := jsonparser.GetString(value, "node")
		mloc  := sortietranslate(mloc1,"sortieloc",lang)

		variants = append(variants, Sortievariant{
			MissionType:     mtype,
			MissionMod:      mmod[0],
			MissionModDesc:  mmod[1],
			MissionLocation: mloc[0],
		})
	}, "Sorties", "[0]", "Variants")
	active := true
	w := Sortie{ID: id, Started: started,
		Ends: ended, Boss: boss[1], Faction: boss[0], Reward: reward, Variants: variants,
		Active: active}
	sortie = append(sortie, w)

	topicf := "/wf/" + lang + "/"+ platform + "/sorties"
	messageJSON, _ := json.Marshal(sortie)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseSyndicateMissions(platformno int, platform string, c mqtt.Client) {
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
	data := apidata[platformno]
	var syndicates []SyndicateMissions
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		syndicatecheck, _ := jsonparser.GetString(value, "syndicate")
		if syndicatecheck != "Ostrons" && syndicatecheck != "Solaris United" {
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
	}, "syndicateMissions")

	topicf := "/wf/" + platform + "/"+ langtest + "/syndicates"
	messageJSON, _ := json.Marshal(syndicates)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()

}
func parseInvasions(platformno int, platform string, c mqtt.Client) {
	type Invasion struct {
		ID                  string
		Location            string
		MissionType         string
		Completed           bool
		Started             string
		VsInfested          bool
		AttackerRewardItem  string `json:",omitempty"`
		AttackerRewardCount int64  `json:",omitempty"`
		AttackerMissionInfo string `json:",omitempty"`
		DefenderRewardItem  string `json:",omitempty"`
		DefenderRewardCount int64  `json:",omitempty"`
		DefenderMissionInfo string `json:",omitempty"`
		Completion          float64
	}

	data := apidata[platformno]
	invasioncheck, _, _, _ := jsonparser.Get(data, "invasions")
	if len(invasioncheck) == 0 {
		topicf := "/wf/" + platform + "/"+ langtest + "/invasions"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		return
	}
	var invasions []Invasion
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		iscomplete, _ := jsonparser.GetBoolean(value, "completed")
		if iscomplete != true {
			attackeritem := ""
			attackeritemcount := int64(0)
			defenderitem := ""
			defenderitemcount := int64(0)
			id, _ := jsonparser.GetString(value, "id")
			started, _ := jsonparser.GetString(value, "activation")
			location, _ := jsonparser.GetString(value, "node")
			missiontype, _ := jsonparser.GetString(value, "desc")
			completed, _ := jsonparser.GetBoolean(value, "completed")
			vsinfested, _ := jsonparser.GetBoolean(value, "vsInfestation")
			_, _, _, ierror := jsonparser.Get(value, "attackerReward", "countedItems", "[0]", "type")
			if ierror == nil {
				attackeritem, _ = jsonparser.GetString(value, "attackerReward", "countedItems", "[0]", "type")
				attackeritemcount, _ = jsonparser.GetInt(value, "attackerReward", "countedItems", "[0]", "count")
			}
			attackerfaction, _ := jsonparser.GetString(value, "attackingFaction")
			_, _, _, ierror2 := jsonparser.Get(value, "defenderReward", "countedItems", "[0]", "type")
			if ierror2 == nil {
				defenderitem, _ = jsonparser.GetString(value, "defenderReward", "countedItems", "[0]", "type")
				defenderitemcount, _ = jsonparser.GetInt(value, "defenderReward", "countedItems", "[0]", "count")
			}
			defenderfaction, _ := jsonparser.GetString(value, "defendingFaction")
			completion, _ := jsonparser.GetFloat(value, "completion")
			w := Invasion{id, location, missiontype, completed, started, vsinfested,
				attackeritem, attackeritemcount, attackerfaction,
				defenderitem, defenderitemcount, defenderfaction, completion}
			invasions = append(invasions, w)
		}
	}, "invasions")

	topicf := "/wf/" + platform + "/"+ langtest + "/invasions"
	messageJSON, _ := json.Marshal(invasions)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

	/*
func parseActiveMissions(platformno int, platform string, c mqtt.Client) {
	type ActiveMissions struct {
		ID          string
		Started     int
		Ends        int
		Region      int
		Node        string
		MissionType string
		Modifier    string
	}
	data := &apidata[platformno]
	fsion := gofasion.NewFasion(*data)
	var mission []ActiveMissions
	lang := string("en")
	ActiveMissionsarray := fsion.Get("ActiveMissions").Array()
	fmt.Println(len(ActiveMissionsarray))

	for _, v := range ActiveMissionsarray {
		id := v.Get("_id").Get("$oid").ValueStr()
		started := v.Get("Activation").Get("$date").Get("$numberLong").ValueInt() / 1000
		ended := v.Get("Expiry").Get("$date").Get("$numberLong").ValueInt() / 1000
		region := v.Get("Region").ValueInt()
		node := v.Get("Node").ValueStr()
		missiontype := v.Get("MissionType").ValueStr()
		modifier := v.Get("Modifier").ValueStr()

		w := ActiveMissions{
			ID:          id,
			Started:     started,
			Ends:        ended,
			Region:      region,
			Node:        node,
			MissionType: missiontype,
			Modifier:    modifier,
		}
		mission = append(mission, w)
	}
	fmt.Println(len(mission))
	topicf := "/wf/" + platform + "/"+ langtest + "/missions"
	messageJSON, _ := json.Marshal(mission)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}
	*/
func calcCompletion(count int, goal int, attacker string) (complete float32) {
	y := float32((1 + float32(count)/float32(goal)))
	x := float32(y * 50)
	if attacker == "Infested" {
		x = float32(y * 100)

	}
	//fmt.Println(y)
	return x
}

// PrintMemUsage - only for debug
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func sortietranslate(src string, langtype string, lang string) (ret [2]string) {
	if langtype == "sortiemod" {
		var x1 [2]string
		x1[0] = src
		x1[1] = src

		result, ok := sortiemodtypes[lang][src]
		if ok != false {
			x1[0] = result.(string)
		}
		result2, ok := sortiemoddesc[lang][src]
		if ok != false {
			x1[1] = result2.(string)
		}

		ret = x1
	}

	if langtype == "sortiemodboss" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src

		result, ok := sortiemodbosses[lang][src].(map[string]interface{})
 		if ok != false {
			x1[0] = result["faction"].(string)
			x1[1] = result["name"].(string)

		}

		ret = x1

	}
	
	if langtype == "sortieloc" {
		var x1 [2]string

		x1[0] = src
		x1[1] = src
		result, ok := sortieloc[lang][src].(map[string]interface{})
	fmt.Println("test2",sortieloc[lang][src])
		
 		if ok != false {
			x1[0] = result["value"].(string)
			x1[1] = result["enemy"].(string)

		}
	/**/
		ret = x1

	}
	/**/
	return ret
}
func voidranslate(src string, lang string) (ret [2]string) {
		var x1 [2]string
		x1[0] = src
		x1[1] = src

		result, ok := fissureModifiers[lang][src].(map[string]interface{})
		if ok != false {
			x1[0] = result["value"].(string)
			x1[1] =fmt.Sprintf("%.0f", result["num"].(float64))

		}
		ret = x1
		return ret
}
func missionranslate(src string, lang string) (ret string) {
		var x1 string
		x1 = src

		result, ok := missionTypes[lang][src].(map[string]interface{})
		if ok != false {
			x1 = result["value"].(string)

		}
		ret = x1
		return ret
}
// FloatToString convert
func FloatToString(inputnum float64) string {
	// to convert a float number to a string
    return strconv.FormatFloat(inputnum, 'f', 6, 64)
}

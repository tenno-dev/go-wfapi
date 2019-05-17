package parser

import (
	"encoding/json"
	"fmt"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//Voidtrader - Voidtrader
type Voidtrader struct {
	ID      string
	Started string
	Ends    string
	NPC     string
	Node    string
	Offers  []VoidtraderOffers `json:",omitempty"`
}

//VoidtraderOffers - VoidtraderOffers
type VoidtraderOffers struct {
	Item    string `json:",omitempty"`
	Ducats  int64  `json:",omitempty"`
	Credits int64  `json:",omitempty"`
}

// ParseVoidTrader Parse Void trader
func ParseVoidTrader(platformno int, platform string, c mqtt.Client, lang string) {
	data := datasources.Apidata[platformno]
	var voidtrader []Voidtrader

	_, _, _, voiderr := jsonparser.Get(data, "VoidTraders")
	if voiderr != nil {
		topicf := "/wf/" + lang + "/" + platform + "/voidtrader"
		token := c.Publish(topicf, 0, true, []byte("{}"))
		token.Wait()
		fmt.Println("reached void error")

		return
	}
	fmt.Println("reached sortie start")

	id, _ := jsonparser.GetString(data, "VoidTraders", "[0]", "_id", "$oid")
	started, _ := jsonparser.GetString(data, "VoidTraders", "[0]", "Activation", "$date", "$numberLong")
	ended, _ := jsonparser.GetString(data, "VoidTraders", "[0]", "Expiry", "$date", "$numberLong")
	npc, _ := jsonparser.GetString(data, "VoidTraders", "[0]", "Character")
	location1, _ := jsonparser.GetString(data, "VoidTraders", "[0]", "Node")
	location := helper.Sortietranslate(location1, "sortieloc", lang)
	var voidoffers []VoidtraderOffers

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		item1, _ := jsonparser.GetString(value, "ItemType")
		item := helper.Langtranslate1(item1, lang)
		ducats, _ := jsonparser.GetInt(value, "PrimePrice")
		credits, _ := jsonparser.GetInt(value, "RegularPrice")
		voidoffers = append(voidoffers, VoidtraderOffers{
			Item:    item,
			Ducats:  ducats,
			Credits: credits,
		})
	}, "VoidTraders", "[0]", "Manifest")
	w := Voidtrader{ID: id, Started: started,
		Ends: ended, NPC: npc, Node: location[0], Offers: voidoffers}
	voidtrader = append(voidtrader, w)

	topicf := "/wf/" + lang + "/" + platform + "/voidtrader"
	messageJSON, _ := json.Marshal(voidtrader)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

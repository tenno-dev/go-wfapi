package parser

import (
	"sync"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/helper"
	"github.com/buger/jsonparser"
)

// DarvoDeals struct
type DarvoDeals struct {
	ID              string
	Start           string
	Ends            string
	Item            string
	Itemtest        string
	Price           int64
	DealPrice       int64
	DiscountPercent int64
	Stock           int64
	Sold            int64
}

// Darvodata for http export
var Darvodata = make(map[int]map[string][]DarvoDeals)

// ParseDarvoDeal Parse current Darvo Deal
func ParseDarvoDeal(platformno int, platform string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := Darvodata[platformno]; !ok {
		Darvodata[platformno] = make(map[string][]DarvoDeals)
	}
	data := datasources.Apidata[platformno]
	var deals []DarvoDeals
	// fmt.Println("Darvo  reached")
	errfissures, _ := jsonparser.GetString(data, "DailyDeals")
	if errfissures != "" {
		// fmt.Println("error Darvo reached")
		return
	}
	// fmt.Println("Darvo2 reached")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//id, _ := jsonparser.GetString(value, "id")
		id := "1"
		started, _ := jsonparser.GetString(value, "Activation", "$date", "$numberLong")
		ended, _ := jsonparser.GetString(value, "Expiry", "$date", "$numberLong")
		item, _ := jsonparser.GetString(value, "StoreItem")
		item = helper.Langtranslate1(item, lang)
		// itemtest := helper.Langtranslate1(item, "en")
		// itemdetails := helper.Getitemdetails(itemtest)
		itemdetails := "PH"

		originalprice, _ := jsonparser.GetInt(value, "OriginalPrice")
		dealprice, _ := jsonparser.GetInt(value, "SalePrice")
		stock, _ := jsonparser.GetInt(value, "AmountTotal")
		sold, _ := jsonparser.GetInt(value, "AmountSold")
		discount, _ := jsonparser.GetInt(value, "Discount")

		w := DarvoDeals{id, started, ended, item, itemdetails, originalprice, dealprice,
			discount, stock, sold}
		deals = append(deals, w)
	}, "DailyDeals")
	Darvodata[platformno][lang] = deals

}

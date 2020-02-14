package parser

import (
	"encoding/json"
	"strings"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// News struct
type News struct {
	ID       string
	Message  string
	URL      string
	Date     string
	priority bool
	Image    string
}

// Newsdata - test news output
var Newsdata = make(map[int]map[string][]News)

// ParseNews parsing news data (Called Events in warframe api)
func ParseNews(platformno int, platform string, c mqtt.Client, lang string) {
	if _, ok := Newsdata[platformno]; !ok {
		Newsdata[platformno] = make(map[string][]News)
	}
	data := datasources.Apidata[platformno]
	_, _, _, ernews := jsonparser.Get(data, "Events")
	if ernews != nil {
		// fmt.Println("error ernews reached")
		return
	}
	var errnews2 bool
	var message string

	var news []News

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		message = ""
		jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			newstemp1, _ := jsonparser.GetString(value1, "LanguageCode")

			if newstemp1 == lang {
				message, _ = jsonparser.GetString(value1, "Message")

			}
		}, "Messages")

		if message != "" {
			errnews2 = false
		} else {
			errnews2 = true
		}

		if errnews2 == false {
			image := "http://n9e5v4d8.ssl.hwcdn.net/uploads/e0b4d18d3330bb0e62dcdcb364d5f004.png"
			id, _ := jsonparser.GetString(value, "_id", "$oid")

			url, _ := jsonparser.GetString(value, "Prop")
			_, imgerr := jsonparser.GetString(value, "ImageUrl")
			if imgerr == nil {
				image, _ = jsonparser.GetString(value, "ImageUrl")
			}
			if strings.HasPrefix(image, "https://forums.warframe.com") {
				image = strings.Split(image, "=")[1]
				image = strings.Split(image, "&key")[0]

			}
			date, _ := jsonparser.GetString(value, "Date", "$date", "$numberLong")
			priority, _ := jsonparser.GetBoolean(value, "priority")
			w := News{ID: id, Message: message, URL: url, Date: date, Image: image, priority: priority}
			news = append(news, w)
		}
	}, "Events")
	topicf := "wf/" + lang + "/" + platform + "/news"
	Newsdata[platformno][lang] = news
	messageJSON, _ := json.Marshal(news)
	token := c.Publish(topicf, 0, true, messageJSON)
	token.Wait()
}

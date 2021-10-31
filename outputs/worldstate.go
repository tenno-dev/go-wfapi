package outputs

import (
	"encoding/json"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/gofiber/fiber/v2"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

// Everything test 2
func Everything(c *fiber.Ctx) error {
	platform := c.Params("platform")
	value, _ := intMap[platform]
	c.Type("json")
	return c.Send(datasources.Apidata[value])
}

// DarvoDeals DarvoDeals
func DarvoDeals(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Darvodata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

// News godoc
func News(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Newsdata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

// Alerts Alertsdata
func Alerts(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Alertsdata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

// Fissures Fissuresdata
func Fissures(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Fissuresdata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

// Nightwave data
func Nightwave(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Nightwavedata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

// Penemy sdata
func Penemy(c *fiber.Ctx) error {
	v, t1, t2 := getPlatformValueAndTokens(c)
	messageJSON, _ := json.Marshal(parser.Penemydata[v][t1])
	return sendJSONMessage(c, messageJSON, t2)
}

func sendJSONMessage(c *fiber.Ctx, msg []byte, token2 string) error {
	c.Type("json")
	c.Set("Content-Language", token2)
	return c.Send(msg)
}

func getPlatformValueAndTokens(c *fiber.Ctx) (int, string, string) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]

	return value, token1, token2
}

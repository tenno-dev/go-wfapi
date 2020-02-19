package outputs

import (
	"encoding/json"

	"github.com/bitti09/go-wfapi/datasources"
	"github.com/bitti09/go-wfapi/parser"
	"github.com/gofiber/fiber"
)

var intMap = map[string]int{"pc": 0, "ps4": 1, "xb1": 2, "swi": 3}

// IndexHandler test
func IndexHandler(c *fiber.Ctx) {
	token := c.Get("Accept-Language")
	token = token[0:2]
	c.Send(token)
}

// Everything test 2
func Everything(c *fiber.Ctx) {
	platform := c.Params("platform")
	value, _ := intMap[platform]
	c.Type("json")
	c.Send(datasources.Apidata[value])
}

// DarvoDeals DarvoDeals
func DarvoDeals(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	token1 := token[0:2]
	token2 := token[0:5]
	value, _ := intMap[platform]
	messageJSON, _ := json.Marshal(parser.Darvodata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

// News godoc
func News(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]
	messageJSON, _ := json.Marshal(parser.Newsdata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

// Alerts Alertsdata
func Alerts(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]
	messageJSON, _ := json.Marshal(parser.Alertsdata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

// Fissures Fissuresdata
func Fissures(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]
	messageJSON, _ := json.Marshal(parser.Fissuresdata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

// Nightwave data
func Nightwave(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]
	messageJSON, _ := json.Marshal(parser.Nightwavedata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

// Penemy sdata
func Penemy(c *fiber.Ctx) {
	platform := c.Params("platform")
	token := c.Get("Accept-Language")
	value, _ := intMap[platform]
	token1 := token[0:2]
	token2 := token[0:5]
	messageJSON, _ := json.Marshal(parser.Penemydata[value][token1])
	c.Type("json")
	c.Set("Content-Language", token2)
	c.Send(messageJSON)
}

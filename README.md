# go-wfapi
WIP  Warframe API Parser with  MQTT publisher

## Current Status:

### Parser:
 - [x] News
 - [x] Sorties
 - [x] Void Fissures
 - [ ] Alerts (waiting for api response)
 - [ ] Darvo's Deals
 - [ ] Nightwave
 - [ ] Syndicate Missions
 - [ ] Invasions
 - [ ] Void Trader
 - [ ] World Events


### General

- [ ] Code Rewrite  
- [ ] Code Spliting
- [ ] Rewrite of the  current "hacky" ways for json parsing


## Demo

+ Web:  coming soon
+ MQTT Client: 
  + Host: mybitti(.)de
  + Port: 1884
  + Protocol: wss
  + Data path: /wf/{lang}/{platform}/
  + Lang: {"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"}
  + Platform: {"pc", "ps4", "xb1", "swi"}
  + Tested Client: [MQTT Explorer](https://mqtt-explorer.com/)


# go-wfapi

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6cd589b3d93c4185af48815346a99d44)](https://app.codacy.com/manual/Bitti09/go-wfapi?utm_source=github.com&utm_medium=referral&utm_content=Bitti09/go-wfapi&utm_campaign=Badge_Grade_Settings)

WIP  Warframe API Parser with  MQTT publisher

## Current Status

### Parser

-   [x] News
-   [x] Sorties
-   [x] Void Fissures
-   [ ] Alerts (waiting for api response)
-   [x] Darvo's Deals ( untranslated )
-   [x] Nightwave 
-   [x] Syndicate Missions ( untranslated )
-   [x] Invasions  ( half translated )
-   [x] Void Trader
-   [-] World Events (partly)
-   [x] Kuva Missions (***new***)
-   [x] Arbitration (***new***)

## Demo

-   Web: <https://api.mybitti.de>

    -   Ping:   /
    -   Worldstate(unparsed) /:platform
    -   Darvos Deal: /:platform/platform/  (Accept-Language required)
    -   News /:platform/news/ (Accept-Language required)
    -   Alerts /:platform/alerts/ (Accept-Language required)
    -   Fissures /:platform/fissures/ (Accept-Language required)
    -   Nightwave /:platform/nightwave/ (Accept-Language required)

    only first two chars  of  (Accept-Language) is used.

-   MQTT Client:
    -   Host: api-wf(.)mybitti(.)de
    -   Port: 8084
    -   Protocol: wss
    -   Data path: wf/{lang}/{platform}/
    -   Lang: {"en", "de", "es", "fr","it","ko","pl","pt","ru","zh"}
    -   Platform: {"pc", "ps4", "xb1", "swi"}
    -   Tested Client: [MQTT Explorer](https://mqtt-explorer.com/)

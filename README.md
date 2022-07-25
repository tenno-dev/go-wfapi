# go-wfapi

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/80fda50c42614ce582c2813bd7847904)](https://app.codacy.com/gh/tenno-dev/go-wfapi?utm_source=github.com&utm_medium=referral&utm_content=tenno-dev/go-wfapi&utm_campaign=Badge_Grade)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/tenno-dev/go-wfapi)

WIP  Warframe API Parser  -- currently rewriting/resturcturing the codebase

## Current Status

### Parser

-   [x] News
-   [x] Sorties
-   [x] Void Fissures
-   [ ] Alerts (waiting for api response)
-   [x] Darvo's Deals
-   [x] Nightwave 
-   [ ] Syndicate Missions ( incomplete )
-   [x] Invasions  ( half translated )
-   [x] Void Trader
-   [x] World Events
-   [ ] Arbitration
-   [ ] Timers (incomplete)

## Demo

-   Web:
    -  URL: [api.tenno.dev](api.tenno.dev)
    -   Ping:   /
    -   Worldstate(unparsed) /:platform
    -   Darvos Deal: /:platform/platform/  (Accept-Language required)
    -   News /:platform/news/ (Accept-Language required)
    -   Alerts /:platform/alerts/ (Accept-Language required)
    -   Fissures /:platform/fissures/ (Accept-Language required)
    -   Nightwave /:platform/nightwave/ (Accept-Language required)
    -   Everythin  in one JSON  /:platform/test
    - Lang Select via Query:  ?lang=(language)  { currently only **en** and **de** is  possible}
  
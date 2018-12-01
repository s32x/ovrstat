# ovrstat

[![Circle CI](https://circleci.com/gh/s32x/ovrstat/tree/master.svg?style=svg)](https://circleci.com/gh/s32x/ovrstat/tree/master)
[![GoDoc](https://godoc.org/github.com/s32x/ovrstat/goow?status.svg)](https://godoc.org/github.com/s32x/ovrstat/goow)

![](web/assets/ovrstatdarksmall.png "ovrstat")

ovrstat is a simple web scraper for the Overwatch stats site that parses and serves the data retrieved as JSON. Included is the go package used to scrape the info for usage in any go binary.

Note: This is a single endpoint web-scraping API that takes the full payload of information that we retrieve from Blizzard and passes it through to you in a single response. Things like caching and splitting data across multiple responses could likely improve performance, but in pursuit of keeping things simple, ovrstat does not implement them.

### Running with Docker
```
docker run -p 8080:8080 s23x/ovrstat
```
### Installing
```
go get s32x.com/ovrstat
ovrstat
```
### Usage

You have two options for using the API: 
* Import the child dependency used in this API and use the API we host on Heroku
* Host your own Ovrstat API using the public docker image `s32x/ovrstat`.

Below is an example of using the REST endpoint (note: CASE matters for the username/tag):
```
https://ovrstat.com/stats/pc/us/Viz-1213
https://ovrstat.com/stats/xbl/Lt%20Evolution
https://ovrstat.com/stats/psn/TayuyaBreast
```

And here is an example of using the included go library:
```go
player, _ := ovrstat.PCStats("us", "Viz-1213")
player2, _ := ovrstat.ConsoleStats(ovrstat.PlatformXBL, "Lt%20Evolution")
player3, _ := ovrstat.ConsoleStats(ovrstat.PlatformPSN, "TayuyaBreast")
```
Both above examples should return to you a PlayerStats struct containing detailed statistics for the specified Overwatch player.

## Full Go example

```go
package main

import (
	"log"

	"s32x.com/ovrstat/ovrstat"
)

func main() {
	log.Println(ovrstat.PCStats("us", "Viz-1213"))
	log.Println(ovrstat.ConsoleStats(ovrstat.PlatformXBL, "Lt%20Evolution"))
	log.Println(ovrstat.ConsoleStats(ovrstat.PlatformPSN, "TayuyaBreast"))
}
```

## Disclaimer
ovrstat isn’t endorsed by Blizzard and doesn’t reflect the views or opinions of Blizzard or anyone officially involved in producing or managing Overwatch. Overwatch and Blizzard are trademarks or registered trademarks of Blizzard Entertainment, Inc. Overwatch © Blizzard Entertainment, Inc.

The BSD 3-clause License
========================

Copyright (c) 2018, Steven Wolfe. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

 - Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

 - Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

 - Neither the name of ovrstat nor the names of its contributors may
   be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

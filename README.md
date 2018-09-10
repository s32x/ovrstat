# ovrstat

[![GoDoc](https://godoc.org/github.com/starboy/ovrstat/goow?status.svg)](https://godoc.org/github.com/starboy/ovrstat/goow)

![](web/assets/ovrstatdarksmall.png "ovrstat")

ovrstat is a simple web scraper for the Overwatch stats site that parses and serves the data retrieved as JSON. Included is the go package used to scrape the info for usage in any go binary.

Note: This is a single endpoint web-scraping API that takes the full payload of information that we retrieve from Blizzard and passes it through to you in a single response. Things like caching and splitting data across multiple responses could likely improve performance, but in pursuit of keeping things simple, ovrstat does not implement them.

### Running with Docker
```
docker run -p 8080:8080 starboy/ovrstat
```
### Installing
```
go get github.com/starboy/ovrstat/ovrstat
```
### Usage

You have two options for using the API: 
* Import the child dependency used in this API and use the API we host on Heroku
* Host your own Ovrstat API using the public docker image `starboy/ovrstat`.

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

	"github.com/starboy/ovrstat/ovrstat"
)

func main() {
	log.Println(ovrstat.PCStats("us", "Viz-1213"))
	log.Println(ovrstat.ConsoleStats(ovrstat.PlatformXBL, "Lt%20Evolution"))
	log.Println(ovrstat.ConsoleStats(ovrstat.PlatformPSN, "TayuyaBreast"))
}
```

## Disclaimer
ovrstat isn’t endorsed by Blizzard and doesn’t reflect the views or opinions of Blizzard or anyone officially involved in producing or managing Overwatch. Overwatch and Blizzard are trademarks or registered trademarks of Blizzard Entertainment, Inc. Overwatch © Blizzard Entertainment, Inc.

The MIT License (MIT)
=====================

Copyright © 2018 Steven Wolfe

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the “Software”), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

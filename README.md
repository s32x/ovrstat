# ovrstat

[![GoDoc](https://godoc.org/github.com/sdwolfe32/ovrstat/goow?status.svg)](https://godoc.org/github.com/sdwolfe32/ovrstat/goow)

![alt text](/img/ovrstatdarksmall.png "ovrstat")

ovrstat is a simple web scraper for the Overwatch stats site that parses and serves the data retrieved as JSON. Also included is goow, a binding used to retrieve the stats that can be used as an Overwatch API go dep.

Note: As this is a web-scraping API I saw no reason to serve separate data across multiple requests. While caching could be an option to save bandwidth on your end, I didn't see any reason not to give you back as much information as we retrieve from Blizzard, thus there is only one endpoint currently.

### Running with Docker
```
docker run sdwolfe32/ovrstat
```
### Installing
```
go get github.com/sdwolfe32/ovrstat/goow
```
### Usage

You have two options for using the API, Either import the child dependency used in this API, use the API we host on heroku, or host your own ovrstat API using the public docker image `sdwolfe32/ovrstat`.

Below is an example of using the REST endpoint:
```
http://ovrstat.com/v1/stats/pc/us/Viz-1213
http://ovrstat.com/v1/stats/xbox/Viz-1213
```

And here is an example of using the included go library:
```
player, _ := ovrstat.PCStats("us", "Viz-1213")
player2, _ := ovrstat.ConsoleStats("xbox")
```
Both above examples should return to you a PlayerStats struct containing detailed statistics for the specified Overwatch player.

## Full Go example

```
package main

import (
	"log"

	"github.com/sdwolfe32/ovrstat/ovrstat"
)

func main() {
	log.Println(ovrstat.PCStats("us", "Viz-1213"))
}
```

## Disclaimer
ovrstat isn’t endorsed by Blizzard and doesn’t reflect the views or opinions of Blizzard or anyone officially involved in producing or managing Overwatch. Overwatch and Blizzard  are trademarks or registered trademarks of Blizzard Entertainment, Inc. Overwatch © Blizzard Entertainment, Inc.

The MIT License (MIT)
=====================

Copyright © 2017 Steven Wolfe

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

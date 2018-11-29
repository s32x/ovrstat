<div align="center">
    <img src="logo.png" height="101" width="350" /><br/>
    <a href="https://godoc.org/github.com/starboy/httpclient">
        <img src="https://godoc.org/github.com/starboy/httpclient?status.svg" />
    </a>
</div>

httpclient is a simple convenience package for performing http/api requests in Go. It wraps the standard libraries net/http package to avoid the repetitive request->decode->closebody logic you're likely so familiar with. Using the lib is very simple - just as with the net/http package you can define a client of your own or use the default. Below is a very basic example.

### Usage

```go
package main

import "github.com/starboy/httpclient"

func main() {
	s, err := httpclient.GetString("https://api.github.com/users/starboy/repos")
	if err != nil {
		panic(err)
	}
	println(s)
}
```

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
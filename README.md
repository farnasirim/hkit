# HKIT

`hkit` is an http toolkit which implements handy, easy to configure facades.

To Install, do as you would with go modules:

```bash
vgo get github.com/farnasirim/hkit
```
You can use old fashioned `go get` too if you want.

## Docs

hkit currently implements utilities for the following use cases:
- [Logging](#logging)

### <a name="logging"></a> Logging

You can use hkit.Logger with any `http.Handler` or `http.HandlerFunc` as easy as the following:
#### net/http
```go
package main

import (
	"net/http"

	"github.com/farnasirim/hkit"
)

func main() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"x": "2", "y": {"a": "b"}}`))
	}

	wrappedHandler := hkit.NewLogger(handlerFunc)

	http.ListenAndServe(":8000", wrappedHandler)
}
```
And then
```bash
curl localhost:8000 -d '{"some": "json"}'
```

will result in
```
Method: POST
remote address: [::1]:39252

Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.61.1
Accept: */*
Content-Length: 16

{"some": "json"}

{"x": "2", "y": {"a": "b"}}
```

#### logrus
This is the easiest way you can use hkit with logrus:
```go
package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/colonelmo/hkit"
)

func main() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"x": "2", "y": {"a": "b"}}`))
	}

	prettyLogWriter := log.StandardLogger().Writer()

	wrappedHandler := hkit.NewLogger(handlerFunc).SetWriter(prettyLogWriter)

	http.ListenAndServe(":8000", wrappedHandler)
}
```

And the previous curl will result in:
<pre><font color="#00AAAA">INFO</font>[0021] Method: POST                                 
<font color="#00AAAA">INFO</font>[0021] remote address: [::1]:39358                  
<font color="#00AAAA">INFO</font>[0021]                                              
<font color="#00AAAA">INFO</font>[0021] User-Agent: curl/7.61.1                      
<font color="#00AAAA">INFO</font>[0021] Accept: */*                                  
<font color="#00AAAA">INFO</font>[0021] Content-Length: 16                           
<font color="#00AAAA">INFO</font>[0021] Content-Type: application/x-www-form-urlencoded 
<font color="#00AAAA">INFO</font>[0021]                                              
<font color="#00AAAA">INFO</font>[0021] {&quot;some&quot;: &quot;json&quot;}                             
<font color="#00AAAA">INFO</font>[0021]   </pre>

being written to the stdout.


# License

MIT

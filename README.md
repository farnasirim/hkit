# HKIT

an http toolkit containing useful http utilities.

hkit currently implements utilities for the following use cases:
- [Logging](#logging)

### <a name="logging"></a> Logging

You can use hkit.Logger with any `http.Handler` or `http.HandlerFunc` as easy as the following:
#### net/http
```go

package main

import (
	"net/http"

	"github.com/colonelmo/hkit"
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
remote address: [::1]:57770

User-Agent: curl/7.54.1
Accept: */*
Content-Length: 16
Content-Type: application/x-www-form-urlencoded

{"some": "json"}
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

Which will result in:
```
INFO[0001] Method: POST
INFO[0001] remote address: [::1]:42424
INFO[0001]
INFO[0001] User-Agent: curl/7.54.1
INFO[0001] Accept: */*
INFO[0001] Content-Length: 16
INFO[0001] Content-Type: application/x-www-form-urlencoded
INFO[0001]
```

being written to the stdout.


# License

BSD-3-clause

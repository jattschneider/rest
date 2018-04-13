# Go Rest

Rest is a Go library that implements a REST JSON client.

## Getting Started

Just a quick example how to use the rest library:

#### main.go
```
package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jattschneider/rest"
)

func init() {
	flag.Parse()
	rest.Authentication = func(r *http.Request) {
		r.SetBasicAuth("user", "password")
	}
}

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		fmt.Printf("BasicAuth usr: %v pwd: %v ok?: %v\n", username, password, ok)
	}))

	re, err := rest.Get(ts.URL)
	if err != nil || re.StatusCode != http.StatusOK {
		return
	}

	payload := rest.EncodeJSON(&struct{ SomeProperty string }{SomeProperty: "struct property value"})

	re, err = rest.Put(ts.URL, payload)
	if err != nil || re.StatusCode != http.StatusOK {
		return
	}

	re, err = rest.Post(ts.URL, payload)
	if err != nil || re.StatusCode != http.StatusOK {
		return
	}

	re, err = rest.Patch(ts.URL, payload)
	if err != nil || re.StatusCode != http.StatusOK {
		return
	}
}
```

```
$ go run main.go
```

### Installing

```
go get -v github.com/jattschneider/rest
```

## Built With

* [Go](https://golang.org/) - The Go Programming Language
* [dep](https://golang.github.io/dep/) - Dependency management for Go
* [glog](https://github.com/golang/glog) - Leveled execution logs for Go

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/jattschneider/config/tags). 

## Authors

* **José Augusto Schneider** - *Initial work* - [jattschneider](https://github.com/jattschneider)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

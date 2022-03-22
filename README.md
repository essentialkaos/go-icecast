<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-icecast.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/go-icecast"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev" /></a>
  <a href="https://kaos.sh/r/go-icecast"><img src="https://kaos.sh/r/go-icecast.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/go-icecast/ci"><img src="https://kaos.sh/w/go-icecast/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/go-icecast/codeql"><img src="https://kaos.sh/w/go-icecast/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/b/go-icecast"><img src="https://kaos.sh/b/b2237e1d-2089-40f3-bfa1-f66bc79c68a8.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/c/go-icecast"><img src="https://kaos.sh/c/go-icecast.svg" alt="Coverage Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-icecast` is a Go package for working with [Icecast Admin API](http://icecast.org/docs/icecast-2.4.1/admin-interface.html).

### Installation

Make sure you have a working Go 1.17+ workspace (_[instructions](https://golang.org/doc/install)_), then:

````
go get -d github.com/essentialkaos/go-icecast/v2
````

For update to latest stable release, do:

```
go get -d -u github.com/essentialkaos/go-icecast/v2
```

### Usage example

```go
package main

import (
  "fmt"
  ic "pkg.re/essentialkaos/go-icecast/v2"
)

func main() {
  api, err := ic.NewAPI("https://127.0.0.1:8000", "admin", "MySuppaPAssWOrd")
  api.SetUserAgent("MyApp", "1.2.3")

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  stats, err := api.GetStats()

  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }

  fmt.Println("%-v\n", stats)
}
```

### Build Status

| Branch     | Status |
|------------|--------|
| `master` (_Stable_) | [![CI](https://kaos.sh/w/go-icecast/ci.svg?branch=master)](https://kaos.sh/w/go-icecast/ci?query=branch:master) |
| `develop` (_Unstable_) | [![CI](https://kaos.sh/w/go-icecast/ci.svg?branch=develop)](https://kaos.sh/w/go-icecast/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>

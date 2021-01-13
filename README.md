<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-icecast.svg"/></a></p>

<p align="center">
  <a href="https://pkg.re/essentialkaos/go-icecast.v2?docs"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/go-icecast"><img src="https://goreportcard.com/badge/github.com/essentialkaos/go-icecast"></a>
  <a href="https://github.com/essentialkaos/go-icecast/actions"><img src="https://github.com/essentialkaos/go-icecast/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://github.com/essentialkaos/go-icecast/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/go-icecast/workflows/CodeQL/badge.svg" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-go-icecast-master"><img alt="codebeat badge" src="https://codebeat.co/badges/b2237e1d-2089-40f3-bfa1-f66bc79c68a8"></a>
  <a href='https://coveralls.io/github/essentialkaos/go-icecast?branch=develop'><img src='https://coveralls.io/repos/github/essentialkaos/go-icecast/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-icecast` is a Go package for working with [Icecast Admin API](http://icecast.org/docs/icecast-2.4.1/admin-interface.html).

### Installation

Make sure you have a working Go 1.14+ workspace (_[instructions](https://golang.org/doc/install)_), then:

````
go get pkg.re/essentialkaos/go-icecast.v2
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/go-icecast.v2
```

### Usage example

```go
package main

import (
  "fmt"
  ic "pkg.re/essentialkaos/go-icecast.v2"
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
| `master` (_Stable_) | [![CI](https://github.com/essentialkaos/go-icecast/workflows/CI/badge.svg?branch=master)](https://github.com/essentialkaos/go-icecast/actions) |
| `develop` (_Unstable_) | [![CI](https://github.com/essentialkaos/go-icecast/workflows/CI/badge.svg?branch=develop)](https://github.com/essentialkaos/go-icecast/actions) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>

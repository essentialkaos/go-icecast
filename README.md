<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-icecast.svg"/></a></p>

<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/go-icecast.v1"><img src="https://godoc.org/pkg.re/essentialkaos/go-icecast.v1?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/go-icecast"><img src="https://goreportcard.com/badge/github.com/essentialkaos/go-icecast"></a>
  <a href="https://travis-ci.com/essentialkaos/go-icecast"><img src="https://travis-ci.com/essentialkaos/go-icecast.svg"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-go-icecast-master"><img alt="codebeat badge" src="https://codebeat.co/badges/b2237e1d-2089-40f3-bfa1-f66bc79c68a8"></a>
  <a href='https://coveralls.io/github/essentialkaos/go-icecast?branch=develop'><img src='https://coveralls.io/repos/github/essentialkaos/go-icecast/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`go-icecast` is a Go package for working with [Icecast Admin API](http://icecast.org/docs/icecast-2.4.1/admin-interface.html).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.13+ workspace (_[instructions](https://golang.org/doc/install)_), then:

````
go get pkg.re/essentialkaos/go-icecast.v1
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/go-icecast.v1
```

### Usage example

```go
package main

import (
  "fmt"
  ic "pkg.re/essentialkaos/go-icecast.v1"
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
| `master` (_Stable_) | [![Build Status](https://travis-ci.com/essentialkaos/go-icecast.svg?branch=master)](https://travis-ci.com/essentialkaos/go-icecast) |
| `develop` (_Unstable_) | [![Build Status](https://travis-ci.com/essentialkaos/go-icecast.svg?branch=develop)](https://travis-ci.com/essentialkaos/go-icecast) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>

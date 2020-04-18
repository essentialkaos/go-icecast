package icecast

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// API is Confluence API struct
type API struct {
	Client *fasthttp.Client // Client is client for http requests

	url       string // confluence URL
	basicAuth string // basic auth
}

// ////////////////////////////////////////////////////////////////////////////////// //

// API errors
var (
	ErrInitEmptyURL      = errors.New("URL can't be empty")
	ErrInitEmptyApp      = errors.New("App can't be empty")
	ErrInitEmptyPassword = errors.New("Password can't be empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewAPI creates new API struct
func NewAPI(url, app, password string) (*API, error) {
	switch {
	case url == "":
		return nil, ErrInitEmptyURL
	case app == "":
		return nil, ErrInitEmptyApp
	case password == "":
		return nil, ErrInitEmptyPassword
	}

	return &API{
		Client: &fasthttp.Client{
			Name:                getUserAgent("", ""),
			MaxIdleConnDuration: 5 * time.Second,
			ReadTimeout:         3 * time.Second,
			WriteTimeout:        3 * time.Second,
			MaxConnsPerHost:     150,
		},

		url:       url,
		basicAuth: genBasicAuthHeader(app, password),
	}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// SetUserAgent set user-agent string based on app name and version
func (api *API) SetUserAgent(app, version string) {
	api.Client.Name = getUserAgent(app, version)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetInfo fetches info about Icecast server
func (api *API) GetInfo() (*Server, error) {
	data, err := api.doRequest("GET", "/admin/stats")

	if err != nil {
		return nil, err
	}

	return parseStatsData(data)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// parseStatsData parses raw XML data with Icecast stats
func parseStatsData(data []byte) (*Server, error) {
	server := &iceServer{}
	err := xml.Unmarshal(data, server)

	if err != nil {
		return nil, err
	}

	return convertStats(server), nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// codebeat:disable[ARITY]

// doRequest creates and executes request
func (api *API) doRequest(method, uri string) ([]byte, error) {
	req := api.acquireRequest(method, uri)
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	err := api.Client.Do(req, resp)

	if err != nil {
		return nil, err
	}

	statusCode := resp.StatusCode()

	if statusCode != 200 {
		return nil, fmt.Errorf("Server return status code %d", statusCode)
	}

	return resp.Body(), nil
}

// acquireRequest acquire new request with given params
func (api *API) acquireRequest(method, uri string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(api.url + uri)

	if method != "GET" {
		req.Header.SetMethod(method)
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/xml")
		req.Header.Add("Accept", "application/xml")
	}

	// Set auth header
	req.Header.Add("Authorization", "Basic "+api.basicAuth)

	return req
}

// getUserAgent generate user-agent string for client
func getUserAgent(app, version string) string {
	if app != "" && version != "" {
		return fmt.Sprintf(
			"%s/%s %s/%s (go; %s; %s-%s)",
			app, version, NAME, VERSION, runtime.Version(),
			runtime.GOARCH, runtime.GOOS,
		)
	}

	return fmt.Sprintf(
		"%s/%s (go; %s; %s-%s)",
		NAME, VERSION, runtime.Version(),
		runtime.GOARCH, runtime.GOOS,
	)
}

// genBasicAuthHeader generate basic auth header
func genBasicAuthHeader(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

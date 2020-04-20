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
	"net/url"
	"runtime"
	"strconv"
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
	data, err := api.doRequest("/stats")

	if err != nil {
		return nil, err
	}

	return parseStatsData(data)
}

// ListMounts fetches info about mounted sources
func (api *API) ListMounts() ([]*Mount, error) {
	data, err := api.doRequest("/listmounts")

	if err != nil {
		return nil, err
	}

	return parseMountsData(data)
}

// ListClients fetches list of listeners connected to given mount point
func (api *API) ListClients(mount string) ([]*Listener, error) {
	data, err := api.doRequest("/listclients?mount=" + mount)

	if err != nil {
		return nil, err
	}

	return parseClientListData(data)
}

// UpdateMeta updates meta for given mount source
func (api *API) UpdateMeta(mount, artist, title string) error {
	url := "/metadata?mode=updinfo&mount=" + mount + "&artist=" + esc(artist) + "&title=" + esc(title)
	data, err := api.doRequest(url)

	if err != nil {
		return err
	}

	return checkResponseData(data)
}

// MoveClients moves clients from one source to another
func (api *API) MoveClients(from, to string) error {
	url := "/moveclients?mount=" + from + "&destination=" + to
	data, err := api.doRequest(url)

	if err != nil {
		return err
	}

	return checkResponseData(data)
}

// KillClient kills client with given ID connected to given mount point
func (api *API) KillClient(mount string, id int) error {
	data, err := api.doRequest("/killclient?mount=" + mount + "&id=" + strconv.Itoa(id))

	if err != nil {
		return err
	}

	return checkResponseData(data)
}

// KillSource kills the source with given mount point
func (api *API) KillSource(mount string) error {
	data, err := api.doRequest("/killsource?mount=" + mount)

	if err != nil {
		return err
	}

	return checkResponseData(data)
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

// parseMountsData parses raw XML data with Icecast mounts info
func parseMountsData(data []byte) ([]*Mount, error) {
	mounts := &iceMounts{}
	err := xml.Unmarshal(data, mounts)

	if err != nil {
		return nil, err
	}

	return mounts.Mounts, nil
}

// parseClientListData parses raw XML data with Icecast listeners
func parseClientListData(data []byte) ([]*Listener, error) {
	listeners := &iceListeners{}
	err := xml.Unmarshal(data, listeners)

	if err != nil {
		return nil, err
	}

	return listeners.Listeners, nil
}

// checkResponseData checks default Icecast response for errors
func checkResponseData(data []byte) error {
	response := &iceResponse{}
	err := xml.Unmarshal(data, response)

	if err != nil {
		return err
	}

	if response.Return == 1 {
		return nil
	}

	return fmt.Errorf(response.Message)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// codebeat:disable[ARITY]

// doRequest creates and executes request
func (api *API) doRequest(uri string) ([]byte, error) {
	req := api.acquireRequest(uri)
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
func (api *API) acquireRequest(uri string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(api.url + "/admin" + uri)

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

// esc escapes the string so it can be safely placed inside a URL query
func esc(s string) string {
	return url.QueryEscape(s)
}

package icecast

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/essentialkaos/ek/v13/req"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// USER_AGENT is user-agent of API client
const USER_AGENT = "go-icecast/3"

// ////////////////////////////////////////////////////////////////////////////////// //

// API is Icecast API client
type API struct {
	engine   *req.Engine
	url      string
	user     string
	password string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrEmptyURL      = errors.New("URL is empty")
	ErrEmptyUser     = errors.New("Username is empty")
	ErrEmptyPassword = errors.New("Password is empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewAPI creates new API struct
func NewAPI(url, user, password string) (*API, error) {
	switch {
	case url == "":
		return nil, ErrEmptyURL
	case user == "":
		return nil, ErrEmptyUser
	case password == "":
		return nil, ErrEmptyPassword
	}

	engine := &req.Engine{}
	engine.SetUserAgent("go-icecast", "3")

	return &API{engine: engine, url: url, user: user, password: password}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// SetUserAgent set user-agent string based on app name and version
func (api *API) SetUserAgent(app, version string) {
	api.engine.SetUserAgent(app, version, USER_AGENT)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetStats fetches info about Icecast server
func (api *API) GetStats() (*Stats, error) {
	stats := &iceStats{}

	err := api.doRequest("/stats", nil, stats)

	if err != nil {
		return nil, err
	}

	return convertStats(stats), nil
}

// ListMounts fetches info about mounted sources
func (api *API) ListMounts() ([]*Mount, error) {
	mounts := &iceMounts{}

	err := api.doRequest("/listmounts", nil, mounts)

	if err != nil {
		return nil, err
	}

	return mounts.Mounts, nil
}

// ListClients fetches list of listeners connected to given mount point
func (api *API) ListClients(mount string) ([]*Listener, error) {
	listeners := &iceListeners{}

	err := api.doRequest("/listclients", req.Query{"mount": mount}, listeners)

	if err != nil {
		return nil, err
	}

	return listeners.Listeners, nil
}

// UpdateMeta updates meta for given mount source
func (api *API) UpdateMeta(mount string, meta TrackMeta) error {
	query := meta.ToQuery()
	query["mode"] = "updinfo"
	query["mount"] = mount

	response := &iceResponse{}

	err := api.doRequest("/metadata", query, response)

	if err != nil {
		return err
	}

	return parseResponse(response)
}

// UpdateFallback updates fallback for given mount source
func (api *API) UpdateFallback(mount, fallback string) error {
	response := &iceResponse{}

	err := api.doRequest(
		"/fallback",
		req.Query{
			"mount":    mount,
			"fallback": fallback,
		},
		response,
	)

	if err != nil {
		return err
	}

	return parseResponse(response)
}

// MoveClients moves clients from one source to another
func (api *API) MoveClients(mount, dest string) error {
	response := &iceResponse{}

	err := api.doRequest(
		"/moveclients",
		req.Query{
			"mount":       mount,
			"destination": dest,
		},
		response,
	)

	if err != nil {
		return err
	}

	return parseResponse(response)
}

// KillClient kills client with given ID connected to given mount point
func (api *API) KillClient(mount string, id int) error {
	response := &iceResponse{}

	err := api.doRequest(
		"/killclient",
		req.Query{
			"mount": mount,
			"id":    id,
		},
		response,
	)

	if err != nil {
		return err
	}

	return parseResponse(response)
}

// KillSource kills the source with given mount point
func (api *API) KillSource(mount string) error {
	response := &iceResponse{}

	err := api.doRequest("/killsource", req.Query{"mount": mount}, response)

	if err != nil {
		return err
	}

	return parseResponse(response)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// doRequest sends request to Icecast API
func (api *API) doRequest(endpoint string, query req.Query, response any) error {
	resp, err := api.engine.Get(req.Request{
		URL:    api.url + "/admin" + endpoint,
		Auth:   req.AuthBasic{api.user, api.password},
		Query:  query,
		Accept: req.CONTENT_TYPE_XML,
	})

	if err != nil {
		return fmt.Errorf("Can't send request to Icecast API: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned non-ok status code %d", resp.StatusCode)
	}

	xmlDec := xml.NewDecoder(resp.Body)
	err = xmlDec.Decode(response)

	if err != nil {
		return fmt.Errorf("Can't parse API response: %w", err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// parseResponse parses default Icecast response
func parseResponse(resp *iceResponse) error {
	if resp == nil {
		return fmt.Errorf("Response has no data")
	}

	if resp.Return != 1 {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

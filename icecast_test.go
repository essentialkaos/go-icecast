package icecast

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	. "pkg.re/essentialkaos/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _DEFAULT_PORT = "33100"
const _DEFAULT_USER = "source"
const _DEFAULT_PASS = "hackme"

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type IcecastSuite struct {
	client *API
}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&IcecastSuite{})

var statsError bool
var mountError bool

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *IcecastSuite) SetUpSuite(c *C) {
	var err error

	port := _DEFAULT_PORT

	if os.Getenv("TEST_PORT") != "" {
		port = os.Getenv("TEST_PORT")
	}

	url := fmt.Sprintf("http://127.0.0.1:%s", port)
	s.client, err = NewAPI(url, _DEFAULT_USER, _DEFAULT_PASS)

	if err != nil {
		c.Fatalf("Error while API client initialization: %v", err)
	}

	s.client.SetUserAgent("go-icecast-tester", "1.0.0")

	go runHTTPServer(c, port)

	time.Sleep(time.Second)
}

func (s *IcecastSuite) TestBasicErrors(c *C) {
	_, err := NewAPI("", "john", "pass")
	c.Assert(err, NotNil)

	_, err = NewAPI("http://domain.com", "", "pass")
	c.Assert(err, NotNil)

	_, err = NewAPI("http://domain.com", "john", "")
	c.Assert(err, NotNil)

	client, err := NewAPI("http://127.0.0.1:40000", "john", "pass")

	c.Assert(err, IsNil)
	c.Assert(client, NotNil)

	err = client.KillSource("/")

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestGetStats(c *C) {
	ic, err := s.client.GetStats()

	c.Assert(err, IsNil)
	c.Assert(ic, NotNil)

	c.Assert(ic.Admin, Equals, "icemaster@localhost")
	c.Assert(ic.Host, Equals, "localhost")
	c.Assert(ic.Location, Equals, "Earth")
	c.Assert(ic.Started.Unix(), Equals, int64(1587116898))

	c.Assert(ic.Info.ID, Equals, "Icecast 2.4.0-kh12")
	c.Assert(ic.Info.Build, Equals, 20190712000901)

	c.Assert(ic.Stats.BannedIPs, Equals, 2)
	c.Assert(ic.Stats.ClientConnections, Equals, 1027)
	c.Assert(ic.Stats.Clients, Equals, 1)
	c.Assert(ic.Stats.Connections, Equals, 741)
	c.Assert(ic.Stats.FileConnections, Equals, 139)
	c.Assert(ic.Stats.ListenerConnections, Equals, 25)
	c.Assert(ic.Stats.Listeners, Equals, 39)
	c.Assert(ic.Stats.OutgoingBitrate, Equals, 17869824)
	c.Assert(ic.Stats.SourceClientConnections, Equals, 3)
	c.Assert(ic.Stats.SourceRelayConnections, Equals, 4)
	c.Assert(ic.Stats.SourceTotalConnections, Equals, 6)
	c.Assert(ic.Stats.Sources, Equals, 6)
	c.Assert(ic.Stats.Stats, Equals, 2)
	c.Assert(ic.Stats.StatsConnections, Equals, 3)
	c.Assert(ic.Stats.StreamBytesRead, Equals, 259204096)
	c.Assert(ic.Stats.StreamBytesSent, Equals, 341397504)

	c.Assert(ic.Sources, HasLen, 2)

	ics := ic.GetSource("source1.ogg")

	c.Assert(ics, NotNil)

	c.Assert(ics.AudioInfo.Bitrate, Equals, 320000)
	c.Assert(ics.AudioInfo.Channels, Equals, 2)
	c.Assert(ics.AudioInfo.SampleRate, Equals, 48000)
	c.Assert(ics.AudioInfo.RawInfo, Equals, "ice-samplerate=48000;ice-bitrate=Quality 0;ice-channels=2")

	c.Assert(ics.IceAudioInfo.Bitrate, Equals, 320000)
	c.Assert(ics.IceAudioInfo.Channels, Equals, 2)
	c.Assert(ics.IceAudioInfo.SampleRate, Equals, 48000)
	c.Assert(ics.IceAudioInfo.RawInfo, Equals, "")

	c.Assert(ics.Track.Artist, Equals, "Nico & Vinz")
	c.Assert(ics.Track.Title, Equals, "Am I Wrong (Gryffin Remix) RA")
	c.Assert(ics.Track.RawInfo, Equals, "Nico & Vinz - Am I Wrong (Gryffin Remix) RA")

	c.Assert(ics.Info.Name, Equals, "Stream #1")
	c.Assert(ics.Info.Description, Equals, "My Super Stream")
	c.Assert(ics.Info.Type, Equals, "application/ogg")
	c.Assert(ics.Info.SubType, Equals, "Vorbis")
	c.Assert(ics.Info.URL, Equals, "https://domain.com")

	c.Assert(ics.Stats.Connected, Equals, 16)
	c.Assert(ics.Stats.IncomingBitrate, Equals, 320000)
	c.Assert(ics.Stats.OutgoingBitrate, Equals, 319042560)
	c.Assert(ics.Stats.ListenerConnections, Equals, 20)
	c.Assert(ics.Stats.ListenerPeak, Equals, 40)
	c.Assert(ics.Stats.Listeners, Equals, 16)
	c.Assert(ics.Stats.MaxListeners, Equals, -1)
	c.Assert(ics.Stats.QueueSize, Equals, 5)
	c.Assert(ics.Stats.SlowListeners, Equals, 5)
	c.Assert(ics.Stats.TotalBytesRead, Equals, 4655111)
	c.Assert(ics.Stats.TotalBytesSent, Equals, 1567151)

	c.Assert(ics.Bitrate, Equals, "Quality 0")
	c.Assert(ics.Genre, Equals, "Various Styles")
	c.Assert(ics.ListenURL, Equals, "http://localhost:8000/source.ogg")
	c.Assert(ics.MetadataUpdated.Unix(), Equals, int64(1587210604))
	c.Assert(ics.StreamStarted.Unix(), Equals, int64(1587210603))
	c.Assert(ics.Public, Equals, true)
	c.Assert(ics.SourceIP, Equals, "192.168.1.97")
	c.Assert(ics.UserAgent, Equals, "Native Instruments IceCast Uplink")

	ics = ic.GetSource("source1.aac")

	c.Assert(ics.AudioInfo.Bitrate, Equals, 320000)
	c.Assert(ics.AudioInfo.Channels, Equals, 1)
	c.Assert(ics.AudioInfo.SampleRate, Equals, 32000)
	c.Assert(ics.AudioInfo.CodecID, Equals, 10)

	c.Assert(ic.GetSource("/source1.ogg"), NotNil)

	ic = &Stats{}
	c.Assert(ic.GetSource("/source1.ogg"), IsNil)

	ic, err = s.client.GetStats()

	c.Assert(err, NotNil)
	c.Assert(ic, IsNil)
}

func (s *IcecastSuite) TestListClients(c *C) {
	listeners, err := s.client.ListClients("/source1.ogg")

	c.Assert(err, IsNil)
	c.Assert(listeners[0].ID, Equals, 757)
	c.Assert(listeners[0].IP, Equals, "192.168.1.22")
	c.Assert(listeners[0].UserAgent, Equals, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	c.Assert(listeners[0].Referer, Equals, "http://192.168.1.11:8000/source1.ogg")
	c.Assert(listeners[0].Lag, Equals, 0)
	c.Assert(listeners[0].Connected, Equals, 419)
	c.Assert(listeners[1].ID, Equals, 764)
	c.Assert(listeners[1].IP, Equals, "192.168.1.33")
	c.Assert(listeners[1].UserAgent, Equals, "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	c.Assert(listeners[1].Referer, Equals, "")
	c.Assert(listeners[1].Lag, Equals, 0)
	c.Assert(listeners[1].Connected, Equals, 312)

	listeners, err = s.client.ListClients("/source2.ogg")

	c.Assert(err, NotNil)
	c.Assert(listeners, HasLen, 0)
}

func (s *IcecastSuite) TestListMounts(c *C) {
	mounts, err := s.client.ListMounts()

	c.Assert(err, IsNil)
	c.Assert(mounts, HasLen, 1)
	c.Assert(mounts[0].Path, Equals, "/source1.ogg")
	c.Assert(mounts[0].Listeners, Equals, 48)
	c.Assert(mounts[0].Connected, Equals, 879)
	c.Assert(mounts[0].ContentType, Equals, "application/ogg")

	mounts, err = s.client.ListMounts()

	c.Assert(err, NotNil)
	c.Assert(mounts, HasLen, 0)
}

func (s *IcecastSuite) TestMetadata(c *C) {
	meta := TrackMeta{Artist: "Future Engineers", Title: "Source Code"}

	err := s.client.UpdateMeta("/source1.ogg", meta)

	c.Assert(err, IsNil)

	err = s.client.UpdateMeta("/source2.ogg", meta)

	c.Assert(err, NotNil)

	err = s.client.UpdateMeta("/source99.ogg", meta)

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestFallback(c *C) {
	err := s.client.UpdateFallback("/source1.ogg", "/source2.ogg")

	c.Assert(err, IsNil)

	err = s.client.UpdateFallback("/source2.ogg", "/source2.ogg")

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestMoveClients(c *C) {
	err := s.client.MoveClients("/source1.ogg", "/source2.ogg")

	c.Assert(err, IsNil)

	err = s.client.MoveClients("/source1.ogg", "/source3.ogg")

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestKillClient(c *C) {
	err := s.client.KillClient("/source1.ogg", 100)

	c.Assert(err, IsNil)

	err = s.client.KillClient("/source1.ogg", 101)

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestKillSource(c *C) {
	err := s.client.KillSource("/source1.ogg")

	c.Assert(err, IsNil)

	err = s.client.KillSource("/source2.ogg")

	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestUnmarshallingErrors(c *C) {
	_, err := parseStatsData([]byte("ABCD"))
	c.Assert(err, NotNil)

	_, err = parseMountsData([]byte("ABCD"))
	c.Assert(err, NotNil)

	_, err = parseClientListData([]byte("ABCD"))
	c.Assert(err, NotNil)

	err = checkResponseData([]byte("ABCD"))
	c.Assert(err, NotNil)
}

func (s *IcecastSuite) TestAux(c *C) {
	c.Assert(parseMax("unlimited"), Equals, -1)
	c.Assert(parseMax("1000"), Equals, 1000)
}

func (s *IcecastSuite) TestMetaEncoder(c *C) {
	meta := TrackMeta{
		Song:    "A",
		Title:   "CD",
		Artist:  "AB",
		URL:     "http://domain.com",
		Artwork: "http://domain.com/cover.jpg",
		Charset: "utf-8",
		Intro:   "intro.ogg",
	}

	c.Assert(meta.ToQuery(), Equals, `song=A&title=CD&artist=AB&url=http%3A%2F%2Fdomain.com&artwork=http%3A%2F%2Fdomain.com%2Fcover.jpg&charset=utf-8&intro=intro.ogg`)
	c.Assert(TrackMeta{}.ToQuery(), Equals, "song=Unknown")
}

// ////////////////////////////////////////////////////////////////////////////////// //

func runHTTPServer(c *C, port string) {
	server := &http.Server{
		Handler:        http.NewServeMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		c.Fatal(err.Error())
	}

	server.Handler.(*http.ServeMux).HandleFunc("/admin/metadata", handlerMetadata)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/fallback", handlerFallback)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/listclients", handlerListClients)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/moveclients", handlerMoveClients)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/killclient", handlerKillClient)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/killsource", handlerKillSource)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/stats", handlerStats)
	server.Handler.(*http.ServeMux).HandleFunc("/admin/listmounts", handlerListMounts)

	err = server.Serve(listener)

	if err != nil {
		c.Fatal(err.Error())
	}
}

func handlerMetadata(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mode := r.URL.Query().Get("mode")
	mount := r.URL.Query().Get("mount")
	artist := r.URL.Query().Get("artist")
	title := r.URL.Query().Get("title")

	if mount == "/source99.ogg" {
		w.WriteHeader(200)
		w.Write(getResponseData("metadata_error.xml"))
		return
	}

	switch {
	case mode != "updinfo",
		mount != "/source1.ogg",
		artist == "",
		title == "":
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	w.Write(getResponseData("metadata.xml"))
}

func handlerFallback(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mount := r.URL.Query().Get("mount")
	fallback := r.URL.Query().Get("fallback")

	if mount != "/source1.ogg" || fallback == "" {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
}

func handlerListClients(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mount := r.URL.Query().Get("mount")

	if mount != "/source1.ogg" {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	w.Write(getResponseData("listclients.xml"))
}

func handlerMoveClients(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mount := r.URL.Query().Get("mount")
	destination := r.URL.Query().Get("destination")

	switch {
	case mount != "/source1.ogg",
		destination != "/source2.ogg":
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	w.Write(getResponseData("moveclients.xml"))
}

func handlerKillClient(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mount := r.URL.Query().Get("mount")
	id := r.URL.Query().Get("id")

	switch {
	case mount != "/source1.ogg",
		id != "100":
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	w.Write(getResponseData("killclient.xml"))
}

func handlerKillSource(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	mount := r.URL.Query().Get("mount")

	if mount != "/source1.ogg" {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	w.Write(getResponseData("killsource.xml"))
}

func handlerStats(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	if statsError {
		w.WriteHeader(400)
		return
	}

	statsError = true

	w.WriteHeader(200)
	w.Write(getResponseData("stats.xml"))
}

func handlerListMounts(w http.ResponseWriter, r *http.Request) {
	if !isBasicAuthSet(r) {
		w.WriteHeader(403)
		return
	}

	if mountError {
		w.WriteHeader(400)
		return
	}

	mountError = true

	w.WriteHeader(200)
	w.Write(getResponseData("listmounts.xml"))
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isBasicAuthSet(r *http.Request) bool {
	user, pass, _ := r.BasicAuth()
	return user == _DEFAULT_USER && pass == _DEFAULT_PASS
}

func getResponseData(filename string) []byte {
	data, _ := ioutil.ReadFile("testdata/" + filename)
	return data
}

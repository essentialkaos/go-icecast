package icecast

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"io/ioutil"
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type IcecastSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&IcecastSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *IcecastSuite) TestBasic(c *C) {
	api, err := NewAPI("https://domain.com", "john", "qwerty1234")

	c.Assert(api, NotNil)
	c.Assert(err, IsNil)
}

func (s *IcecastSuite) TestStatsParser(c *C) {
	data, err := ioutil.ReadFile("testdata/stats.xml")

	c.Assert(err, IsNil)
	c.Assert(data, Not(HasLen), 0)

	ic, err := parseStatsData(data)

	c.Assert(err, IsNil)
	c.Assert(ic, NotNil)

	c.Assert(ic.Admin, Equals, "icemaster@localhost")
	c.Assert(ic.Host, Equals, "localhost")
	c.Assert(ic.Location, Equals, "Earth")
	c.Assert(ic.Start.Unix(), Equals, int64(1587116898))

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

	c.Assert(ic.Sources, HasLen, 1)

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
	c.Assert(ics.StreamStart.Unix(), Equals, int64(1587210603))
	c.Assert(ics.Public, Equals, true)
	c.Assert(ics.SourceIP, Equals, "192.168.1.97")
	c.Assert(ics.UserAgent, Equals, "Native Instruments IceCast Uplink")
}

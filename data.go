package icecast

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"strconv"
	"strings"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _DATE_FORMAT = "2/Jan/2006:15:04:05 -0700"

// ////////////////////////////////////////////////////////////////////////////////// //

// Stats contains info about Icecast Server
type Stats struct {
	Admin    string
	Host     string
	Started  time.Time
	Location string

	Info    *ServerInfo
	Stats   *ServerStats
	Sources Sources
}

// ServerInfo contains basic info about Icecast Server
type ServerInfo struct {
	ID    string
	Build int
}

// ServerStats contains overall Icecast Server statistics
type ServerStats struct {
	BannedIPs               int
	ClientConnections       int
	Clients                 int
	Connections             int
	FileConnections         int
	ListenerConnections     int
	Listeners               int
	OutgoingBitrate         int
	SourceClientConnections int
	SourceRelayConnections  int
	SourceTotalConnections  int
	Sources                 int
	Stats                   int
	StatsConnections        int
	StreamBytesRead         int
	StreamBytesSent         int
}

// Sources contains info about all sources
type Sources map[string]*Source

// Source contains info about source
type Source struct {
	MetadataUpdated time.Time
	StreamStarted   time.Time
	Bitrate         string
	Genre           string
	ListenURL       string
	SourceIP        string
	UserAgent       string
	AudioInfo       *AudioInfo
	IceAudioInfo    *AudioInfo
	Info            *SourceInfo
	Stats           *SourceStats
	Track           *TrackInfo
	Public          bool
}

// SourceInfo contains basic source info
type SourceInfo struct {
	Name        string
	Description string
	Type        string
	URL         string
	SubType     string
}

// SourceStats contains source statistics
type SourceStats struct {
	Connected           int
	IncomingBitrate     int
	OutgoingBitrate     int
	ListenerConnections int
	ListenerPeak        int
	Listeners           int
	MaxListeners        int
	QueueSize           int
	SlowListeners       int
	TotalBytesRead      int
	TotalBytesSent      int
}

// AudioInfo contains basic info about stream
type AudioInfo struct {
	Bitrate    int
	Channels   int
	SampleRate int
	CodecID    int
	RawInfo    string
}

// TrackInfo contains info about current playing track
type TrackInfo struct {
	Artist      string
	Title       string
	Artwork     string
	MetadataURL string
	RawInfo     string
}

// Mount contains basic info about source mount
type Mount struct {
	Path        string `xml:"mount,attr"`
	Listeners   int    `xml:"listeners"`
	Connected   int    `xml:"Connected"`
	ContentType string `xml:"content-type"`
}

// Listener contains info about listener
type Listener struct {
	ID        int    `xml:"ID"`
	IP        string `xml:"IP"`
	UserAgent string `xml:"UserAgent"`
	Referer   string `xml:"Referer"`
	Lag       int    `xml:"lag"`
	Connected int    `xml:"Connected"`
}

// TrackMeta contains track meta
type TrackMeta struct {
	Song    string
	Title   string
	Artist  string
	URL     string
	Artwork string
	Charset string
	Intro   string
}

// ////////////////////////////////////////////////////////////////////////////////// //

type iceStats struct {
	Admin                   string       `xml:"admin"`
	BannedIPs               int          `xml:"banned_IPs"`
	Build                   int          `xml:"build"`
	ClientConnections       int          `xml:"client_connections"`
	Clients                 int          `xml:"clients"`
	Connections             int          `xml:"connections"`
	FileConnections         int          `xml:"file_connections"`
	Host                    string       `xml:"host"`
	ListenerConnections     int          `xml:"listener_connections"`
	Listeners               int          `xml:"listeners"`
	Location                string       `xml:"location"`
	OutgoingKbitrate        int          `xml:"outgoing_kbitrate"`
	ServerID                string       `xml:"server_id"`
	ServerStart             string       `xml:"server_start"`
	SourceClientConnections int          `xml:"source_client_connections"`
	SourceRelayConnections  int          `xml:"source_relay_connections"`
	SourceTotalConnections  int          `xml:"source_total_connections"`
	Sources                 int          `xml:"sources"`
	Stats                   int          `xml:"stats"`
	StatsConnections        int          `xml:"stats_connections"`
	StreamKbytesRead        int          `xml:"stream_kbytes_read"`
	StreamKbytesSent        int          `xml:"stream_kbytes_sent"`
	SourcesData             []*iceSource `xml:"source"`
}

type iceSource struct {
	Mount               string `xml:"mount,attr"`
	Artist              string `xml:"artist"`
	Title               string `xml:"title"`
	Artwork             string `xml:"artwork"`
	AudioBitrate        int    `xml:"audio_bitrate"`
	AudioChannels       int    `xml:"audio_channels"`
	AudioInfo           string `xml:"audio_info"`
	AudioCodecID        int    `xml:"audio_codecid"`
	AudioSamplerate     int    `xml:"audio_samplerate"`
	MPEGChannels        int    `xml:"mpeg_channels"`
	MPEGSamplerate      int    `xml:"mpeg_samplerate"`
	Bitrate             string `xml:"bitrate"`
	Connected           int    `xml:"connected"`
	Genre               string `xml:"genre"`
	IceBitrate          int    `xml:"ice-bitrate"`
	IceChannels         int    `xml:"ice-channels"`
	IceSamplerate       int    `xml:"ice-samplerate"`
	IncomingBitrate     int    `xml:"incoming_bitrate"`
	ListenerConnections int    `xml:"listener_connections"`
	ListenerPeak        int    `xml:"listener_peak"`
	Listeners           int    `xml:"listeners"`
	ListenURL           string `xml:"listenurl"`
	MaxListeners        string `xml:"max_listeners"`
	MetadataUpdated     string `xml:"metadata_updated"`
	MetadataURL         string `xml:"metadata_url"`
	OutgoingKbitrate    int    `xml:"outgoing_kbitrate"`
	Public              int    `xml:"public"`
	QueueSize           int    `xml:"queue_size"`
	ServerDescription   string `xml:"server_description"`
	ServerName          string `xml:"server_name"`
	ServerType          string `xml:"server_type"`
	ServerURL           string `xml:"server_url"`
	SlowListeners       int    `xml:"slow_listeners"`
	SourceIP            string `xml:"source_ip"`
	StreamStart         string `xml:"stream_start"`
	Subtype             string `xml:"subtype"`
	TotalBytesRead      int    `xml:"total_bytes_read"`
	TotalBytesSent      int    `xml:"total_bytes_sent"`
	UserAgent           string `xml:"user_agent"`
	YpCurrentlyPlaying  string `xml:"yp_currently_playing"`
}

type iceMounts struct {
	Mounts []*Mount `xml:"source"`
}

type iceListeners struct {
	Listeners []*Listener `xml:"source>listener"`
}

type iceResponse struct {
	Message string `xml:"message"`
	Return  int    `xml:"return"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetSource tries to find source with given mount point
func (s *Stats) GetSource(mount string) *Source {
	if s.Sources == nil {
		return nil
	}

	source, ok := s.Sources[mount]

	if ok {
		return source
	}

	return s.Sources["/"+mount]
}

// ToQuery encodes meta for URL query
func (m TrackMeta) ToQuery() string {
	var result string

	if m.Song != "" {
		result += "song=" + esc(m.Song) + "&"
	}

	if m.Title != "" {
		result += "title=" + esc(m.Title) + "&"
	}

	if m.Artist != "" {
		result += "artist=" + esc(m.Artist) + "&"
	}

	if m.URL != "" {
		result += "url=" + esc(m.URL) + "&"
	}

	if m.Artwork != "" {
		result += "artwork=" + esc(m.Artwork) + "&"
	}

	if m.Charset != "" {
		result += "charset=" + esc(m.Charset) + "&"
	}

	if m.Intro != "" {
		result += "intro=" + esc(m.Intro) + "&"
	}

	if len(result) == 0 {
		return "song=Unknown"
	}

	return result[:len(result)-1]
}

// ////////////////////////////////////////////////////////////////////////////////// //

// convertStats converts data from internal format
func convertStats(sv *iceStats) *Stats {
	result := &Stats{
		Admin:    sv.Admin,
		Host:     sv.Host,
		Started:  parseDate(sv.ServerStart),
		Location: sv.Location,
		Info: &ServerInfo{
			ID:    sv.ServerID,
			Build: sv.Build,
		},
		Stats: &ServerStats{
			BannedIPs:               sv.BannedIPs,
			ClientConnections:       sv.ClientConnections,
			Clients:                 sv.Clients,
			Connections:             sv.Connections,
			FileConnections:         sv.FileConnections,
			ListenerConnections:     sv.ListenerConnections,
			Listeners:               sv.Listeners,
			OutgoingBitrate:         sv.OutgoingKbitrate * 1024,
			SourceClientConnections: sv.SourceClientConnections,
			SourceRelayConnections:  sv.SourceRelayConnections,
			SourceTotalConnections:  sv.SourceTotalConnections,
			Sources:                 sv.Sources,
			Stats:                   sv.Stats,
			StatsConnections:        sv.StatsConnections,
			StreamBytesRead:         sv.StreamKbytesRead * 1024,
			StreamBytesSent:         sv.StreamKbytesSent * 1024,
		},
	}

	if len(sv.SourcesData) != 0 {
		result.Sources = make(Sources)
	}

	for _, s := range sv.SourcesData {
		result.Sources[s.Mount] = &Source{
			AudioInfo: &AudioInfo{
				Bitrate:    s.AudioBitrate,
				Channels:   s.AudioChannels,
				SampleRate: s.AudioSamplerate,
				CodecID:    s.AudioCodecID,
				RawInfo:    s.AudioInfo,
			},
			IceAudioInfo: &AudioInfo{
				Bitrate:    s.IceBitrate * 1000,
				Channels:   s.IceChannels,
				SampleRate: s.IceSamplerate,
			},
			Track: &TrackInfo{
				Artist:      s.Artist,
				Title:       s.Title,
				Artwork:     s.Artwork,
				MetadataURL: s.MetadataURL,
				RawInfo:     s.YpCurrentlyPlaying,
			},
			Info: &SourceInfo{
				Name:        s.ServerName,
				Description: s.ServerDescription,
				Type:        s.ServerType,
				URL:         s.ServerURL,
				SubType:     s.Subtype,
			},
			Stats: &SourceStats{
				Connected:           s.Connected,
				IncomingBitrate:     s.IncomingBitrate,
				OutgoingBitrate:     s.OutgoingKbitrate * 1024,
				ListenerConnections: s.ListenerConnections,
				ListenerPeak:        s.ListenerPeak,
				Listeners:           s.Listeners,
				MaxListeners:        parseMax(s.MaxListeners),
				QueueSize:           s.QueueSize,
				SlowListeners:       s.SlowListeners,
				TotalBytesRead:      s.TotalBytesRead,
				TotalBytesSent:      s.TotalBytesSent,
			},
			Bitrate:         s.Bitrate,
			Genre:           s.Genre,
			ListenURL:       s.ListenURL,
			MetadataUpdated: parseDate(s.MetadataUpdated),
			StreamStarted:   parseDate(s.StreamStart),
			Public:          s.Public == 1,
			SourceIP:        s.SourceIP,
			UserAgent:       s.UserAgent,
		}

		if s.MPEGSamplerate != 0 && s.AudioSamplerate == 0 {
			result.Sources[s.Mount].AudioInfo.SampleRate = s.MPEGSamplerate
		}

		if s.MPEGChannels != 0 && s.AudioChannels == 0 {
			result.Sources[s.Mount].AudioInfo.Channels = s.MPEGChannels
		}

		if s.AudioBitrate == 0 && isNum(s.Bitrate) {
			bitrate, _ := strconv.Atoi(s.Bitrate)

			if bitrate > 0 && bitrate < 1024 {
				bitrate *= 1000
			}

			result.Sources[s.Mount].AudioInfo.Bitrate = bitrate
		}
	}

	return result
}

// parseMax parse value with possible "unlimited" value
func parseMax(data string) int {
	if data == "unlimited" {
		return -1
	}

	n, _ := strconv.Atoi(data)

	return n
}

// isNum returns true if given string contains number
func isNum(s string) bool {
	return strings.Trim(s, "0123456789") == ""
}

// parseDate parses date
func parseDate(date string) time.Time {
	result, _ := time.Parse(_DATE_FORMAT, date)
	return result
}

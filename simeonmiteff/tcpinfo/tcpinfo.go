package tcpinfo

import (
	"encoding/json"
	"time"
)

type Info struct {
	State               string        `json:"state"`                  // Connection state
	Options             []Option      `json:"options,omitempty"`      // Requesting options
	PeerOptions         []Option      `json:"peer_options,omitempty"` // Options requested from peer
	SenderMSS           uint64        `json:"snd_mss"`                // Maximum segment size for sender in bytes
	ReceiverMSS         uint64        `json:"rcv_mss"`                // Maximum segment size for receiver in bytes
	RTT                 time.Duration `json:"rtt"`                    // Round-trip time
	RTTVar              time.Duration `json:"rttvar"`                 // Round-trip time variation
	RTO                 time.Duration `json:"rto"`                    // Retransmission timeout
	ATO                 time.Duration `json:"ato"`                    // Delayed acknowledgement timeout [Linux only]
	LastDataSent        time.Duration `json:"last_data_sent"`         // Since last data sent [Linux only]
	LastDataReceived    time.Duration `json:"last_data_rcvd"`         // Since last data received [FreeBSD and Linux]
	LastAckReceived     time.Duration `json:"last_ack_rcvd"`          // Since last ack received [Linux only]
	ReceiverWindow      uint64        `json:"rcv_wnd"`                // advertised receiver window in bytes
	SenderSSThreshold   uint64        `json:"snd_ssthresh"`           // slow start threshold for sender in bytes or # of segments
	ReceiverSSThreshold uint64        `json:"rcv_ssthresh"`           // slow start threshold for receiver in bytes [Linux only]
	SenderWindowBytes   uint64        `json:"snd_cwnd_bytes"`         // congestion window for sender in bytes [Darwin and FreeBSD]
	SenderWindowSegs    uint64        `json:"snd_cwnd_segs"`          // congestion window for sender in # of segments [Linux and NetBSD]
	Sys                 *SysInfo      `json:"sys,omitempty"`          // Platform-specific information
}

type Option struct {
	Kind  string `json:"kind"`
	Value uint64 `json:"value"`
}

// MarshalJSON implements the MarshalJSON method of json.Marshaler
// interface.
func (i *Info) MarshalJSON() ([]byte, error) {
	raw := make(map[string]interface{})
	raw["state"] = i.State
	if len(i.Options) > 0 {
		opts := make(map[string]interface{})
		for _, opt := range i.Options {
			opts[opt.Kind] = opt
		}
		raw["options"] = opts
	}
	if len(i.PeerOptions) > 0 {
		opts := make(map[string]interface{})
		for _, opt := range i.PeerOptions {
			opts[opt.Kind] = opt
		}
		raw["peer_options"] = opts
	}
	raw["snd_mss"] = i.SenderMSS
	raw["rcv_mss"] = i.ReceiverMSS
	raw["rtt"] = i.RTT
	raw["rttvar"] = i.RTTVar
	raw["rto"] = i.RTO
	raw["ato"] = i.ATO
	raw["last_data_sent"] = i.LastDataSent
	raw["last_data_rcvd"] = i.LastDataReceived
	raw["last_ack_rcvd"] = i.LastAckReceived
	raw["rcv_wnd"] = i.ReceiverWindow
	raw["snd_ssthresh"] = i.SenderSSThreshold
	raw["rcv_ssthresh"] = i.ReceiverSSThreshold
	raw["snd_cwnd_bytes"] = i.SenderWindowBytes
	raw["snd_cwnd_segs"] = i.SenderWindowSegs
	if i.Sys != nil {
		raw["sys"] = i.Sys
	}
	return json.Marshal(&raw)
}

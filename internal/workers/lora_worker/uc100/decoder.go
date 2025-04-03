package uc100

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/johandrevandeventer/lora-worker/internal/workers/types"
)

type UC100 struct {
	EndDeviceIDs  EndDeviceIDs  `json:"end_device_ids"`
	ReceivedAt    time.Time     `json:"received_at"`
	UplinkMessage UplinkMessage `json:"uplink_message"`
}

// Device identifiers
type EndDeviceIDs struct {
	DeviceID       string         `json:"device_id"`
	ApplicationIDs ApplicationIDs `json:"application_ids"`
	DevEUI         string         `json:"dev_eui"`
	JoinEUI        string         `json:"join_eui"`
	DevAddr        string         `json:"dev_addr"`
}

type ApplicationIDs struct {
	ApplicationID string `json:"application_id"`
}

// Uplink message structure
type UplinkMessage struct {
	SessionKeyID    string         `json:"session_key_id"`
	FPort           int            `json:"f_port"`
	FCnt            int            `json:"f_cnt"`
	FrmPayload      string         `json:"frm_payload"`
	DecodedPayload  map[string]any `json:"decoded_payload"`
	RxMetadata      []RxMetadata   `json:"rx_metadata"`
	Settings        Settings       `json:"settings"`
	ReceivedAt      time.Time      `json:"received_at"`
	Confirmed       bool           `json:"confirmed"`
	ConsumedAirtime string         `json:"consumed_airtime"`
	NetworkIDs      NetworkIDs     `json:"network_ids"`
}

// Metadata for received packets
type RxMetadata struct {
	GatewayIDs   GatewayIDs `json:"gateway_ids"`
	Time         time.Time  `json:"time"`
	Timestamp    int64      `json:"timestamp"`
	RSSI         int        `json:"rssi"`
	ChannelRSSI  int        `json:"channel_rssi"`
	SNR          float64    `json:"snr"`
	UplinkToken  string     `json:"uplink_token"`
	ChannelIndex int        `json:"channel_index"`
}

type GatewayIDs struct {
	GatewayID string `json:"gateway_id"`
	EUI       string `json:"eui"`
}

type Settings struct {
	DataRate   DataRate  `json:"data_rate"`
	CodingRate string    `json:"coding_rate"`
	Frequency  string    `json:"frequency"`
	Timestamp  int64     `json:"timestamp"`
	Time       time.Time `json:"time"`
}

type DataRate struct {
	Lora Lora `json:"lora"`
}

type Lora struct {
	Bandwidth       int `json:"bandwidth"`
	SpreadingFactor int `json:"spreading_factor"`
}

type NetworkIDs struct {
	NetID string `json:"net_id"`
}

// Decoder processes Lora payloads
func Decoder(payload json.RawMessage) (decodedPayloadInfo *types.DecodedPayloadInfo, err error) {
	var data UC100
	if err := json.Unmarshal(payload, &data); err != nil {
		return decodedPayloadInfo, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &types.DecodedPayloadInfo{
		RawPayload: payload,
	}, nil
}

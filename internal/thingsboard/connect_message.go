package thingsboard

import "encoding/json"

const ConnectTopic = "v1/gateway/connect"

type ConnectMessage struct {
	DeviceName string `json:"device"`
}

func (m ConnectMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

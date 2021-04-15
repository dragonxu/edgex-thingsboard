package thingsboard

type TelemetryMessage map[string][]DeviceData

type DeviceData struct {
	TS     int                    `json:"ts"`
	Values map[string]interface{} `json:"values"`
}

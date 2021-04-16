package thingsboard

type TelemetryMessage map[string][]DeviceDataPoint // deviceName -> []DeviceDataPoint

type DeviceDataPoint struct {
	TS     int                    `json:"ts"`
	Values map[string]interface{} `json:"values"`
}

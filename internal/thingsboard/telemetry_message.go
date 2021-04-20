package thingsboard

import (
	"encoding/json"
	"errors"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"strconv"
)

const TelemetryTopic = "v1/gateway/telemetry"

type TelemetryMessage map[string][]DeviceDataPoint // deviceName -> []DeviceDataPoint

func (m TelemetryMessage) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

type DeviceDataPoint struct {
	TS     int                    `json:"ts"`
	Values map[string]interface{} `json:"values"`
}

func AdapterTelemetryMessage(event *contract.Event) (TelemetryMessage, error) {
	msgs := make(TelemetryMessage)
	for _, reading := range event.Readings {
		value, err := adapterTelemetryValue(reading.ValueType, reading.Value)
		if err != nil {
			return nil, err
		}

		ts := int(reading.Origin / 1000000)
		dataPoint := DeviceDataPoint{
			TS: ts,
			Values: map[string]interface{}{
				reading.Name: value,
			},
		}
		if devicePoints, ok := msgs[event.Device]; ok {
			devicePoints = append(devicePoints, dataPoint)
		} else {
			msgs[reading.Device] = []DeviceDataPoint{dataPoint}
		}
	}
	return msgs, nil
}

func adapterTelemetryValue(valueType string, value string) (v interface{}, err error) {
	switch valueType {
	case contract.ValueTypeBool:
		return strconv.ParseBool("value")
	case contract.ValueTypeString:
		return value, nil
	case contract.ValueTypeUint8:
		n, err := strconv.Atoi(value)
		return uint8(n), err
	case contract.ValueTypeUint16:
		n, err := strconv.Atoi(value)
		return uint16(n), err
	case contract.ValueTypeUint32:
		n, err := strconv.Atoi(value)
		return uint32(n), err
	case contract.ValueTypeUint64:
		n, err := strconv.Atoi(value)
		return uint64(n), err
	case contract.ValueTypeInt8:
		n, err := strconv.Atoi(value)
		return int8(n), err
	case contract.ValueTypeInt16:
		n, err := strconv.Atoi(value)
		return int16(n), err
	case contract.ValueTypeInt32:
		n, err := strconv.Atoi(value)
		return int32(n), err
	case contract.ValueTypeInt64:
		n, err := strconv.Atoi(value)
		return int64(n), err
	case contract.ValueTypeFloat32:
		n, err := strconv.ParseFloat(value, 32)
		return float32(n), err
	case contract.ValueTypeFloat64:
		n, err := strconv.ParseFloat(value, 32)
		return n, err
	case contract.ValueTypeBinary:
		return value, nil
	case contract.ValueTypeBoolArray:
		var arr []bool
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeStringArray:
		var arr []string
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeUint8Array:
		var arr []uint8
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeUint16Array:
		var arr []uint16
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeUint32Array:
		var arr []uint32
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeUint64Array:
		var arr []uint64
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeInt8Array:
		var arr []int8
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeInt16Array:
		var arr []int16
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeInt32Array:
		var arr []int32
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeInt64Array:
		var arr []int64
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeFloat32Array:
		var arr []float32
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	case contract.ValueTypeFloat64Array:
		var arr []float64
		err := json.Unmarshal([]byte(value), &arr)
		return value, err
	default:
		return value, errors.New("unsupported value type")
	}
}

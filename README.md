# Edgex-thingsboard

[中文](README.cn.md)

A micro service that connects Edgex to Thingsboard by MQTT. 

- Connect Edgex devices to Thingsboard
- Report Edgex devices events to Thingsboard
- Handle RPC requests from Thingsboard

## Usage

Before start service, you need to configure Thingsboard MQTT client.

Configure by file:

```
[Mqtt]
Address = "tcp://localhost:1883"
Username = "edgex-thingsboard"
ClientId = "client-id"
Timeout = 10000
```

Or configure by environment variables:

```
MQTT_ADDRESS: tcp://localhost:1883
MQTT_USERNAME: edgex-thingsboard
MQTT_CLIENTID: client-id
MQTT_TIMEOUT: "10000"
```

While: 

|Arguments|Name|Description|
|---|---|---|
| Mqtt.Address | MQTT Address | |
| Mqtt.Username | Username　| |
| Mqtt.ClientId | Client ID | |
| Mqtt.Timeout | Timeout | unit: millisecond |

## Build

If you are using zeroMQ as your message bus be sure to first [install the zeroMQ library](https://github.com/edgexfoundry/edgex-go#zeromq).

## Internal

### Connecting Devices

1. Edgex-thingsboard will send MQTT messages to Thingsboard while starting:

```json
{
  "device": "Virtual-Sensor-01"
}
```

While:

|Arguments|Name|Description|
|---|---|---|
| device | Device Name ||

### RPC

1. Thingsboard will send MQTT messages to Edgex-thingsboard:

```json
{
  "device": "Virtual-Sensor-01",
  "data": {
    "id": 4,
    "method": "GET",
    "service": "edgex-core-command",
    "uri": "/api/version",
    "params": {},
    "timeout": 10000
  }
}
```

While:

|Arguments|Name|Description|
|---|---|---|
| device | Device Name ||
| data.id | Request ID ||
| data.service | Service Name ||
| data.uri | HTTP URI ||
| data.method | HTTP Method ||
| data.params | HTTP Parameters ||
| data.timeout | HTTP Timeout | unit: millisecond |

|Service Key|Service Name|Api Doc|
|---|---|---|
| edgex-core-command | Core Command Service | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-core-data | Core Data Service | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-core-metadata | Meta Data Service| https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-support-notifications | Notification Service | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-support-scheduler | Scheduler Service | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/1.2.1 |
| edgex-sys-mgmt-agent | System Management Service | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/1.2.1 |

2. After Edgex-thingsboard processed the request, it will reply MQTT messages to Thingsboard:

```json
{
  "device": "Virtual-Sensor-01",
  "id": 4,
  "data": {
    "http_status": 200,
    "success": true,
    "message": "",
    "result": {}
  }
}
```

While:

|Arguments|Name|Description|
|---|---|---|
| id | Request ID ||
| device | Device Name ||
| data.http_status | HTTP Status ||
| data.success | Response Status ||
| data.message | Response Error Message ||
| data.result | Response Data ||

### Telemetry

Edgex-thingsboard will report devices events to Thingsboard:

```json
{
  "Device A": [{
    "ts": 1483228800000,
    "values": {
      "temperature": 42,
      "humidity": 80
    }
  }, {
    "ts": 1483228801000,
    "values": {
      "temperature": 43,
      "humidity": 82
    }
  }],
  "Device B": [{
    "ts": 1483228800000,
    "values": {
      "temperature": 42,
      "humidity": 80
    }
  }]
}
```
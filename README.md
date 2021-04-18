# edgex-thingsboard

用于将Edgex网关接入Thingsboard物联网平台。主要包括基于MQTT的遥测数据接入，和控制RPC接入。

## 使用方式

启动服务时，需配置Thingsboard服务端MQTT的连接信息。

如使用配置文件方式：

```
[Mqtt]
Address = "tcp://localhost:1883"
Username = "edgex-thingsboard"
RpcRequestTopic = "v1/gateway/rpc"
RpcResponseTopic = "v1/gateway/rpc"
TelemetryTopic = "v1/gateway/telemetry"
Timeout = 10000
```

或使用环境变量方式：

```
MQTT_ADDRESS: tcp://192.168.5.88:1883
MQTT_USERNAME: edgex-control-agent
MQTT_RPCREQUESTTOPIC = "v1/gateway/rpc"
MQTT_RPCRESPONSETOPIC = "v1/gateway/rpc"
MQTT_TELEMETRYTOPIC = "v1/gateway/telemetry"
MQTT_TIMEOUT: "10000"
```

其中：

|参数|名称|描述|
|---|---|---|
| Mqtt.Address | MQTT Broker 地址| |
| Mqtt.Username | 用户名　| |
| Mqtt.RpcRequestTopic | 控制请求主题 | |
| Mqtt.RpcResponseTopic | 控制响应主题 | |
| Mqtt.TelemetryTopic | 遥测数据主题 | |
| Mqtt.Timeout | 超时时间 | 单位为毫秒 |

## 实现原理

### 控制RPC

1. Thingsboard会按如下格式发送MQTT消息给Edgex：

发送消息： 

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

其中： 

|参数|名称|描述|
|---|---|---|
| device | 设备名称 ||
| data.id | 请求ID ||
| data.service | 微服务名称 ||
| data.uri | HTTP接口地址 ||
| data.method | HTTP请求方法 ||
| data.params | HTTP请求参数 ||
| data.timeout | HTTP请求超时时间 | 单位为毫秒 |

|service值|对应微服务名称|微服务接口地址|
|---|---|---|
| edgex-core-command | 命令微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-core-data | 核心数据微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-core-metadata | 元数据微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-support-notifications | 通知微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command |
| edgex-support-scheduler | 调度微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/1.2.1 |
| edgex-sys-mgmt-agent | 系统管理微服务 | https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/1.2.1 |

2. Edgex处理完RPC消息后，会返回如下MQTT消息给Thingsboard：

```json
{
  "device": "Virtual-Sensor-01",
  "id": 4,
  "data": {
    "success": true
  }
}
```

其中：

|参数|名称|描述|
|---|---|---|
| id | 请求ID ||
| device | 设备名称 ||
| data.success | 响应结果 ||


### 遥测数据

Edgex会将遥测数据按如下格式发往给Thingsboard：

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
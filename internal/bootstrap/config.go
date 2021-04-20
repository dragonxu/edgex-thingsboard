package bootstrap

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	bootstrapConfig "github.com/edgexfoundry/go-mod-bootstrap/config"
	msgTypes "github.com/edgexfoundry/go-mod-messaging/pkg/types"
)

type WritableInfo struct {
	LogLevel string
}

type MessageQueueInfo struct {
	Protocol string
	Host     string
	Port     int
	Type     string
	Topic    string
	Optional map[string]string
}

func (m MessageQueueInfo) URL() string {
	return fmt.Sprintf("%s://%s:%v", m.Protocol, m.Host, m.Port)
}

type ThingsboardMQTTInfo struct {
	Address  string
	Username string
	ClientId string
	Timeout  int
}

func (t ThingsboardMQTTInfo) GetMQTTOption() *mqtt.ClientOptions {
	return mqtt.NewClientOptions().
		AddBroker(t.Address).
		SetUsername(t.Username).
		SetClientID(t.ClientId)
}

type ConfigurationClients map[string]bootstrapConfig.ClientInfo

type ConfigurationStruct struct {
	Writable        WritableInfo
	Service         bootstrapConfig.ServiceInfo
	MessageQueue    MessageQueueInfo
	Registry        bootstrapConfig.RegistryInfo
	ThingsBoardMQTT ThingsboardMQTTInfo
	Clients         ConfigurationClients
}

func (c *ConfigurationStruct) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ConfigurationStruct)
	if ok {
		if configuration.Service.Port == 0 {
			return false
		}
		*c = *configuration
	}
	return ok
}

func (c *ConfigurationStruct) EmptyWritablePtr() interface{} {
	return &WritableInfo{}
}

func (c *ConfigurationStruct) UpdateWritableFromRaw(rawWritable interface{}) bool {
	writable, ok := rawWritable.(*WritableInfo)
	if ok {
		c.Writable = *writable
	}
	return ok
}

func (c *ConfigurationStruct) GetBootstrap() bootstrapConfig.BootstrapConfiguration {
	return bootstrapConfig.BootstrapConfiguration{
		Clients:  c.Clients,
		Service:  c.Service,
		Registry: c.Registry,
	}
}

func (c *ConfigurationStruct) GetLogLevel() string {
	return c.Writable.LogLevel
}

func (c *ConfigurationStruct) GetRegistryInfo() bootstrapConfig.RegistryInfo {
	return c.Registry
}

func (c ConfigurationStruct) GetMessagingConfig() msgTypes.MessageBusConfig {
	return msgTypes.MessageBusConfig{
		Type: c.MessageQueue.Type,
		SubscribeHost: msgTypes.HostInfo{
			Host:     c.MessageQueue.Host,
			Port:     c.MessageQueue.Port,
			Protocol: c.MessageQueue.Protocol,
		},
		Optional: c.MessageQueue.Optional,
	}
}

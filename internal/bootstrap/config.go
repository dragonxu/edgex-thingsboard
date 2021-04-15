package bootstrap

import (
	bootstrapConfig "github.com/edgexfoundry/go-mod-bootstrap/config"
)

type ConfigurationClients map[string]bootstrapConfig.ClientInfo

type WritableInfo struct {
	LogLevel string
}

type MQTTInfo struct {
	Address          string
	Username         string
	RPCRequestTopic  string
	RPCResponseTopic string
	TelemetryTopic   string
	Timeout          int
}

type ConfigurationStruct struct {
	Writable WritableInfo
	Service  bootstrapConfig.ServiceInfo
	Mqtt     MQTTInfo
	Registry bootstrapConfig.RegistryInfo
	Clients  ConfigurationClients
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

func (c *ConfigurationStruct) GetMQTTInfo() MQTTInfo {
	return c.Mqtt
}

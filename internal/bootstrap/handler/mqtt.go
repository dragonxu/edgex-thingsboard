package handler

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/interfaces"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"sync"
	"time"
)

type MQTT struct {
	mqttInfo       interfaces.MQTTOption
	connectTimeout time.Duration
}

func NewMQTT(mqttInfo interfaces.MQTTOption, connectTimeout time.Duration) *MQTT {
	return &MQTT{
		mqttInfo:       mqttInfo,
		connectTimeout: connectTimeout,
	}
}

func (p *MQTT) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, startupTimer startup.Timer, dic *di.Container) bool {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	var client mqtt.Client
	for startupTimer.HasNotElapsed() {
		var err error
		opt := p.mqttInfo.GetMQTTOption()
		client, err = mqtt_server.NewClient(opt, p.connectTimeout)
		if err == nil {
			break
		}
		client = nil
		lc.Warn(fmt.Sprintf("couldn't create mqtt client: %s", err.Error()))
		startupTimer.SleepForInterval()
	}

	if client == nil {
		lc.Error(fmt.Sprintf("failed to create mqtt client in allotted time"))
		return false
	}

	dic.Update(di.ServiceConstructorMap{
		container.MQTTName: func(get di.Get) interface{} {
			return client
		},
	})

	lc.Info("mqtt client connected")
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		for {
			select {
			case <-ctx.Done():
				client.Disconnect(uint(p.connectTimeout / time.Millisecond))
				lc.Info("mqtt client disconnected")
				return
			}
		}
	}()

	return true
}

package pubsub

import (
	"context"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/interfaces"
	"github.com/inspii/edgex-thingsboard/pkg/messaging"
	"github.com/inspii/edgex-thingsboard/pkg/messaging/mqtt"
	"io"
	"sync"
	"time"
)

type httpServer interface {
	IsRunning() bool
}

type PubSub struct {
	httpServer httpServer
	mqttInfo   interfaces.MQTTInfo
}

func NewPubSub(httpServer httpServer, mqttInfo interfaces.MQTTInfo) PubSub {
	return PubSub{
		httpServer: httpServer,
		mqttInfo:   mqttInfo,
	}
}

func (p PubSub) newPubSub(lc logger.LoggingClient) (messaging.PubSub, error) {
	info := p.mqttInfo.GetMQTTInfo()
	timeout := time.Duration(info.Timeout) * time.Millisecond
	client, err := mqtt.NewClient(info.Address, info.Username, timeout, lc)
	if err != nil {
		return nil, err
	}
	return mqtt.NewPubSub(client), nil
}

func (p PubSub) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, startupTimer startup.Timer, dic *di.Container) bool {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	// initialize database.
	var pubsub messaging.PubSub
	for startupTimer.HasNotElapsed() {
		var err error
		pubsub, err = p.newPubSub(lc)
		if err == nil {
			break
		}
		pubsub = nil
		lc.Warn(fmt.Sprintf("couldn't create mqtt client: %s", err.Error()))
		startupTimer.SleepForInterval()
	}

	if pubsub == nil {
		lc.Error(fmt.Sprintf("failed to create mqtt client in allotted time"))
		return false
	}

	dic.Update(di.ServiceConstructorMap{
		container.PubSubName: func(get di.Get) interface{} {
			return pubsub
		},
	})

	lc.Info("pubsub connected")
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		for {
			if p.httpServer.IsRunning() == false {
				if closer, ok := pubsub.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						lc.Error(fmt.Sprintf("failed to close pubsub: %s", err.Error()))
					}
				}
				break
			}
			time.Sleep(time.Second)
		}
		lc.Info("pubsub disconnected")
	}()

	return true
}

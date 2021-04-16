package handler

import (
	"context"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-messaging/messaging"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/interfaces"
	"sync"
)

type Messaging struct {
	messagingInfo interfaces.MessagingInfo
}

func NewMessaging(mqttInfo interfaces.MessagingInfo) Messaging {
	return Messaging{
		messagingInfo: mqttInfo,
	}
}

func (p Messaging) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, startupTimer startup.Timer, dic *di.Container) bool {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	var client messaging.MessageClient
	for startupTimer.HasNotElapsed() {
		var err error
		config := p.messagingInfo.GetMessagingConfig()
		client, err = messaging.NewMessageClient(config)
		if err == nil {
			break
		}
		client = nil
		lc.Warn(fmt.Sprintf("couldn't create messaging client: %s", err.Error()))
		startupTimer.SleepForInterval()
	}

	if client == nil {
		lc.Error(fmt.Sprintf("failed to create messaging client in allotted time"))
		return false
	}

	dic.Update(di.ServiceConstructorMap{
		container.MessagingName: func(get di.Get) interface{} {
			return client
		},
	})

	lc.Info("messaging client connected")
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		for {
			select {
			case <-ctx.Done():
				if err := client.Disconnect(); err != nil {
					lc.Error("failed to disconnect messaging client")
					return
				}
				lc.Info("messaging client disconnected")
				return
			}
		}
	}()

	return true
}

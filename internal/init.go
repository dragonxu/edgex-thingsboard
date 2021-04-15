package internal

import (
	"context"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	contracts "github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/gorilla/mux"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"sync"
)

type Bootstrap struct {
	router *mux.Router
}

func NewBootstrap(router *mux.Router) *Bootstrap {
	return &Bootstrap{
		router: router,
	}
}

func (b *Bootstrap) BootstrapHandler(ctx context.Context, _ *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	loadRoutes(b.router, dic)

	configuration := container.ConfigurationFrom(dic.Get)

	dic.Update(di.ServiceConstructorMap{
		container.ClientsName: func(get di.Get) interface{} {
			return bootstrap.NewServiceRoutes()
		},
	})
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	if err := mqttForwardHTTP(dic); err != nil {
		lc.Error(fmt.Sprintf("pubsub forward http: %s", err.Error()))
		return false
	}

	serviceClients := container.ClientsFrom(dic.Get)
	for serviceKey, serviceName := range b.listDefaultServices() {
		serviceClient := configuration.Clients[serviceName]
		serviceAddr := fmt.Sprintf("%s://%s:%d", serviceClient.Protocol, serviceClient.Host, serviceClient.Port)
		serviceClients.Set(serviceKey, serviceAddr)
	}

	return true
}

func (Bootstrap) listDefaultServices() map[string]string {
	return map[string]string{
		contracts.CoreCommandServiceKey:           "CoreCommand",
		contracts.CoreDataServiceKey:              "CoreData",
		contracts.CoreMetaDataServiceKey:          "CoreMetadata",
		contracts.SupportSchedulerServiceKey:      "SupportScheduler",
		contracts.SupportNotificationsServiceKey:  "SupportNotification",
		contracts.SystemManagementAgentServiceKey: "SystemMgmtAgent",
	}
}

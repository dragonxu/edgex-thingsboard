package internal

import (
	"context"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	contracts "github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/gorilla/mux"
	bootstrap2 "github.com/inspii/edgex-thingsboard/internal/bootstrap"
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
	serviceRoutes := bootstrap2.NewServiceRoutes()
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	for serviceKey, serviceName := range b.listDefaultServices() {
		serviceClient := configuration.Clients[serviceName]
		serviceAddr := fmt.Sprintf("%s://%s:%d", serviceClient.Protocol, serviceClient.Host, serviceClient.Port)
		serviceRoutes.Set(serviceKey, serviceAddr)
	}

	dic.Update(di.ServiceConstructorMap{
		container.ServiceRoutesName: func(get di.Get) interface{} {
			return bootstrap2.NewServiceRoutes()
		},
		container.ServiceRoutesName: func(get di.Get) interface{} {
			return serviceRoutes
		},
		container.MetaDataDeviceClientName: func(get di.Get) interface{} {
			metadataURL, ok := serviceRoutes.Get(contracts.CoreMetaDataServiceKey)
			if !ok {
				lc.Error("meta data client not set")
			}
			return metadata.NewDeviceClient(local.New(metadataURL + clients.ApiDeviceRoute))
		},
	})

	thingsboardGateway := NewThingsboardGateway(dic)
	if err := thingsboardGateway.Serve(); err != nil {
		lc.Error(fmt.Sprintf("servve thingsboard gateway: %s", err.Error()))
		return false
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

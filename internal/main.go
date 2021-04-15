package internal

import (
	"context"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/httpserver"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/message"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/testing"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/gorilla/mux"
	edgex_thingsboard "github.com/inspii/edgex-thingsboard"
	bootstrap2 "github.com/inspii/edgex-thingsboard/internal/bootstrap"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/handler/pubsub"
	"os"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *mux.Router, readyStream chan<- bool) {
	startupTimer := startup.NewStartUpTimer(ControlAgentServiceKey)

	f := flags.New()
	f.Parse(os.Args[1:])
	configuration := &bootstrap2.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	httpServer := httpserver.NewBootstrap(router, true)

	bootstrap.Run(
		ctx,
		cancel,
		f,
		ControlAgentServiceKey,
		ConfigStemCore+ConfigMajorVersion,
		configuration,
		startupTimer,
		dic,
		[]interfaces.BootstrapHandler{
			pubsub.NewPubSub(httpServer, configuration).BootstrapHandler,
			NewBootstrap(router).BootstrapHandler,
			httpServer.BootstrapHandler,
			message.NewBootstrap(ControlAgentServiceKey, edgex_thingsboard.Version).BootstrapHandler,
			testing.NewBootstrap(httpServer, readyStream).BootstrapHandler,
		})
}

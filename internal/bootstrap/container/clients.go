package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap"
)

var ClientsName = di.TypeInstanceToName((*bootstrap.ServiceRoutes)(nil))

func ClientsFrom(get di.Get) *bootstrap.ServiceRoutes {
	return get(ClientsName).(*bootstrap.ServiceRoutes)
}

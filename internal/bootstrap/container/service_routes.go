package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap"
)

var ServiceRoutesName = di.TypeInstanceToName((*bootstrap.ServiceRoutes)(nil))

func ServiceRoutesFrom(get di.Get) *bootstrap.ServiceRoutes {
	return get(ServiceRoutesName).(*bootstrap.ServiceRoutes)
}

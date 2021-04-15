package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap"
)

var ConfigurationName = di.TypeInstanceToName(bootstrap.ConfigurationStruct{})

func ConfigurationFrom(get di.Get) *bootstrap.ConfigurationStruct {
	return get(ConfigurationName).(*bootstrap.ConfigurationStruct)
}

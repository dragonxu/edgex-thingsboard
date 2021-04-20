package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
)

var MetaDataDeviceClientName = di.TypeInstanceToName((*metadata.DeviceClient)(nil))

func MetaDataDeviceClientFrom(get di.Get) metadata.DeviceClient {
	return get(MetaDataDeviceClientName).(metadata.DeviceClient)
}

package device

import (
	"context"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"time"
)

func ListDevices(dic *di.Container) ([]models.Device, error) {
	client := container.MetaDataDeviceClientFrom(dic.Get)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Devices(ctx)
}

package dependencies

import (
	"github.com/warehouse/user-service/internal/handler/grpc"
)

func (d *dependencies) UserHandler() *grpc.UserHandler {
	if d.userHandler == nil {
		d.userHandler = grpc.NewUserHandler(
			d.cfg.Timeouts,
			d.log,
			d.UserService(),
		)
	}

	return d.userHandler
}

package dependencies

import (
	"github.com/warehouse/user-service/internal/service/user"
)

func (d *dependencies) UserService() user.Service {
	if d.userService == nil {
		d.userService = user.NewService(
			*d.cfg,
			d.log,
			d.PgxTransactionRepo(),
			d.UserRepo(),
		)
	}

	return d.userService
}

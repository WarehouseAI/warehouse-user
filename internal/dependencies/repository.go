package dependencies

import (
	"github.com/warehouse/user-service/internal/repository/operations/transactions"
	"github.com/warehouse/user-service/internal/repository/operations/user"
)

func (d *dependencies) PgxTransactionRepo() transactions.Repository {
	if d.pgxTransactionRepo == nil {
		d.pgxTransactionRepo = transactions.NewPgxRepository(d.PostgresClient())
	}
	return d.pgxTransactionRepo
}

func (d *dependencies) UserRepo() user.Repository {
	if d.userRepo == nil {
		d.userRepo = user.NewPGRepository(d.log, d.PostgresClient())
	}
	return d.userRepo
}

package transactions

import (
	"context"

	"github.com/warehouse/user-service/internal/db"

	"github.com/jmoiron/sqlx"
)

type (
	Tx struct {
		*sqlx.Tx
		ctx context.Context
	}

	repositoryPG struct {
		client *db.PostgresClient
	}
)

func NewPgxRepository(client *db.PostgresClient) Repository {
	return &repositoryPG{
		client: client,
	}
}

func (repo *repositoryPG) StartTransaction(ctx context.Context) (Transaction, error) {
	tx, err := repo.client.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Tx:  tx,
		ctx: ctx,
	}, nil
}

// Txm - get transaction method
func (t *Tx) Txm() *sqlx.Tx {
	return t.Tx
}

func (t *Tx) Rollback() {
	_ = t.Tx.Rollback()
}

package user

import (
	"context"

	"github.com/warehouse/user-service/internal/db"
	"github.com/warehouse/user-service/internal/pkg/errors/repository_errors"
	"github.com/warehouse/user-service/internal/pkg/logger"
	"github.com/warehouse/user-service/internal/repository/models"
	"github.com/warehouse/user-service/internal/repository/operations/transactions"
)

type repositoryPG struct {
	log logger.Logger
	pg  *db.PostgresClient
}

func NewPGRepository(log logger.Logger, client *db.PostgresClient) Repository {
	return &repositoryPG{
		pg:  client,
		log: log.Named("pg_users"),
	}
}

func (r *repositoryPG) Create(
	ctx context.Context,
	tx transactions.Transaction,
	u models.User,
) (models.User, error) {
	query := `
    INSERT INTO users (firstname, lastname, username, email, hash, role)
    VALUES(:firstname, :lastname, :username, :email, :hash, :role)
  `

	res, err := tx.Txm().NamedExecContext(ctx, query, u)
	if err != nil {
		return models.User{}, r.log.ErrorRepo(err, repository_errors.PostgresqlExecRaw, query)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.User{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	if rowsAffected != 1 {
		return models.User{}, r.log.ErrorRepo(err, repository_errors.PostgresqlRowsAffectedRaw, query)
	}

	return u, nil
}

func (r *repositoryPG) GetByUsername(
	ctx context.Context,
	tx transactions.Transaction,
	username string,
) (models.User, error) {
	cond := `WHERE u.username = $1`
	list, err := r.getUserByCondition(ctx, tx.Txm(), cond, username)
	if err != nil {
		return models.User{}, err
	}

	if len(list) == 0 {
		return models.User{}, repository_errors.PostgresqlNotFound
	}

	return list[0], nil
}

func (r *repositoryPG) GetById(
	ctx context.Context,
	tx transactions.Transaction,
	id string,
) (models.User, error) {
	cond := `WHERE u.id = $1`
	list, err := r.getUserByCondition(ctx, tx.Txm(), cond, id)
	if err != nil {
		return models.User{}, err
	}

	if len(list) == 0 {
		return models.User{}, repository_errors.PostgresqlNotFound
	}

	return list[0], nil
}

func (r *repositoryPG) GetByEmail(
	ctx context.Context,
	tx transactions.Transaction,
	email string,
) (models.User, error) {
	cond := "WHERE u.email = $1"
	list, err := r.getUserByCondition(ctx, tx.Txm(), cond, email)
	if err != nil {
		return models.User{}, err
	}

	if len(list) == 0 {
		return models.User{}, repository_errors.PostgresqlNotFound
	}

	return list[0], nil
}

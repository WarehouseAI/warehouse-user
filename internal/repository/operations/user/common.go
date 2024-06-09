package user

import (
	"context"
	"fmt"

	"github.com/warehouse/user-service/internal/pkg/errors/repository_errors"
	"github.com/warehouse/user-service/internal/repository/models"

	"github.com/jmoiron/sqlx"
)

func (r *repositoryPG) getUserByCondition(
	ctx context.Context,
	executor sqlx.ExtContext,
	condition string,
	params ...interface{},
) ([]models.User, error) {
	query := `
    SELECT u.id, u.firstname, u.lastname, u.username, u.hash, u.role, u.created_at, u.updated_at
    FROM users as u
  `
	query = fmt.Sprintf("%s %s", query, condition)

	var list []models.User
	err := sqlx.SelectContext(ctx, executor, &list, query, params...)
	if err != nil {
		return []models.User{}, r.log.ErrorRepo(err, repository_errors.PostgresqlGetRaw, query)
	}

	return list, nil
}

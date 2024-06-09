package user

import (
	"context"

	"github.com/warehouse/user-service/internal/repository/models"
	"github.com/warehouse/user-service/internal/repository/operations/transactions"
)

type Repository interface {
	Create(ctx context.Context, tx transactions.Transaction, user models.User) (models.User, error)
	GetByUsername(ctx context.Context, tx transactions.Transaction, username string) (models.User, error)
	GetByEmail(ctx context.Context, tx transactions.Transaction, email string) (models.User, error)
	GetById(ctx context.Context, tx transactions.Transaction, id string) (models.User, error)
}

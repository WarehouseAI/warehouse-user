package user

import (
	"context"

	"github.com/warehouse/user-service/internal/config"
	"github.com/warehouse/user-service/internal/domain"
	"github.com/warehouse/user-service/internal/handler/models"
	"github.com/warehouse/user-service/internal/pkg/errors"
	"github.com/warehouse/user-service/internal/pkg/logger"
	"github.com/warehouse/user-service/internal/repository/operations/transactions"
	userRepo "github.com/warehouse/user-service/internal/repository/operations/user"
)

type (
	Service interface {
		Create(ctx context.Context, reqData models.CreateUserRequest) (domain.User, *errors.Error)
		GetByEmail(ctx context.Context, email string) (domain.User, *errors.Error)
		GetById(ctx context.Context, id string) (domain.User, *errors.Error)
		GetByUsername(ctx context.Context, username string) (domain.User, *errors.Error)
	}

	service struct {
		cfg config.Config

		txRepo   transactions.Repository
		userRepo userRepo.Repository

		log logger.Logger
	}
)

func NewService(
	cfg config.Config,
	log logger.Logger,
	txRepo transactions.Repository,
	userRepo userRepo.Repository,
) Service {
	return &service{
		cfg:      cfg,
		log:      log,
		txRepo:   txRepo,
		userRepo: userRepo,
	}
}

func (s *service) Create(ctx context.Context, reqData models.CreateUserRequest) (domain.User, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	u := domain.User{
		Firstname: reqData.Firstname,
		Lastname:  reqData.Lastname,
		Username:  reqData.Username,
		Email:     reqData.Email,
		Hash:      reqData.Hash,
		Role:      domain.RoleUser,
	}

	createdUser, err := s.userRepo.Create(ctx, tx, u.ToModel())
	if err != nil {
		return domain.User{}, errors.DatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}

	return domain.User{}.FromModel(createdUser), nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (domain.User, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	user, err := s.userRepo.GetByEmail(ctx, tx, email)
	if err != nil {
		return domain.User{}, errors.DatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}

	return domain.User{}.FromModel(user), nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (domain.User, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	user, err := s.userRepo.GetByUsername(ctx, tx, username)
	if err != nil {
		return domain.User{}, errors.DatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}

	return domain.User{}.FromModel(user), nil
}

func (s *service) GetById(ctx context.Context, id string) (domain.User, *errors.Error) {
	tx, err := s.txRepo.StartTransaction(ctx)
	if err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}
	defer tx.Rollback()

	user, err := s.userRepo.GetById(ctx, tx, id)
	if err != nil {
		return domain.User{}, errors.DatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, s.log.ServiceTxError(err)
	}

	return domain.User{}.FromModel(user), nil
}

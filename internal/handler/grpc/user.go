package grpc

import (
	"context"

	"github.com/warehouse/user-service/internal/config"
	handler_converters "github.com/warehouse/user-service/internal/handler/converters"
	"github.com/warehouse/user-service/internal/handler/models"
	"github.com/warehouse/user-service/internal/pkg/logger"
	"github.com/warehouse/user-service/internal/service/user"
	"github.com/warehouse/user-service/internal/warehousepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	UserHandler struct {
		warehousepb.UnimplementedUserServiceServer
		timeouts config.Timeouts
		log      logger.Logger
		userSvc  user.Service
	}
)

func NewUserHandler(
	timeouts config.Timeouts,
	log logger.Logger,
	userSvc user.Service,
) *UserHandler {
	return &UserHandler{
		timeouts: timeouts,
		log:      log,
		userSvc:  userSvc,
	}
}

func (s *UserHandler) CreateUser(ctx context.Context, req *warehousepb.CreateUserRequest) (*warehousepb.User, error) {
	if req == nil || req.Email == "" || req.Hash == "" || req.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty required request data")
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.timeouts.RequestTimeout)
	defer cancel()

	createdUser, e := s.userSvc.Create(ctx, models.CreateUserRequest{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Username,
		Hash:      req.Hash,
		Email:     req.Email,
	})
	if e != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(e)
	}

	return createdUser.ToProto(), nil
}

func (s *UserHandler) GetUserByEmail(ctx context.Context, req *warehousepb.GetUserByEmailRequest) (*warehousepb.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty required request data")
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.timeouts.RequestTimeout)
	defer cancel()

	user, e := s.userSvc.GetByEmail(ctx, req.Email)
	if e != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(e)
	}

	return user.ToProto(), nil
}

func (s *UserHandler) GetUserByLogin(ctx context.Context, req *warehousepb.GetUserByLoginRequest) (*warehousepb.User, error) {
	if req == nil || req.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty required request data")
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.timeouts.RequestTimeout)
	defer cancel()

	user, e := s.userSvc.GetByUsername(ctx, req.Username)
	if e != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(e)
	}

	return user.ToProto(), nil
}

func (s *UserHandler) GetUserById(ctx context.Context, req *warehousepb.GetUserByIdRequest) (*warehousepb.User, error) {
	if req == nil || req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty required request data")
	}

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.timeouts.RequestTimeout)
	defer cancel()

	user, e := s.userSvc.GetById(ctx, req.Id)
	if e != nil {
		return nil, handler_converters.MakeStatusFromErrorsError(e)
	}

	return user.ToProto(), nil
}

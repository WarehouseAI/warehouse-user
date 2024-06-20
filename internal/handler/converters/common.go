package converters

import (
	"github.com/warehouse/user-service/internal/handler/models"
	"github.com/warehouse/user-service/internal/pkg/errors"
	"github.com/warehouse/user-service/internal/pkg/errors/repository_errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeJsonErrorResponseWithErrorsError(err *errors.Error) models.ErrorResponse {
	res := models.ErrorResponse{
		Code:   err.Code,
		Reason: err.Reason,
	}

	if err.Details != nil {
		res.Details = err.Details.Error()
	}

	return res
}

func MakeStatusFromErrorsError(err *errors.Error) error {
	details := err.Reason
	code := codes.Internal

	if err.Details != nil {
		details = err.Details.Error()
		if err.Details == repository_errors.PostgresqlNotFound {
			code = codes.NotFound
		}
	}

	return status.Errorf(code, details)
}

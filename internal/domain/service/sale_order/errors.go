package sale_order

import (
	"github.com/kiaplayer/clean-architecture-example/internal/domain/errors"
)

type ErrValidation struct{ errors.AppError }

func NewErrValidation(reason string, cause error) *ErrValidation {
	return &ErrValidation{errors.NewAppError(reason, cause)}
}

//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package get_sale_order

import (
	"context"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
)

type saleOrderService interface {
	GetOrderByID(ctx context.Context, id uint64) (*document.SaleOrder, error)
}

type UseCase struct {
	saleOrderService saleOrderService
}

func NewUseCase(sos saleOrderService) *UseCase {
	return &UseCase{
		saleOrderService: sos,
	}
}

func (u *UseCase) Handle(ctx context.Context, id uint64) (saleOrder *document.SaleOrder, err error) {
	return u.saleOrderService.GetOrderByID(ctx, id)
}

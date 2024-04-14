//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package create_sale_order

import (
	"context"
	"time"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

type timeGenerator interface {
	NowDate() time.Time
}

type numberGenerator interface {
	GenerateNumber(date time.Time, company reference.Company) string
}

type saleOrderService interface {
	CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error)
}

type UseCase struct {
	timeGenerator    timeGenerator
	numberGenerator  numberGenerator
	saleOrderService saleOrderService
}

func NewUseCase(tg timeGenerator, ng numberGenerator, sos saleOrderService) *UseCase {
	return &UseCase{
		timeGenerator:    tg,
		numberGenerator:  ng,
		saleOrderService: sos,
	}
}

func (u *UseCase) Handle(
	ctx context.Context,
	saleOrder *document.SaleOrder,
) (saleOrderUpdated *document.SaleOrder, err error) {
	saleOrder.Date = u.timeGenerator.NowDate()
	saleOrder.Number = u.numberGenerator.GenerateNumber(saleOrder.Date, saleOrder.Company)
	saleOrderUpdated, err = u.saleOrderService.CreateOrder(ctx, saleOrder)
	if err != nil {
		return nil, err
	}
	return saleOrderUpdated, nil
}

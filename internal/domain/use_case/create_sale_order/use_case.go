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
	products []document.SaleOrderProduct,
	company reference.Company,
	customer reference.Customer,
	appendUser *reference.User,
) (saleOrder *document.SaleOrder, err error) {
	docDate := u.timeGenerator.NowDate()
	docNumber := u.numberGenerator.GenerateNumber(docDate, company)
	saleOrder = &document.SaleOrder{
		Document: document.Document{
			Number:        docNumber,
			Date:          docDate,
			Status:        document.StatusDraft,
			BasisDocument: nil,
			Company:       company,
			AppendUser:    appendUser,
			ChangeUser:    appendUser,
		},
		Customer: customer,
		Products: products,
	}
	saleOrder, err = u.saleOrderService.CreateOrder(ctx, saleOrder)
	if err != nil {
		return nil, err
	}
	return saleOrder, nil
}

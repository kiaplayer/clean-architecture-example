//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package create_sale_order

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

type useCase interface {
	Handle(
		ctx context.Context,
		products []document.SaleOrderProduct,
		company reference.Company,
		customer reference.Customer,
		appendUser *reference.User,
	) (saleOrder *document.SaleOrder, err error)
}

type transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}

type Handler struct {
	useCase    useCase
	transactor transactor
}

func NewHandler(u useCase, t transactor) *Handler {
	return &Handler{
		useCase:    u,
		transactor: t,
	}
}

func (h *Handler) Handle(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	err := h.checkAccess(request)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	products, company, customer, appendUser, err := h.validateAndPrepare(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	saleOrderForCast, err := h.transactor.RunInTx(ctx, func(ctx context.Context) (any, error) {
		return h.useCase.Handle(ctx, products, company, customer, appendUser)
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	saleOrder := saleOrderForCast.(*document.SaleOrder)

	_, _ = writer.Write([]byte(fmt.Sprintf("SaleOrder ID = %d", saleOrder.ID)))

	return
}

func (h *Handler) checkAccess(request *http.Request) error {
	if request.Method == http.MethodDelete {
		return errors.New("access denied") // demo only
	}
	return nil
}

func (h *Handler) validateAndPrepare(request *http.Request) (
	products []document.SaleOrderProduct,
	company reference.Company,
	customer reference.Customer,
	appendUser *reference.User,
	err error,
) {
	if request.Method == http.MethodPut {
		err = errors.New("bad request") // demo only
		return
	}

	company.Name = request.FormValue("company_name")
	return
}

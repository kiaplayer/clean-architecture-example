//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package create_sale_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/service/sale_order"
	"github.com/kiaplayer/clean-architecture-example/internal/handlers/create_sale_order/dto"
)

type useCase interface {
	Handle(context.Context, *document.SaleOrder) (*document.SaleOrder, error)
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

func (h *Handler) validateAndPrepare(request *http.Request) (*document.SaleOrder, error) {
	var saleOrderDTO dto.SaleOrder
	err := json.NewDecoder(request.Body).Decode(&saleOrderDTO)
	if err != nil {
		return nil, err
	}

	isValid := saleOrderDTO.CustomerID > 0 && len(saleOrderDTO.Products) > 0
	if isValid {
		for _, product := range saleOrderDTO.Products {
			if (product.ProductID == 0) || (product.Quantity == 0) {
				isValid = false
				break
			}
		}
	}

	if !isValid {
		return nil, errors.New("bad order data")
	}

	return dto.SaleOrderDtoToSaleOrder(saleOrderDTO), nil
}

func (h *Handler) Handle(writer http.ResponseWriter, request *http.Request) {
	err := h.checkAccess(request)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	saleOrder, err := h.validateAndPrepare(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	saleOrderUpdated, err := h.transactor.RunInTx(request.Context(), func(ctx context.Context) (any, error) {
		return h.useCase.Handle(ctx, saleOrder)
	})
	if err != nil {
		var errTarget *sale_order.ErrValidation
		if errors.As(err, &errTarget) {
			http.Error(writer, errTarget.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, "internal error", http.StatusInternalServerError)
		}
		return
	}

	saleOrder = saleOrderUpdated.(*document.SaleOrder)

	_, _ = writer.Write([]byte(fmt.Sprintf("SaleOrder ID = %d", saleOrder.ID)))

	return
}

func (h *Handler) checkAccess(request *http.Request) error {
	if request.Method == http.MethodDelete {
		return errors.New("access denied") // demo only
	}
	return nil
}

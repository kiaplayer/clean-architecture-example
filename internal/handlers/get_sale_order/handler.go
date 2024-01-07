//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package get_sale_order

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
)

type useCase interface {
	Handle(
		ctx context.Context,
		id uint64,
	) (saleOrder *document.SaleOrder, err error)
}

type Handler struct {
	useCase useCase
}

func NewHandler(u useCase) *Handler {
	return &Handler{
		useCase: u,
	}
}

func (h *Handler) Handle(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	err := h.checkAccess(request)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	saleOrderID, err := h.validateAndPrepare(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	saleOrder, err := h.useCase.Handle(ctx, saleOrderID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = writer.Write([]byte(fmt.Sprintf("SaleOrder ID = %d", saleOrder.ID)))

	return
}

func (h *Handler) checkAccess(request *http.Request) error {
	if request.Method == http.MethodDelete {
		return errors.New("access denied") // demo only
	}
	return nil
}

func (h *Handler) validateAndPrepare(request *http.Request) (uint64, error) {
	id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("bad id: %w", err)
	}
	if id <= 0 {
		return 0, errors.New("bad id")
	}

	return uint64(id), nil
}

package sale_order

import (
	"context"
	"fmt"
	"slices"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
)

type repository interface {
	CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error)
}

type Service struct {
	repository repository
}

func NewService(r repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error) {
	order, err := s.ValidateOrder(order)
	if err != nil {
		return order, err
	}
	return s.repository.CreateOrder(ctx, order)
	// TODO: Additinal logic (reserve, email...)
}

func (s *Service) ValidateOrder(order *document.SaleOrder) (*document.SaleOrder, error) {
	if !slices.Contains(document.ValidStatuses, order.Status) {
		return order, fmt.Errorf("bad status: %d", order.Status)
	}
	return order, nil
}

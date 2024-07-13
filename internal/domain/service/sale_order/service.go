//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mocks/$GOFILE
package sale_order

import (
	"context"
	"fmt"
	"slices"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
)

type repository interface {
	CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error)
	GetByID(ctx context.Context, id uint64) (*document.SaleOrder, error)
}

type productRepository interface {
	Exists(ctx context.Context, id uint64) (bool, error)
}

type Service struct {
	repository        repository
	productRepository productRepository
}

func NewService(r repository, pr productRepository) *Service {
	return &Service{
		repository:        r,
		productRepository: pr,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error) {
	order, err := s.ValidateOrder(ctx, order)
	if err != nil {
		return order, err
	}
	savedSaleOrder, err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	// TODO: Additinal logic goes here
	// - Reserve products
	// - Send emails to customer and manager (via pgq)
	// - etc
	return savedSaleOrder, nil
}

func (s *Service) ValidateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error) {
	if !slices.Contains(document.ValidStatuses, order.Status) {
		return nil, fmt.Errorf("bad status: %d", order.Status)
	}

	if len(order.Products) > 0 {
		productsIDs := make([]uint64, 0, len(order.Products))

		for _, product := range order.Products {
			productsIDs = append(productsIDs, product.Product.ID)
		}

		for _, productID := range slices.Compact(productsIDs) {
			exists, err := s.productRepository.Exists(ctx, productID)
			if err != nil {
				return nil, err
			}

			if !exists {
				return nil, fmt.Errorf("bad product id: %d", productID)
			}
		}
	}

	return order, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id uint64) (*document.SaleOrder, error) {
	return s.repository.GetByID(ctx, id)
}

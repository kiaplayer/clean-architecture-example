package main

import (
	"database/sql"

	"github.com/kiaplayer/clean-architecture-example/internal/adapters/repositories/sale_order"
	saleorderservice "github.com/kiaplayer/clean-architecture-example/internal/domain/service/sale_order"
	createsaleorderusecase "github.com/kiaplayer/clean-architecture-example/internal/domain/use_case/create_sale_order"
	"github.com/kiaplayer/clean-architecture-example/internal/handlers/create_sale_order"
	"github.com/kiaplayer/clean-architecture-example/pkg/generators"
	"github.com/kiaplayer/clean-architecture-example/pkg/storage/db"
)

func main() {
	conn := &sql.DB{}

	transactor := db.NewTransactor(conn)
	timeGenerator := generators.NewTimeGenerator()
	numberGenerator := generators.NewNumberGenerator()
	saleOrderRepo := sale_order.NewRepository(conn)
	saleOrderService := saleorderservice.NewService(saleOrderRepo)
	useCase := createsaleorderusecase.NewUseCase(timeGenerator, numberGenerator, saleOrderService)

	create_sale_order.NewHandler(useCase, transactor)
}

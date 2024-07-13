package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kiaplayer/clean-architecture-example/internal/adapters/repositories/document/sale_order"
	"github.com/kiaplayer/clean-architecture-example/internal/adapters/repositories/reference/product"
	saleorderservice "github.com/kiaplayer/clean-architecture-example/internal/domain/service/sale_order"
	createsaleorderusecase "github.com/kiaplayer/clean-architecture-example/internal/domain/use_case/create_sale_order"
	getsaleorderusecase "github.com/kiaplayer/clean-architecture-example/internal/domain/use_case/get_sale_order"
	"github.com/kiaplayer/clean-architecture-example/internal/handlers/create_sale_order"
	"github.com/kiaplayer/clean-architecture-example/internal/handlers/get_sale_order"
	"github.com/kiaplayer/clean-architecture-example/pkg/generators"
	"github.com/kiaplayer/clean-architecture-example/pkg/storage/db"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dbConn, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(dbConn *sql.DB) {
		closeErr := dbConn.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}(dbConn)

	transactor := db.NewTransactor(dbConn)
	timeGenerator := generators.NewTimeGenerator()
	numberGenerator := generators.NewNumberGenerator()
	saleOrderRepo := sale_order.NewRepository(dbConn)
	productRepo := product.NewRepository(dbConn)
	saleOrderService := saleorderservice.NewService(saleOrderRepo, productRepo)

	createSaleOrderHandler := create_sale_order.NewHandler(
		createsaleorderusecase.NewUseCase(timeGenerator, numberGenerator, saleOrderService),
		transactor,
	)
	getSaleOrderHandler := get_sale_order.NewHandler(getsaleorderusecase.NewUseCase(saleOrderService))

	srvMux := http.NewServeMux()
	srvMux.HandleFunc("POST /sale-order", createSaleOrderHandler.Handle)
	srvMux.HandleFunc("GET /sale-order", getSaleOrderHandler.Handle)

	srv := http.Server{
		Addr:    os.Getenv("SERVICE_ADDR"),
		Handler: srvMux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Service shutting down...")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Service started at: %s", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

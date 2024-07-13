//go:build integration

package sale_order

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
)

const testDBFilePath = "sqlite_test.db"

type TestRepositorySuite struct {
	suite.Suite
	db *sql.DB
}

func TestRepositoryByTestSuite(t *testing.T) {
	suite.Run(t, new(TestRepositorySuite))
}

func (rts *TestRepositorySuite) SetupSuite() {
	_ = os.Remove(testDBFilePath)

	dbConn, err := sql.Open("sqlite3", testDBFilePath)
	if err != nil {
		rts.Failf("cannot open db connection before tests: %s", err.Error())
	}

	driver, err := sqlite3.WithInstance(dbConn, &sqlite3.Config{})
	if err != nil {
		rts.Failf("cannot init db driver before tests: %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../db/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		rts.Failf("cannot init db migrator before tests: %s", err.Error())
	}

	err = m.Up()
	if err != nil {
		rts.Failf("cannot apply db migrations before tests: %s", err.Error())
	}

	rts.db = dbConn
}

func (rts *TestRepositorySuite) TearDownSuite() {
	err := rts.db.Close()
	if err != nil {
		rts.Failf("tear down suite: %s", err.Error())
	}
	_ = os.Remove(testDBFilePath)
}

func (rts *TestRepositorySuite) TestCreateOrder_Success() {
	// arrange
	ctx := context.Background()

	tx, _ := rts.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	repository := NewRepository(rts.db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "1",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	// act
	actual, err := repository.CreateOrder(ctx, saleOrder)

	// assert
	rts.NoError(err)
	rts.NotNil(actual)
	rts.NotEqual(0, actual.ID)
}

func (rts *TestRepositorySuite) TestCreateOrder_InsertError() {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())

	tx, _ := rts.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	repository := NewRepository(rts.db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "2",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	// act
	cancel()
	actual, err := repository.CreateOrder(ctx, saleOrder)

	// assert
	rts.ErrorContains(err, "context canceled")
	rts.Nil(actual)
}

func (rts *TestRepositorySuite) TestGetByID_Success() {
	// arrange
	ctx := context.Background()

	tx, _ := rts.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	repository := NewRepository(tx)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "3",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	// act & assert
	saleOrder, err := repository.CreateOrder(ctx, saleOrder)
	rts.NoError(err)
	rts.NotNil(saleOrder)
	rts.NotEqual(0, saleOrder.ID)

	saleOrder, err = repository.GetByID(ctx, saleOrder.ID)
	rts.NoError(err)
	rts.NotNil(saleOrder)
}

func (rts *TestRepositorySuite) TestGetByID_NotFound() {
	// arrange
	ctx := context.Background()

	repository := NewRepository(rts.db)

	// act
	actual, err := repository.GetByID(ctx, 999)

	// assert
	rts.NoError(err)
	rts.Nil(actual)
}

func (rts *TestRepositorySuite) TestGetByID_BadStatusError() {
	// arrange
	ctx := context.Background()

	tx, _ := rts.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	repository := NewRepository(tx)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "4",
			Date:   time.Now(),
			Status: -999,
		},
	}

	// act & assert
	saleOrder, err := repository.CreateOrder(ctx, saleOrder)
	rts.NoError(err)
	rts.NotNil(saleOrder)
	rts.NotEqual(0, saleOrder.ID)

	saleOrder, err = repository.GetByID(ctx, saleOrder.ID)
	rts.ErrorContains(err, "bad status")
	rts.Nil(saleOrder)
}

func (rts *TestRepositorySuite) TestGetByID_QueryError() {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())

	repository := NewRepository(rts.db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "5",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	// act
	cancel()
	actual, err := repository.GetByID(ctx, saleOrder.ID)

	// assert
	rts.ErrorContains(err, "context canceled")
	rts.Nil(actual)
}

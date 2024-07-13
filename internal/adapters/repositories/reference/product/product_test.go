package product

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExists_Success(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	id := uint64(1)

	queryResult := sqlmock.NewRows([]string{"cnt"}).AddRow(id)

	mock.
		ExpectQuery("SELECT id FROM product WHERE id = ?").
		WithArgs(id).
		WillReturnRows(queryResult)

	// act
	actual, err := repository.Exists(ctx, id)

	// assert
	assert.NoError(t, err)
	assert.True(t, actual)
}

func TestExists_NotFound(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	id := uint64(1)

	queryResult := sqlmock.NewRows([]string{"cnt"})

	mock.
		ExpectQuery("SELECT id FROM product WHERE id = ?").
		WithArgs(id).
		WillReturnRows(queryResult)

	// act
	actual, err := repository.Exists(ctx, id)

	// assert
	assert.NoError(t, err)
	assert.False(t, actual)
}

func TestExists_Error(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	id := uint64(1)

	queryError := errors.New("some query error")

	mock.
		ExpectQuery("SELECT id FROM product WHERE id = ?").
		WithArgs(id).
		WillReturnError(queryError)

	// act
	actual, err := repository.Exists(ctx, id)

	// assert
	assert.False(t, actual)
	assert.ErrorContains(t, err, queryError.Error())
}

package repository

import (
	"avito-testTask/models"
	// "database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestProductPostgres_ProductCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewProductPostgres(db)

	testTime := time.Now()
	productID := uuid.New()
	receptionID := uuid.New()

	input := models.Product{
		Id:          productID,
		DateTime:    testTime,
		Type:        "электроника",
		ReceptionId: receptionID,
	}

	type mockBehavior func(m sqlmock.Sqlmock, product models.Product)

	testTable := []struct {
		name         string
		input        models.Product
		mockBehavior mockBehavior
		expected     models.Product
		wantErr      bool
	}{
		{
			name:  "OK",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, product models.Product) {
				m.ExpectQuery(`INSERT INTO products`).
					WithArgs(product.Id, product.DateTime, product.Type, product.ReceptionId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
						AddRow(product.Id, product.DateTime, product.Type, product.ReceptionId))
			},
			expected: input,
			wantErr:  false,
		},
		{
			name:  "DB Error",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, product models.Product) {
				m.ExpectQuery(`INSERT INTO products`).
					WithArgs(product.Id, product.DateTime, product.Type, product.ReceptionId).
					WillReturnError(fmt.Errorf("insert error"))
			},
			expected: models.Product{},
			wantErr:  true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, tt.input)

			got, err := r.ProductCreate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ProductCreate() got = %v, want %v", got, tt.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

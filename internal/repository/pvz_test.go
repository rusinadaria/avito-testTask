package repository

import (
	"avito-testTask/models"
	"fmt"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestPVZCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewPVZPostgres(db)

	testDate := time.Now()
	testUUID, _ := uuid.Parse("1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d")
	input := models.PVZ{
		Id:               testUUID,
		RegistrationDate: testDate,
		City:             "Москва",
	}

	type mockBehavior func(m sqlmock.Sqlmock, pvz models.PVZ)

	testCases := []struct {
		name         string
		input        models.PVZ
		mockBehavior mockBehavior
		expected     models.PVZ
		wantErr      bool
	}{
		{
			name:  "OK",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, pvz models.PVZ) {
				m.ExpectQuery(`INSERT INTO pvz`).
					WithArgs(pvz.Id, pvz.RegistrationDate, pvz.City).
					WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
						AddRow(pvz.Id, pvz.RegistrationDate, pvz.City))
			},
			expected: input,
			wantErr:  false,
		},
		{
			name:  "Insert error",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, pvz models.PVZ) {
				m.ExpectQuery(`INSERT INTO pvz`).
					WithArgs(pvz.Id, pvz.RegistrationDate, pvz.City).
					WillReturnError(fmt.Errorf("insert failed"))
			},
			expected: models.PVZ{},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.input)

			result, err := r.PVZCreate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("PVZCreate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if result != tc.expected {
				t.Errorf("PVZCreate() result = %+v, expected %+v", result, tc.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

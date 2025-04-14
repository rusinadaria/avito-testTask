package repository

import (
	"avito-testTask/models"
	// "database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestReceptionCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewReceptionPostgres(db)

	testTime := time.Now()
	receptionID := uuid.New()
	pvzID := uuid.New()
	status := models.InProgress

	input := models.Reception{
		Id:       receptionID,
		DateTime: testTime,
		PvzId:    pvzID,
		Status:   status,
	}

	type mockBehavior func(m sqlmock.Sqlmock, reception models.Reception)

	testCases := []struct {
		name         string
		input        models.Reception
		mockBehavior mockBehavior
		expected     models.Reception
		wantErr      bool
	}{
		{
			name:  "OK",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, reception models.Reception) {
				m.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
						AddRow(reception.Id, reception.DateTime, reception.PvzId, reception.Status))
			},
			expected: input,
			wantErr:  false,
		},
		{
			name:  "Insert error",
			input: input,
			mockBehavior: func(m sqlmock.Sqlmock, reception models.Reception) {
				m.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnError(fmt.Errorf("insert failed"))
			},
			expected: models.Reception{},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.input)

			result, err := r.ReceptionCreate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReceptionCreate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if result != tc.expected {
				t.Errorf("ReceptionCreate() result = %+v, expected %+v", result, tc.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCloseReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewReceptionPostgres(db)

	testTime := time.Now()
	receptionID := uuid.New()
	pvzID := uuid.New()
	closedStatus := models.Close

	type mockBehavior func(m sqlmock.Sqlmock, pvzId uuid.UUID)

	testCases := []struct {
		name         string
		pvzId        uuid.UUID
		mockBehavior mockBehavior
		expected     models.Reception
		wantErr      bool
	}{
		{
			name:  "OK",
			pvzId: pvzID,
			mockBehavior: func(m sqlmock.Sqlmock, pvzId uuid.UUID) {
				m.ExpectQuery(`UPDATE receptions\s+SET status = 'close'`).
					WithArgs(pvzId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
						AddRow(receptionID, testTime, pvzId, closedStatus))
			},
			expected: models.Reception{
				Id:       receptionID,
				DateTime: testTime,
				PvzId:    pvzID,
				Status:   closedStatus,
			},
			wantErr: false,
		},
		{
			name:  "Query Error",
			pvzId: pvzID,
			mockBehavior: func(m sqlmock.Sqlmock, pvzId uuid.UUID) {
				m.ExpectQuery(`UPDATE receptions\s+SET status = 'close'`).
					WithArgs(pvzId).
					WillReturnError(fmt.Errorf("update failed"))
			},
			expected: models.Reception{},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.pvzId)

			result, err := r.CloseReception(tc.pvzId)
			if (err != nil) != tc.wantErr {
				t.Errorf("CloseReception() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if result != tc.expected {
				t.Errorf("CloseReception() result = %+v, expected %+v", result, tc.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

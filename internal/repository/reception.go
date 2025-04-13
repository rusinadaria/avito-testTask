package repository

import (
	// "github.com/jmoiron/sqlx"
	"avito-testTask/models"
	"database/sql"
	// "log"
	"fmt"
	"github.com/google/uuid"
)

type ReceptionPostgres struct {
	db *sql.DB
}

func NewReceptionPostgres(db *sql.DB) *ReceptionPostgres {
	return &ReceptionPostgres{db: db}
}

func (r *ReceptionPostgres) GetLastReceptionByPVZ(pvzId uuid.UUID) (models.Reception, error) { // изменить
	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1
		ORDER BY date_time DESC
		LIMIT 1
	`

	var reception models.Reception
	err := r.db.QueryRow(query, pvzId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)

	if err == sql.ErrNoRows {
		return models.Reception{}, nil 
	}
	if err != nil {
		return models.Reception{}, err
	}

	return reception, nil
}

func (r *ReceptionPostgres) GetLastReceptionStatus(pvzId uuid.UUID) (models.Status, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1
		ORDER BY date_time DESC
		LIMIT 1
	`

	var reception models.Reception
	err := r.db.QueryRow(query, pvzId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)
	if err != nil {
		return "", err
	}

	if reception.Status == models.Close {
		return models.Close, nil
	}
	return models.InProgress, nil
}

func (r *ReceptionPostgres) CloseReception(pvzId uuid.UUID) (models.Reception, error) {
	query := `
		UPDATE receptions
		SET status = 'close'
		WHERE id = (
			SELECT id FROM receptions
			WHERE pvz_id = $1
			ORDER BY date_time DESC
			LIMIT 1
		)
		RETURNING id, date_time, pvz_id, status;
	`


	var updatedReception models.Reception
	err := r.db.QueryRow(query, pvzId).Scan(
		&updatedReception.Id,
		&updatedReception.DateTime,
		&updatedReception.PvzId,
		&updatedReception.Status,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return models.Reception{}, fmt.Errorf("failed to update status Reception: %w", err)
	}

	return updatedReception, nil
}

func (r *ReceptionPostgres) ReceptionCreate(reception models.Reception) (models.Reception, error) {
	query := `
		INSERT INTO receptions (id, date_time, pvz_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, date_time, pvz_id, status
	`

	var createdReception models.Reception
	err := r.db.QueryRow(query, reception.Id, reception.DateTime, reception.PvzId, reception.Status).Scan(
		&createdReception.Id,
		&createdReception.DateTime,
		&createdReception.PvzId,
		&createdReception.Status,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return models.Reception{}, fmt.Errorf("failed to insert Receptions: %w", err)
	}

	return createdReception, nil
}
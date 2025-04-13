package repository

import (
	// "github.com/jmoiron/sqlx"
	"avito-testTask/models"
	"database/sql"
	// "log"
	"fmt"
	// "time"
	// "github.com/google/uuid"
)

type PVZPostgres struct {
	db *sql.DB
}

func NewPVZPostgres(db *sql.DB) *PVZPostgres {
	return &PVZPostgres{db: db}
}

func (r *PVZPostgres) PVZCreate(pvz models.PVZ) (models.PVZ, error) {
	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)
		RETURNING id, registration_date, city
	`

	var createdPVZ models.PVZ
	err := r.db.QueryRow(query, pvz.Id, pvz.RegistrationDate, pvz.City).Scan(
		&createdPVZ.Id,
		&createdPVZ.RegistrationDate,
		&createdPVZ.City,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return models.PVZ{}, fmt.Errorf("failed to insert PVZ: %w", err)
	}

	return createdPVZ, nil
}
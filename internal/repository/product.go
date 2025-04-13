package repository

import (
	// "github.com/jmoiron/sqlx"
	"avito-testTask/models"
	"database/sql"
	// "log"
	"fmt"
	"github.com/google/uuid"
)

type ProductPostgres struct {
	db *sql.DB
}

func NewProductPostgres(db *sql.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}


func (r *ProductPostgres) ProductCreate(product models.Product) (models.Product, error) {
	query := `
		INSERT INTO products (id, date_time, type, reception_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, date_time, type, reception_id
	`

	var createdProduct models.Product
	err := r.db.QueryRow(query, product.Id, product.DateTime, product.Type, product.ReceptionId).Scan(
		&createdProduct.Id,
		&createdProduct.DateTime,
		&createdProduct.Type,
		&createdProduct.ReceptionId,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return models.Product{}, fmt.Errorf("failed to insert Product: %w", err)
	}

	return createdProduct, nil
}

func (r *ProductPostgres) ProductDelete(PvzId uuid.UUID) error {
	query := `
		DELETE FROM products
		WHERE id = (
			SELECT id FROM products
			WHERE reception_id = (
				SELECT id FROM receptions
				WHERE pvz_id = $1
				ORDER BY date_time DESC
				LIMIT 1
			)
			ORDER BY date_time DESC
			LIMIT 1
		);
	`
	
	_, err := r.db.Exec(query, PvzId)
	if err != nil {
		fmt.Println("SQL error:", err)
		return fmt.Errorf("failed to delete Product: %w", err)
	}

	return nil
}

package repository

import (
	// "github.com/jmoiron/sqlx"
	"avito-testTask/models"
	"database/sql"
	// "log"
	"fmt"
	// "github.com/google/uuid"
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

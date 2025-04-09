package repository


import (
	"database/sql"
	// "avito-testTask/models"
)

type Authorization interface {
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		// Authorization: NewAuthPostgres(db),
	}
}
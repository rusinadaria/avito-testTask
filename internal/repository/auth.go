package repository

import (
	// "github.com/jmoiron/sqlx"
	// "avito-testTask/models"
	"database/sql"
	// "log"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
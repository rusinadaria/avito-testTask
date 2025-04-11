package repository

import (
	"avito-testTask/models"
	"database/sql"
	"github.com/google/uuid"
	// "avito-testTask/models"
)

type PVZ interface {
	PVZCreate(models.PVZ) (models.PVZ, error)
}

type Reception interface {
	ReceptionCreate(models.Reception) (models.Reception, error)
	GetLastReceptionByPVZ(uuid.UUID) (models.Reception, error)
	GetLastReceptionStatus(pvzId uuid.UUID) (models.Status, error)
	CloseReception(pvzId uuid.UUID) (models.Reception, error)
}

type Product interface {
	ProductCreate(product models.Product) (models.Product, error)
}

type Repository struct {
	PVZ
	Reception
	Product
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		PVZ: NewPVZPostgres(db),
		Reception: NewReceptionPostgres(db),
		Product: NewProductPostgres(db),
	}
}
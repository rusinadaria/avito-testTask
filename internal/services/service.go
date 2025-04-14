package services

import (
	"avito-testTask/internal/repository"
	"avito-testTask/models"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type PVZ interface {
	CreatePVZ(models.PVZ) (models.PVZ, error)
	GetPvzList(startDate *time.Time, endDate *time.Time, page int, limit int) ([]models.PVZWithReceptions, error)

}

type Reception interface {
	CreateReception(uuid.UUID) (models.Reception, error)
	CheckReception(pvzId uuid.UUID) (models.Reception, error)
}

type Product interface {
	AddProduct(Type models.Type, PvzId uuid.UUID) (models.Product, error)
	DeleteProduct(PvzId uuid.UUID) error
}

type Service struct {
	PVZ
	Reception
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		PVZ: NewPVZService(repos.PVZ),
		Reception: NewReceptionService(repos.Reception),
		Product: NewProductService(repos.Product, repos.Reception),
	}
}

package services

import (
	"avito-testTask/internal/repository"
	"avito-testTask/models"
	"fmt"
	"time"
)

type PVZService struct {
	repo repository.PVZ
}

func NewPVZService(repo repository.PVZ) *PVZService {
	return &PVZService{repo: repo}
}


func (s *PVZService) CreatePVZ(pvz models.PVZ) (models.PVZ, error) {
	pvz, err := s.repo.PVZCreate(pvz)
	if err != nil {
		fmt.Println("Ошибка в сервисе при создании ПВЗ")
		return models.PVZ{}, err
	}
	return pvz, nil
}

func (s *PVZService) GetPvzList(startDate *time.Time, endDate *time.Time, page int, limit int) ([]models.PVZWithReceptions, error) {
	pvzList, err := s.repo.GetPvz(startDate, endDate, limit, page)
	if err != nil {
		fmt.Println("Ошибка в сервисе при создании ПВЗ")
		return []models.PVZWithReceptions{}, err
	}
	return pvzList, nil
}
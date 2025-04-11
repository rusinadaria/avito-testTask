package services

import (
	"avito-testTask/internal/repository"
	"avito-testTask/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	// "fmt"
)

type ReceptionService struct {
	repo repository.Reception
}

func NewReceptionService(repo repository.Reception) *ReceptionService {
	return &ReceptionService{repo: repo}
}

func (s *ReceptionService) CreateReception(pvzId uuid.UUID) (models.Reception, error) {
	// Проверка на последнюю приёмку
	lastReception, err := s.repo.GetLastReceptionByPVZ(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	if lastReception.Status == models.InProgress {
		return models.Reception{}, fmt.Errorf("невозможно создать новую приёмку: последняя не закрыта")
	}

	var reception models.Reception
	reception.Id = uuid.New()
	reception.DateTime = time.Now().UTC()
	reception.PvzId = pvzId
	reception.Status = models.InProgress
	createdReception, err := s.repo.ReceptionCreate(reception)
	if err != nil {
		fmt.Println("Не удалось создать приемку в сервисе")
		return models.Reception{}, err
	}
	return createdReception, nil
}

func (s *ReceptionService) CheckReception(pvzId uuid.UUID) (models.Reception, error) {
	// Проверка на статус последней приемки
	status, err := s.repo.GetLastReceptionStatus(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	if status == models.Close {
		fmt.Println("Приемка уже закрыта")
		return models.Reception{}, err
	} else {
		updatedReception, err := s.repo.CloseReception(pvzId)
		if err != nil {
			fmt.Println("Ошибка при обновлении статуса приемки")
			return models.Reception{}, err
		} else {
			return updatedReception, nil
		}
	}
}
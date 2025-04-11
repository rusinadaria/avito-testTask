package services

import (
	// "avito-testTask/internal/handlers"
	"avito-testTask/internal/repository"
	"avito-testTask/models"
	"fmt"
	"time"
	"github.com/google/uuid"
	// "fmt"
)

type ProductService struct {
	repo repository.Product
	repos repository.Reception
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}



func (s *ProductService) AddProduct(Type models.Type, PvzId uuid.UUID) (models.Product, error) {

	// При этом товар должен привязываться к последнему незакрытому приёму товаров в рамках текущего ПВЗ.

	// if checkStatusLastReception(pvzId) == in_progress {
	// 	addProduct(requestProduct)
	// }

	// status, err := s.repo.GetLastReceptionStatus(requestProduct.PvzId)
	lastReception, err := s.repos.GetLastReceptionByPVZ(PvzId)
	if err != nil {
		fmt.Println("Ошибка при попытке получить статус в сервисе")
		return models.Product{}, err
	} else {
		if lastReception.Status == models.Close {
			fmt.Println("Не удалось добавить товар, приемка закрыта")
			return models.Product{}, err
		} else {
			var product models.Product
			product.Id = uuid.New()
			product.DateTime = time.Now().UTC()
			product.Type = Type
			product.ReceptionId = lastReception.Id

			createdProduct, err := s.repo.ProductCreate(product)
			if err != nil {
				fmt.Println("Ошибка при попытке добавить товар")
				return models.Product{}, err
			} else {
				return createdProduct, nil
			}
		}
	}

	// Если же нет новой незакрытой приёмки товаров, то в таком случае должна возвращаться ошибка, и товар не должен добавляться в систему.

	// if checkStatusLastReception == close {
	// 	return err
	// }


	// Если последняя приёмка товара все ещё не была закрыта, то результатом должна стать привязка товара к текущему ПВЗ и текущей приёмке с последующем добавлением данных в хранилище.

}
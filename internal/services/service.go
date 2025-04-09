package services

import (
	"avito-testTask/internal/repository"
	// "avito-testTask/models"
)

type Auth interface {
	CreateUser(username string, password string) (int, error)
}

type Service struct {
	Auth
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		// Auth: NewAuthService(repos.Authorization),
	}
}

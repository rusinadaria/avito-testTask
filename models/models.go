package models

import (
	"time"
	"github.com/google/uuid"
)

type Token string

type Role string

const (
	RoleEmployee Role = "employee"
	RoleModerator Role = "moderator"
)

type User struct {
	// Id string `json:"id"`
	Id uuid.UUID `json:"id"`
	Email string `json:"email"`
	Role  Role `json:"role"`
}

type City string

const (
	Moscow City = "Москва"
	SPB City = "Санкт-Петербург"
	Kazan City = "Казань"
)

type PVZ struct {
	Id uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City City `json:"city"`
}

type Status string

const (
	InProgress Status = "in_progress"
	Close Status = "close"
)

type Reception struct {
	Id uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId uuid.UUID `json:"pvzId"`
	Status Status `json:"status"`
}

type Type string

const (
	Electronics Type = "электроника"
	Clothes Type = "одежда"
	Shoes Type = "обувь"
)

type Product struct {
	Id uuid.UUID `json:"id"`
	DateTime    time.Time  `json:"dateTime"`
	Type Type `json:"type"`
	ReceptionId uuid.UUID  `json:"receptionId"`
}

type Error struct {
	Message string `json:"message"`
}

type ReceptionWithProducts struct {
	Reception Reception   `json:"reception"`
	Products  []Product   `json:"products"`
}

type PVZWithReceptions struct {
	PVZ        PVZ               `json:"pvz"`
	Receptions []ReceptionWithProducts  `json:"receptions"`
}

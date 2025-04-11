package models

import (
	"time"
	"github.com/google/uuid"
)

// Token:
//       type: string

type Token string

//     User:
//       type: object
//       properties:
//         id:
//           type: string
//           format: uuid
//         email:
//           type: string
//           format: email
//         role:
//           type: string
//           enum: [employee, moderator]
//       required: [email, role]

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

//     PVZ:
//       type: object
//       properties:
//         id:
//           type: string
//           format: uuid
//         registrationDate:
//           type: string
//           format: date-time
//         city:
//           type: string
//           enum: [Москва, Санкт-Петербург, Казань]
//       required: [city]

type City string

const (
	Moscow City = "Москва"
	SPB City = "Санкт-Петербург"
	Kazan City = "Казань"
)

type PVZ struct {
	// Id string `json:"id"`
	Id uuid.UUID `json:"id"`
	// RegistrationDate string `json:"registrationDate"`
	RegistrationDate time.Time `json:"registrationDate"`
	City City `json:"city"`
}

//     Reception:
//       type: object
//       properties:
//         id:
//           type: string
//           format: uuid
//         dateTime:
//           type: string
//           format: date-time
//         pvzId:
//           type: string
//           format: uuid
//         status:
//           type: string
//           enum: [in_progress, close]
//       required: [dateTime, pvzId, status]

type Status string

const (
	InProgress Status = "in_progress"
	Close Status = "close"
)

type Reception struct {
	// Id string `json:"id"`
	Id uuid.UUID `json:"id"`
	// DateTime string `json:"dateTime"`
	DateTime time.Time `json:"dateTime"`
	// PvzId string `json:"pvzId"`
	PvzId uuid.UUID `json:"pvzId"`
	Status Status `json:"status"`
}

//     Product:
//       type: object
//       properties:
//         id:
//           type: string
//           format: uuid
//         dateTime:
//           type: string
//           format: date-time
//         type:
//           type: string
//           enum: [электроника, одежда, обувь]
//         receptionId:
//           type: string
//           format: uuid
//       required: [type, receptionId]

type Type string

const (
	Electronics Type = "электроника"
	Clothes Type = "одежда"
	Shoes Type = "обувь"
)

type Product struct {
	// Id string `json:"id"`
	Id uuid.UUID `json:"id"`
	// DateTime string `json:"dateTime"`
	DateTime    time.Time  `json:"dateTime"`
	Type Type `json:"type"`
	// ReceptionId string `json:"receptionId"`
	ReceptionId uuid.UUID  `json:"receptionId"`
}

//     Error:
//       type: object
//       properties:
//         message:
//           type: string
//       required: [message]

type Error struct {
	Message string `json:"message"`
}
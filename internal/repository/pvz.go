package repository

import (
	// "github.com/jmoiron/sqlx"
	"avito-testTask/models"
	"database/sql"
	// "log"
	"fmt"
	"time"
	"github.com/google/uuid"
)

type PVZPostgres struct {
	db *sql.DB
}

func NewPVZPostgres(db *sql.DB) *PVZPostgres {
	return &PVZPostgres{db: db}
}

func (r *PVZPostgres) PVZCreate(pvz models.PVZ) (models.PVZ, error) {
	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)
		RETURNING id, registration_date, city
	`

	var createdPVZ models.PVZ
	err := r.db.QueryRow(query, pvz.Id, pvz.RegistrationDate, pvz.City).Scan(
		&createdPVZ.Id,
		&createdPVZ.RegistrationDate,
		&createdPVZ.City,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return models.PVZ{}, fmt.Errorf("failed to insert PVZ: %w", err)
	}

	return createdPVZ, nil
}

// func (r *PVZPostgres) GetPvz() ([]models.PVZWithReceptions, error) {
// 	query := `
// 		SELECT
// 			p.id AS pvz_id,
// 			p.registration_date,
// 			p.city,
// 			r.id AS reception_id,
// 			r.date_time AS reception_date_time,
// 			r.status AS reception_status,
// 			pr.id AS product_id,
// 			pr.date_time AS product_date_time,
// 			pr.type AS product_type
// 		FROM pvz p
// 		LEFT JOIN receptions r ON r.pvz_id = p.id
// 		LEFT JOIN products pr ON pr.reception_id = r.id
// 		ORDER BY p.city, p.registration_date DESC, r.date_time DESC, pr.date_time DESC;
// 	`

// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	type receptionKey struct {
// 		PvzID uuid.UUID
// 		RecID uuid.UUID
// 	}

// 	pvzMap := make(map[uuid.UUID]*models.PVZWithReceptions)
// 	receptionMap := make(map[receptionKey]*models.ReceptionWithProducts)

// 	for rows.Next() {
// 		var (
// 			pvzID           uuid.UUID
// 			pvzRegDate      time.Time
// 			pvzCity         string
// 			receptionID     sql.NullString
// 			receptionTime   sql.NullTime
// 			receptionStatus sql.NullString
// 			productID       sql.NullString
// 			productTime     sql.NullTime
// 			productType     sql.NullString
// 		)

// 		if err := rows.Scan(
// 			&pvzID,
// 			&pvzRegDate,
// 			&pvzCity,
// 			&receptionID,
// 			&receptionTime,
// 			&receptionStatus,
// 			&productID,
// 			&productTime,
// 			&productType,
// 		); err != nil {
// 			return nil, err
// 		}

// 		pvz, exists := pvzMap[pvzID]
// 		if !exists {
// 			pvz = &models.PVZWithReceptions{
// 				PVZ: models.PVZ{
// 					Id:               pvzID,
// 					RegistrationDate: pvzRegDate,
// 					City:             models.City(pvzCity),
// 				},
// 				Receptions: []models.ReceptionWithProducts{},
// 			}
// 			pvzMap[pvzID] = pvz
// 		}

// 		if receptionID.Valid {
// 			recID, _ := uuid.Parse(receptionID.String)
// 			key := receptionKey{PvzID: pvzID, RecID: recID}

// 			reception, exists := receptionMap[key]
// 			if !exists {
// 				reception = &models.ReceptionWithProducts{
// 					Reception: models.Reception{
// 						Id:       recID,
// 						DateTime: receptionTime.Time,
// 						PvzId:    pvzID,
// 						Status:   models.Status(receptionStatus.String),
// 					},
// 					Products: []models.Product{},
// 				}
// 				receptionMap[key] = reception
// 				pvz.Receptions = append(pvz.Receptions, *reception)
// 			}

// 			if productID.Valid {
// 				prodID, _ := uuid.Parse(productID.String)
// 				product := models.Product{
// 					Id:          prodID,
// 					DateTime:    productTime.Time,
// 					Type:        models.Type(productType.String),
// 					ReceptionId: recID,
// 				}
// 				reception.Products = append(reception.Products, product)
// 			}
// 		}
// 	}

// 	var result []models.PVZWithReceptions
// 	for _, pvz := range pvzMap {
// 		result = append(result, *pvz)
// 	}

// 	return result, nil
// }


func (r *PVZPostgres) GetPvz(startDate, endDate *time.Time, limit, offset int) ([]models.PVZWithReceptions, error) {
	query := `
		SELECT
			p.id AS pvz_id,
			p.registration_date,
			p.city,
			r.id AS reception_id,
			r.date_time AS reception_date_time,
			r.status AS reception_status,
			pr.id AS product_id,
			pr.date_time AS product_date_time,
			pr.type AS product_type
		FROM pvz p
		LEFT JOIN receptions r ON r.pvz_id = p.id
		LEFT JOIN products pr ON pr.reception_id = r.id
		WHERE 
			($1::timestamp IS NULL OR r.date_time >= $1) AND
			($2::timestamp IS NULL OR r.date_time <= $2)
		ORDER BY p.city, p.registration_date DESC, r.date_time DESC, pr.date_time DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type receptionKey struct {
		PvzID uuid.UUID
		RecID uuid.UUID
	}

	pvzMap := make(map[uuid.UUID]*models.PVZWithReceptions)
	receptionMap := make(map[receptionKey]*models.ReceptionWithProducts)

	for rows.Next() {
		var (
			pvzID           uuid.UUID
			pvzRegDate      time.Time
			pvzCity         string
			receptionID     sql.NullString
			receptionTime   sql.NullTime
			receptionStatus sql.NullString
			productID       sql.NullString
			productTime     sql.NullTime
			productType     sql.NullString
		)

		if err := rows.Scan(
			&pvzID,
			&pvzRegDate,
			&pvzCity,
			&receptionID,
			&receptionTime,
			&receptionStatus,
			&productID,
			&productTime,
			&productType,
		); err != nil {
			return nil, err
		}

		pvz, exists := pvzMap[pvzID]
		if !exists {
			pvz = &models.PVZWithReceptions{
				PVZ: models.PVZ{
					Id:               pvzID,
					RegistrationDate: pvzRegDate,
					City:             models.City(pvzCity),
				},
				Receptions: []models.ReceptionWithProducts{},
			}
			pvzMap[pvzID] = pvz
		}

		if receptionID.Valid {
			recID, _ := uuid.Parse(receptionID.String)
			key := receptionKey{PvzID: pvzID, RecID: recID}

			reception, exists := receptionMap[key]
			if !exists {
				reception = &models.ReceptionWithProducts{
					Reception: models.Reception{
						Id:       recID,
						DateTime: receptionTime.Time,
						PvzId:    pvzID,
						Status:   models.Status(receptionStatus.String),
					},
					Products: []models.Product{},
				}
				receptionMap[key] = reception
				pvz.Receptions = append(pvz.Receptions, *reception)
			}

			if productID.Valid {
				prodID, _ := uuid.Parse(productID.String)
				product := models.Product{
					Id:          prodID,
					DateTime:    productTime.Time,
					Type:        models.Type(productType.String),
					ReceptionId: recID,
				}
				reception.Products = append(reception.Products, product)
			}
		}
	}

	var result []models.PVZWithReceptions
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return result, nil
}




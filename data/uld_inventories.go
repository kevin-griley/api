package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevin-griley/api/db"
)

type UldInventory struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	UldNumber           string    `json:"uld_number"`
	UldType             string    `json:"uld_type"`
	UldStatus           string    `json:"uld_status"`
	CurrentLocationId   uuid.UUID `json:"current_location_id"`
	CurrentLocationType string    `json:"current_location_type"`
}

func CreateUldInventory(u *UldInventory) (uuid.UUID, error) {
	query := `
		INSERT INTO uld_inventories (id, created_at, updated_at, uld_number, uld_type, uld_status, current_location_id, current_location_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := db.Psql.QueryRow(
		query,
		u.ID,
		u.CreatedAt,
		u.UpdatedAt,
		u.UldNumber,
		u.UldType,
		u.UldStatus,
		u.CurrentLocationId,
		u.CurrentLocationType,
	).Scan(&u.ID)
	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

func GetUldInventoryByUldNumber(uldNumber string) (*UldInventory, error) {
	rows, err := db.Psql.Query(`SELECT * FROM uld_inventories WHERE uld_number = $1`, uldNumber)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUldInventory(rows)
	}
	return nil, fmt.Errorf("uld_inventory %s not found", uldNumber)
}

func NewUldInventory(uldNumber, uldType, uldStatus, currentLocationType string, currentLocationId uuid.UUID) (*UldInventory, error) {
	return &UldInventory{
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		UldNumber:           uldNumber,
		UldType:             uldType,
		UldStatus:           uldStatus,
		CurrentLocationId:   currentLocationId,
		CurrentLocationType: currentLocationType,
	}, nil
}

func scanIntoUldInventory(rows *sql.Rows) (*UldInventory, error) {
	u := new(UldInventory)
	err := rows.Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.UldNumber,
		&u.UldType,
		&u.UldStatus,
		&u.CurrentLocationId,
		&u.CurrentLocationType,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

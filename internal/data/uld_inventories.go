package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func ScanIntoUldInventory(rows *sql.Rows) (*UldInventory, error) {
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

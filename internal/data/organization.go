package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *organizationStoreImpl) CreateRequest(name, address, contactInfo string, organizationType OrganizationType) (*Organization, error) {

	orgId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	uniqueURL := GenerateRandomString(10)

	return &Organization{
		ID:               orgId,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		Name:             name,
		UniqueURL:        uniqueURL,
		Address:          address,
		ContactInfo:      contactInfo,
		OrganizationType: organizationType,
	}, nil
}

func (s *organizationStoreImpl) CreateOrganization(o *Organization) (*Organization, error) {

	data := map[string]any{
		"id":                o.ID,
		"created_at":        o.CreatedAt,
		"updated_at":        o.UpdatedAt,
		"name":              o.Name,
		"unique_url":        o.UniqueURL,
		"address":           o.Address,
		"contact_info":      o.ContactInfo,
		"organization_type": o.OrganizationType,
	}

	query, values, err := BuildInsertQuery("organizations", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrganization(rows)
	}

	return nil, fmt.Errorf("failed to create organization")

}

func (s *organizationStoreImpl) UpdateRequest(name, uniqueURL, address, contactInfo string, organizationType OrganizationType) (*Organization, error) {

	o := new(Organization)

	if name != "" {
		o.Name = name
	}
	if uniqueURL != "" {
		o.UniqueURL = uniqueURL
	}
	if address != "" {
		o.Address = address
	}
	if contactInfo != "" {
		o.ContactInfo = contactInfo
	}
	if organizationType != "" {
		o.OrganizationType = organizationType
	}

	return o, nil
}

func (s *organizationStoreImpl) UpdateOrganization(o *Organization) (*Organization, error) {

	updateData := make(map[string]any)
	updateData["updated_at"] = time.Now().UTC()

	if o.Name != "" {
		updateData["name"] = o.Name
	}
	if o.UniqueURL != "" {
		updateData["unique_url"] = o.UniqueURL
	}
	if o.Address != "" {
		updateData["address"] = o.Address
	}
	if o.ContactInfo != "" {
		updateData["contact_info"] = o.ContactInfo
	}
	if o.OrganizationType != "" {
		updateData["organization_type"] = o.OrganizationType
	}

	conditions := map[string]any{
		"id": o.ID,
	}

	query, values, err := BuildUpdateQuery("organizations", updateData, conditions)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoOrganization(rows)
	}

	return nil, fmt.Errorf("failed to update organization")

}

func (s *organizationStoreImpl) GetOrganizationByID(ID uuid.UUID) (*Organization, error) {

	data := map[string]any{
		"ID": ID,
	}

	query, values, err := BuildSelectQuery("organizations", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrganization(rows)
	}

	return nil, fmt.Errorf("organization %s not found", ID)

}

func (s *organizationStoreImpl) GetOrganizationByName(name string) (*Organization, error) {

	data := map[string]any{
		"name": name,
	}

	query, values, err := BuildSelectQuery("organizations", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrganization(rows)
	}

	return nil, fmt.Errorf("organization %s not found", name)

}

func (s *organizationStoreImpl) GetOrganizationByUniqueURL(uniqueURL string) (*Organization, error) {

	data := map[string]any{
		"unique_url": uniqueURL,
	}

	query, values, err := BuildSelectQuery("organizations", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrganization(rows)
	}

	return nil, fmt.Errorf("organization %s not found", uniqueURL)

}

type organizationStoreImpl struct {
	db *sql.DB
}

var NewOrganizationStore = func(db *sql.DB) OrganizationStore {
	return &organizationStoreImpl{
		db: db,
	}
}

type OrganizationStore interface {
	GetOrganizationByID(ID uuid.UUID) (*Organization, error)
	GetOrganizationByName(name string) (*Organization, error)
	GetOrganizationByUniqueURL(uniqueURL string) (*Organization, error)

	CreateOrganization(o *Organization) (*Organization, error)
	CreateRequest(name, address, contactInfo string, organizationType OrganizationType) (*Organization, error)

	UpdateOrganization(o *Organization) (*Organization, error)
	UpdateRequest(name, uniqueURL, address, contactInfo string, organizationType OrganizationType) (*Organization, error)
}

type OrganizationType string

const (
	Airline   OrganizationType = "Airline"
	Carrier   OrganizationType = "Carrier"
	Warehouse OrganizationType = "Warehouse"
)

type Organization struct {
	ID               uuid.UUID        `json:"id"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Name             string           `json:"name"`
	UniqueURL        string           `json:"unique_url"`
	Address          string           `json:"address"`
	ContactInfo      string           `json:"contact_info"`
	OrganizationType OrganizationType `json:"organization_type"`
}

func scanIntoOrganization(rows *sql.Rows) (*Organization, error) {
	var o Organization
	err := rows.Scan(
		&o.ID,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Name,
		&o.Address,
		&o.ContactInfo,
		&o.OrganizationType,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

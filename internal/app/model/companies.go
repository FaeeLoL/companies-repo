// Package model defines the core data structures used across the application
package model

import "github.com/google/uuid"

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

type Company struct {
	ID             uuid.UUID   `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description,omitempty"`
	EmployeesCount int         `json:"employees_count"`
	Registered     bool        `json:"registered"`
	Type           CompanyType `json:"type"`
}

type CreateCompanyData struct {
	ID             uuid.UUID
	Name           string
	Description    string
	EmployeesCount int
	Registered     bool
	Type           CompanyType
}

type UpdateCompanyData struct {
	ID             *uuid.UUID   `json:"id,omitempty"`
	Name           *string      `json:"name,omitempty"`
	Description    *string      `json:"description,omitempty"`
	EmployeesCount *int         `json:"employees_count,omitempty"`
	Registered     *bool        `json:"registered,omitempty"`
	Type           *CompanyType `json:"type,omitempty"`
}

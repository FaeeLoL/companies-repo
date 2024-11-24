// Package repositories provides data access layers for interacting with the database.
package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/faeelol/companies-store/internal/app/apperrors"
	"github.com/faeelol/companies-store/internal/app/model"
)

type CompanyRepository interface {
	CreateCompany(ctx context.Context, tx *sqlx.Tx, company *Company) error
	GetCompany(ctx context.Context, tx *sqlx.Tx, reqUUID uuid.UUID, name string) (model.Company, error)
	DeleteCompany(ctx context.Context, tx *sqlx.Tx, reqUUID uuid.UUID, name string) error
	UpdateCompany(ctx context.Context, tx *sqlx.Tx, updates model.UpdateCompanyData) error
}

type companyRepository struct {
}

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

func (c CompanyType) toDTO() model.CompanyType {
	return model.CompanyType(c)
}

type Company struct {
	ID             uuid.UUID   `db:"id"`
	Name           string      `db:"name"`
	Description    *string     `db:"description"`
	EmployeesCount int         `db:"employees_count"`
	Registered     bool        `db:"registered"`
	Type           CompanyType `db:"type"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

func (c Company) toDTO() model.Company {
	return model.Company{
		ID:             c.ID,
		Name:           c.Name,
		Description:    *c.Description,
		EmployeesCount: c.EmployeesCount,
		Registered:     c.Registered,
		Type:           c.Type.toDTO(),
	}
}

func NewCompanyRepository() CompanyRepository {
	return &companyRepository{}
}

// CreateCompany creates a new company in the database
func (r *companyRepository) CreateCompany(ctx context.Context, tx *sqlx.Tx, company *Company) error {
	query := `
		INSERT INTO companies (id, name, description, employees_count, registered, type, created_at, updated_at)
		VALUES (:id, :name, :description, :employees_count, :registered, :type, :created_at, :updated_at)
	`

	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	_, err := tx.NamedExecContext(ctx, query, company)
	if err != nil {

		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return apperrors.NewBadRequestError("duplicate key violation: unique constraint failed")
			}
		}
		return apperrors.NewInternalServerError("failed to create company").WithCause(err)
	}
	return nil
}

// GetCompany retrieves a company from the database by UUID or name.
// If both are empty, returns an error.
func (r *companyRepository) GetCompany(
	ctx context.Context,
	tx *sqlx.Tx,
	reqUUID uuid.UUID,
	name string,
) (model.Company, error) {
	var query string
	var args []interface{}

	switch {
	case reqUUID != uuid.Nil:
		query = `
			SELECT id, name, description, employees_count, registered, type, created_at, updated_at
			FROM companies
			WHERE id = $1
		`
		args = append(args, reqUUID)
	case name != "":
		query = `
			SELECT id, name, description, employees_count, registered, type, created_at, updated_at
			FROM companies
			WHERE name = $1
		`
		args = append(args, name)
	default:
		return model.Company{}, apperrors.NewBadRequestError("either uuid or name must be provided")
	}

	var company Company
	err := tx.GetContext(ctx, &company, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Company{}, apperrors.NewNotFoundError("company not found")
		}
		return model.Company{}, apperrors.NewInternalServerError("failed to query company").WithCause(err)
	}

	return company.toDTO(), nil
}

// DeleteCompany removes a company from the database by UUID or name.
// If both are empty, returns an error.
func (r *companyRepository) DeleteCompany(
	ctx context.Context,
	tx *sqlx.Tx,
	reqUUID uuid.UUID,
	name string,
) error {
	var query string
	var args []interface{}

	switch {
	case reqUUID != uuid.Nil:
		query = `
			DELETE FROM companies
			WHERE id = $1
		`
		args = append(args, reqUUID)
	case name != "":
		query = `
			DELETE FROM companies
			WHERE name = $1
		`
		args = append(args, name)
	default:
		return apperrors.NewBadRequestError("either uuid or name must be provided")
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return apperrors.NewInternalServerError("failed to delete company").WithCause(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalServerError("failed to get rows affected").WithCause(err)
	}

	if rowsAffected == 0 {
		return apperrors.NewNotFoundError("company not found")
	}

	return nil
}

func (r *companyRepository) UpdateCompany(
	ctx context.Context,
	tx *sqlx.Tx,
	updates model.UpdateCompanyData,
) error {
	var query string
	var args []any

	var setClauses []string
	if updates.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *updates.Description)
	}
	if updates.EmployeesCount != nil {
		setClauses = append(setClauses, "employees_count = ?")
		args = append(args, *updates.EmployeesCount)
	}
	if updates.Registered != nil {
		setClauses = append(setClauses, "registered = ?")
		args = append(args, *updates.Registered)
	}
	if updates.Type != nil {
		setClauses = append(setClauses, "type = ?")
		args = append(args, *updates.Type)
	}

	if len(setClauses) == 0 {
		return apperrors.NewBadRequestError("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	switch {
	case updates.ID != nil:
		query = fmt.Sprintf("UPDATE companies SET %s WHERE id = ?", strings.Join(setClauses, ", "))
		args = append(args, *updates.ID)
	case updates.Name != nil:
		query = fmt.Sprintf("UPDATE companies SET %s WHERE name = ?", strings.Join(setClauses, ", "))
		args = append(args, *updates.Name)
	default:
		return apperrors.NewBadRequestError("either uuid or name must be provided")
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return apperrors.NewInternalServerError("failed to update company").WithCause(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternalServerError("failed to get rows affected").WithCause(err)
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFoundError("company not found")
	}

	return nil
}

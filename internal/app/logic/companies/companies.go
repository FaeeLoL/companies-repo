// Package companies contains business logic for managing company operations.
package companies

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/faeelol/companies-store/internal/app/database"
	"github.com/faeelol/companies-store/internal/app/database/repositories"
	"github.com/faeelol/companies-store/internal/app/kafka"
	"github.com/faeelol/companies-store/internal/app/model"
	"github.com/faeelol/companies-store/internal/app/rest/middlewares"
)

type EventsProducer interface {
	Publish(ctx context.Context, key string, value []byte) error
}

type Controller struct {
	db          *sqlx.DB
	companyRepo repositories.CompanyRepository
	producer    EventsProducer
}

func NewCompaniesController(db *sqlx.DB, companyRepo repositories.CompanyRepository, producer EventsProducer) *Controller {
	return &Controller{
		db:          db,
		companyRepo: companyRepo,
		producer:    producer,
	}
}

func (c *Controller) CreateCompany(ctx context.Context, company model.CreateCompanyData) error {
	err := database.WithinTransaction(ctx, c.db, func(tx *sqlx.Tx) error {
		return c.companyRepo.CreateCompany(ctx, tx, &repositories.Company{
			ID:             company.ID,
			Name:           company.Name,
			Description:    &company.Description,
			EmployeesCount: company.EmployeesCount,
			Registered:     company.Registered,
			Type:           repositories.CompanyType(company.Type),
		})
	})
	if err != nil {
		return err
	}

	c.PublishEvent(ctx, kafka.CreateCompanyEvent, company.ID.String(), "name", company)

	return err
}

func (c *Controller) GetCompany(ctx context.Context, reqUUID uuid.UUID, name string) (model.Company, error) {
	company := model.Company{}

	err := database.WithinTransaction(ctx, c.db, func(tx *sqlx.Tx) error {
		var txErr error
		company, txErr = c.companyRepo.GetCompany(ctx, tx, reqUUID, name)
		return txErr
	})

	return company, err
}

func (c *Controller) DeleteCompany(ctx context.Context, reqUUID uuid.UUID, name string) error {
	err := database.WithinTransaction(ctx, c.db, func(tx *sqlx.Tx) error {
		return c.companyRepo.DeleteCompany(ctx, tx, reqUUID, name)
	})
	if err != nil {
		return err
	}

	if reqUUID != uuid.Nil {
		c.PublishEvent(ctx, kafka.DeleteCompanyEvent, reqUUID.String(), "uuid", map[string]string{})
	} else if name != "" {
		c.PublishEvent(ctx, kafka.DeleteCompanyEvent, name, "name", map[string]string{})
	}

	return nil
}

func (c *Controller) UpdateCompany(
	ctx context.Context,
	updates model.UpdateCompanyData,
) error {
	err := database.WithinTransaction(ctx, c.db, func(tx *sqlx.Tx) error {
		return c.companyRepo.UpdateCompany(ctx, tx, updates)
	})
	if err != nil {
		return err
	}

	if updates.ID != nil {
		c.PublishEvent(ctx, kafka.UpdateCompanyEvent, updates.ID.String(), "uuid", updates)
	} else if updates.Name != nil {
		c.PublishEvent(ctx, kafka.UpdateCompanyEvent, *updates.Name, "name", updates)
	}

	return nil
}

func (c *Controller) PublishEvent(ctx context.Context, action string, identifier string, idType string, data any) {
	event := map[string]any{
		"action":     action,
		"identifier": identifier,
		"id_type":    idType,
		"data":       data,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		middlewares.GetLoggerFromContext(ctx).Errorf("failed to serialize event: %v", err)
		return
	}

	if err := c.producer.Publish(ctx, action, eventBytes); err != nil {
		middlewares.GetLoggerFromContext(ctx).Errorf("failed to publish event to Kafka: %v", err)
	}
}

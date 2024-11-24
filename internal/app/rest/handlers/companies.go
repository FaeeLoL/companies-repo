package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"

	"github.com/faeelol/companies-store/internal/app/apperrors"
	"github.com/faeelol/companies-store/internal/app/model"
	"github.com/faeelol/companies-store/internal/app/rest/middlewares"
)

type CreateCompaniesController interface {
	CreateCompany(ctx context.Context, company model.CreateCompanyData) error
}

type CreateCompaniesHandler struct {
	schema *gojsonschema.Schema
	ccr    CreateCompaniesController
}

func NewCreateCompaniesHandler(ccr CreateCompaniesController) *CreateCompaniesHandler {
	return &CreateCompaniesHandler{
		schema: mustJSONSchema(createCompaniesSchema),
		ccr:    ccr,
	}
}

type CreateCompanyRequest struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	EmployeesCount int       `json:"employees_count"`
	Registered     bool      `json:"registered"`
	Type           string    `json:"type"`
}

func (ccr *CreateCompanyRequest) ToDTO() model.CreateCompanyData {
	return model.CreateCompanyData{
		ID:             ccr.ID,
		Name:           ccr.Name,
		Description:    ccr.Description,
		EmployeesCount: ccr.EmployeesCount,
		Registered:     ccr.Registered,
		Type:           model.CompanyType(ccr.Type),
	}
}

func (h *CreateCompaniesHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLoggerFromContext(r.Context())
	ctx := r.Context()

	var company CreateCompanyRequest
	err := ParseRequestJSON(r, h.schema, &company)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	err = h.ccr.CreateCompany(ctx, company.ToDTO())
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	RespondCodeAndJSON(rw, http.StatusCreated, nil, nil)
}

type GetCompaniesController interface {
	GetCompany(ctx context.Context, reqUUID uuid.UUID, name string) (model.Company, error)
}

type GetCompaniesHandler struct {
	gcc GetCompaniesController
}

func NewGetCompaniesHandler(gcc GetCompaniesController) *GetCompaniesHandler {
	return &GetCompaniesHandler{
		gcc: gcc,
	}
}

func (h *GetCompaniesHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLoggerFromContext(r.Context())
	ctx := r.Context()

	name, err := getStringParam(r, "name", false)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	id, err := getUUIDParam(r, "uuid", false)
	if err != nil {
		RespondError(rw, err, logger)
	}

	if name == "" && id == uuid.Nil {
		RespondError(rw, apperrors.NewBadRequestError("name or uuid should be provided"), logger)
		return
	}

	res, err := h.gcc.GetCompany(ctx, id, name)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	RespondCodeAndJSON(rw, http.StatusOK, res, nil)
}

type DeleteCompaniesController interface {
	DeleteCompany(ctx context.Context, reqUUID uuid.UUID, name string) error
}

type DeleteCompaniesHandler struct {
	dcc DeleteCompaniesController
}

func NewDeleteCompaniesHandler(dcc DeleteCompaniesController) *DeleteCompaniesHandler {
	return &DeleteCompaniesHandler{
		dcc: dcc,
	}
}

func (h *DeleteCompaniesHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLoggerFromContext(r.Context())
	ctx := r.Context()

	name, err := getStringParam(r, "name", false)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	id, err := getUUIDParam(r, "uuid", false)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	if name == "" && id == uuid.Nil {
		RespondError(rw, apperrors.NewBadRequestError("name or uuid should be provided"), logger)
		return
	}

	err = h.dcc.DeleteCompany(ctx, id, name)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	RespondCodeAndJSON(rw, http.StatusNoContent, nil, nil)
}

type PatchCompaniesController interface {
	UpdateCompany(ctx context.Context, updates model.UpdateCompanyData) error
}

type PatchCompaniesHandler struct {
	pcc    PatchCompaniesController
	schema *gojsonschema.Schema
}

func NewPatchCompaniesHandler(pcc PatchCompaniesController) *PatchCompaniesHandler {
	return &PatchCompaniesHandler{
		pcc:    pcc,
		schema: mustJSONSchema(patchCompaniesSchema),
	}
}

type PatchCompanyRequest struct {
	ID             *uuid.UUID `json:"id"`
	Name           *string    `json:"name"`
	Description    *string    `json:"description,omitempty"`
	EmployeesCount *int       `json:"employees_count"`
	Registered     *bool      `json:"registered"`
	Type           *string    `json:"type"`
}

func (ccr *PatchCompanyRequest) ToDTO() model.UpdateCompanyData {
	var companyType *model.CompanyType
	if ccr.Type != nil {
		convertedCompanyType := model.CompanyType(*ccr.Type)
		companyType = &convertedCompanyType
	}
	return model.UpdateCompanyData{
		ID:             ccr.ID,
		Name:           ccr.Name,
		Description:    ccr.Description,
		EmployeesCount: ccr.EmployeesCount,
		Registered:     ccr.Registered,
		Type:           companyType,
	}
}

func (h *PatchCompaniesHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLoggerFromContext(r.Context())
	ctx := r.Context()

	var updates model.UpdateCompanyData

	err := ParseRequestJSON(r, h.schema, &updates)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	if updates.Name == nil && updates.Description == nil && updates.EmployeesCount == nil && updates.Registered == nil && updates.Type == nil {
		RespondError(rw, apperrors.NewBadRequestError("no fields to update"), logger)
		return
	}

	err = h.pcc.UpdateCompany(ctx, updates)
	if err != nil {
		RespondError(rw, err, logger)
		return
	}

	RespondCodeAndJSON(rw, http.StatusOK, nil, nil)
}

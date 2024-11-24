// Package handlers provides HTTP request handlers for the application,
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"

	"github.com/faeelol/companies-store/internal/app/apperrors"
)

const (
	ContentTypeAppJSON = "application/json"
)

func RespondCodeAndJSON(rw http.ResponseWriter, statusCode int, respData interface{}, logger logrus.FieldLogger) {
	if respData == nil {
		rw.WriteHeader(statusCode)
		return
	}

	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", ContentTypeAppJSON)
	}

	respJSON, err := json.Marshal(respData)
	if err != nil {
		if logger != nil {
			logger.Error(logger.WithField("error", err), "error while marshaling json for response body")
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(statusCode)
	if _, err = rw.Write(respJSON); err != nil {
		if logger != nil {
			logger.Error(logger.WithField("error", err), "error while writing response body")
		}
	}
}

func CheckJSONContentType(r *http.Request) error {
	reqContentType := r.Header.Get("Content-Type")
	if reqContentType == "" {
		return errors.New("Content-Type header is not set")
	}

	contentType, _, err := mime.ParseMediaType(reqContentType)
	if err != nil {
		return fmt.Errorf("failed to parse Content-Type header for request: %s", err)
	}

	if contentType != ContentTypeAppJSON {
		return fmt.Errorf("Content-Type header is not %s", ContentTypeAppJSON)
	}

	return nil
}

func ParseRequestJSON(r *http.Request, schema *gojsonschema.Schema, dst any) error {
	if err := CheckJSONContentType(r); err != nil {
		return apperrors.NewBadRequestError(err.Error())
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apperrors.NewBadRequestError("failed to read request body")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Sprintf("failed to validate JSON schema: %v", err))
	}

	if !result.Valid() {
		var validationErrors []string
		for _, err := range result.Errors() {
			validationErrors = append(validationErrors, err.String())
		}
		return &ValidationError{
			Errors: validationErrors,
		}
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return apperrors.NewBadRequestError("failed to parse JSON body")
	}

	return nil
}

// ValidationError list of errors
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, ";\n")
}

func mustJSONSchema(js []byte) *gojsonschema.Schema {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(js))
	if err != nil {
		panic(err)
	}
	return schema
}

func RespondError(w http.ResponseWriter, err error, logger logrus.FieldLogger) {
	appErr := apperrors.MapToAppError(err)

	if logger != nil && appErr.Code >= http.StatusInternalServerError {
		logger.WithField("error", err).WithField("cause", appErr.Err).Error("internal server error")
	}

	RespondCodeAndJSON(w, appErr.Code, map[string]string{"error": appErr.Message}, logger)
}

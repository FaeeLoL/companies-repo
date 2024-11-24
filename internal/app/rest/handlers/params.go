package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/faeelol/companies-store/internal/app/apperrors"
)

func getStringParam(r *http.Request, name string, isRequired bool) (string, error) {
	v := r.URL.Query().Get(name)
	if v == "" {
		if isRequired {
			return "", apperrors.NewBadRequestError(fmt.Sprintf("missing %s param", name))
		}

		return "", nil
	}

	return v, nil
}

func getUUIDParam(r *http.Request, key string, isRequired bool) (uuid.UUID, error) {
	rawUUID, err := getStringParam(r, key, isRequired)
	if err != nil {
		return uuid.Nil, err
	}
	if rawUUID == "" {
		return uuid.Nil, nil
	}
	resUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, apperrors.NewBadRequestError(fmt.Sprintf("invalid %s param", key))
	}
	return resUUID, nil
}

package httputils

import (
	"actualiteco/actualite/internal/core/domain"
	"fmt"
	"net/http"
	"strconv"
)

func GetUintQueryParam(r *http.Request, key string) (*uint64, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to convert param='%s' by key='%s' in query parameters: %w",
			param,
			key,
			domain.ErrInvalidRequest,
		)
	}

	return &val, nil
}

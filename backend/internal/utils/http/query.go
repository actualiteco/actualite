package httputils

import (
	"actualiteco/actualite/internal/core/domain"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pageQueryValueKey     = "page"
	pageSizeQueryValueKey = "page_size"
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

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}
	return &param
}

func GetPaginationParamsFromQuery(r *http.Request) (*uint64, *uint64, error) {
	page, err := GetUintQueryParam(r, pageQueryValueKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get page from query param: %w", err)
	}

	pageSize, err := GetUintQueryParam(r, pageSizeQueryValueKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get page size from query param: %w", err)
	}

	return page, pageSize, nil
}

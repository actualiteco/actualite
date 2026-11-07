package httprequest

import (
	"actualiteco/actualite/internal/core/domain"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, domain.ErrInvalidRequest)
	}

	err := requestValidator.Struct(dest)
	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, domain.ErrInvalidRequest)
	}

	return nil
}

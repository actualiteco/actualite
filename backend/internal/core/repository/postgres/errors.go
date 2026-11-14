package postgres_pool

import (
	"actualiteco/actualite/internal/core/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

func MapUniqueViolation(err error, constraints map[string]error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != uniqueViolation {
		return err
	}

	if domainErr, exists := constraints[pgErr.ConstraintName]; exists {
		return fmt.Errorf("%s: %w", pgErr.ConstraintName, domainErr)
	}

	return fmt.Errorf("unique violation (%s): %w", pgErr.ConstraintName, domain.ErrConflict)
}

package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
)

var userConstraints = map[string]error{
	"idx_users_username": domain.ErrUsernameAlreadyUsed,
	"idx_users_email":    domain.ErrEmailAlreadyUsed,
}

type UsersRepository struct {
	ConnPool postgres_pool.Pool
}

func NewAuthRepository(conn postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{
		ConnPool: conn,
	}
}

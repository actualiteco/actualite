package search_repository

import postgres_pool "actualiteco/actualite/internal/core/repository/postgres"

type SearchRepository struct {
	pool postgres_pool.Pool
}

func New(pool postgres_pool.Pool) *SearchRepository {
	return &SearchRepository{pool: pool}
}

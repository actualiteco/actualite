package notification_repository

import postgres_pool "actualiteco/actualite/internal/core/repository/postgres"

type NotificationsRepository struct {
	ConnPool postgres_pool.Pool
}

func New(connPool postgres_pool.Pool) *NotificationsRepository {
	return &NotificationsRepository{
		ConnPool: connPool,
	}
}

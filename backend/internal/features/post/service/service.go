package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"io"
	"uuid"
)

type PostsService struct {
	postsRepo PostsRepository
	txManager txManager
	storage   fileStorage
}

type PostsRepository interface {
	Like(ctx context.Context, userID, postID uuid.UUID) error
	Unlike(ctx context.Context, userID, postID uuid.UUID) error
	CreatePost(ctx context.Context, post domain.Post, tags []domain.Tag) (domain.Post, error)
	GetPostByID(ctx context.Context, userID, postID uuid.UUID) (domain.PostWithAuthor, error)
	GetPostsByAuthorID(ctx context.Context, userID, authorID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
	GetReplies(ctx context.Context, userID, parentID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
	GetPostsFeed(ctx context.Context, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
	Delete(ctx context.Context, userID, postID uuid.UUID) error

	GetCountOfPostsByUserID(ctx context.Context, userID uuid.UUID) (uint64, bool, error)
	GetCountOfReplies(ctx context.Context, parentID uuid.UUID) (uint64, bool, error)
	GetCountOfPostsForFeed(ctx context.Context, userID uuid.UUID) (uint64, bool, error)
}

type tagsRepository interface {
}

type txManager interface {
	DoWrite(ctx context.Context, fn func(context.Context) error) error
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type fileStorage interface {
	Upload(ctx context.Context, filename string, file io.Reader) (string, error)
}

func New(postsRepo PostsRepository, txManager txManager, storage fileStorage) *PostsService {
	return &PostsService{
		postsRepo: postsRepo,
		txManager: txManager,
		storage:   storage,
	}
}

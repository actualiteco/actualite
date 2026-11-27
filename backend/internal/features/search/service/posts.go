package search_service

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

func (s *SearchService) SearchPostsByContent(ctx context.Context, query string, userID uuid.UUID,
	page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	if query = strings.TrimSpace(query); query == "" {
		return nil, domain.Pagination{}, fmt.Errorf("query is required: %w", domain.ErrInvalidRequest)
	}

	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	var posts []domain.PostWithAuthor
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		count, err := s.repo.CountPostsWithContent(ctx, query)
		if err != nil {
			return err
		}
		pagination.ParseCount(count)
		if pagination.Page > pagination.TotalPages {
			return nil
		}

		posts, err = s.repo.SearchPosts(ctx, query, userID, limit, offset)
		return err
	})
	return posts, pagination, err
}

func (s *SearchService) SearchPostsByHashtags(ctx context.Context, query string, userID uuid.UUID,
	page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	tags, err := extractTagsFromQuery(query)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("extract tags from query: %v: %w", err, domain.ErrInvalidRequest)
	}

	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	var posts []domain.PostWithAuthor
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		count, err := s.repo.CountPostsWithHashtags(ctx, tags)
		if err != nil {
			return err
		}
		pagination.ParseCount(count)
		if pagination.Page > pagination.TotalPages {
			return nil
		}

		posts, err = s.repo.SearchPostsByHashtags(ctx, tags, userID, limit, offset)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("search posts by hashtags in tx: %w", err)
	}

	return posts, pagination, nil
}

func extractTagsFromQuery(query string) ([]domain.Tag, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	tagsStrs := strings.Split(query, ",")
	for i, tagStr := range tagsStrs {
		tagsStrs[i] = strings.TrimSpace(strings.ToLower(tagStr))
	}

	tags, err := domain.StringsToTags(tagsStrs)
	if err != nil {
		return nil, fmt.Errorf("convert strings to tags: %w", err)
	}

	return tags, nil
}

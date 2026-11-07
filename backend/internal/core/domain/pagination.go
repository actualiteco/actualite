package domain

import (
	"fmt"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 30
)

type Pagination struct {
	Page       uint64
	PageSize   uint64
	Total      uint64
	TotalPages uint64
}

func NewPagination(page, pageSize *uint64) Pagination {
	var p uint64 = DefaultPage
	if page != nil {
		p = *page
	}

	var ps uint64 = DefaultPageSize
	if pageSize != nil {
		ps = *pageSize
	}

	return Pagination{
		Page:     p,
		PageSize: ps,
	}
}

func (p *Pagination) GetOffsetLimit() (uint64, uint64, error) {
	if p.Page == 0 || p.PageSize == 0 {
		return 0, 0, fmt.Errorf("invalid page=%d page_size=%d to get offset, limit", p.Page, p.PageSize)
	}

	offset := (p.Page - 1) * p.PageSize

	return offset, p.PageSize, nil
}

func (p *Pagination) ParseCount(count uint64) {
	p.Total = count
	if p.PageSize == 0 {
		p.TotalPages = 0
		return
	}
	p.TotalPages = (count + p.PageSize - 1) / p.PageSize
}

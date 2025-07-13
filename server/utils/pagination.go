package utils

import (
	"math"
)

type DefaultQueryParams struct {
	Page      int    `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sortBy" json:"sortBy"`
	SortOrder string `form:"sortOrder" json:"sortOrder"`
}

// Pagination represents pagination configuration
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Offset     int `json:"offset"`
}

func BuildPagination(page, limit, total int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	offset := (page - 1) * limit

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Offset:     offset,
	}
}

// optional methods
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

func (p *Pagination) NextPage() int {
	if p.HasNext() {
		return p.Page + 1
	}
	return p.Page
}

func (p *Pagination) PrevPage() int {
	if p.HasPrev() {
		return p.Page - 1
	}
	return p.Page
}

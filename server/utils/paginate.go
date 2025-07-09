package utils

import "github.com/fiqrioemry/asset_management_system_app/server/dto"

func Paginate(total int64, currentPage, limit int) *dto.PaginationResponse {
	if limit <= 0 {
		limit = 10
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginationResponse{
		CurrentPage: currentPage,
		Limit:       limit,
		TotalItems:  int(total),
		TotalPages:  totalPages,
	}
}

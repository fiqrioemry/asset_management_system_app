// ==================== utils/response.go ====================
package utils

import (
	"net/http"

	"github.com/fiqrioemry/asset_management_system_app/server/errors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Code    errors.ErrorCode `json:"code"`
	Errors  map[string]any   `json:"errors,omitempty"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Pagination  *Pagination     `json:"pagination,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	Flags       map[string]bool `json:"flags,omitempty"`
}

// HTTP ERROR RESPONSES ===============
func HandleError(c *gin.Context, err error) {
	LogError(c, err)

	if appErr, ok := errors.IsAppError(err); ok {
		response := ErrorResponse{
			Success: false,
			Message: appErr.Message,
			Code:    appErr.Code,
		}

		if appErr.Context != nil {
			if errorDetails, exists := appErr.Context["errors"]; exists {
				response.Errors = errorDetails.(map[string]any)
			}
		}

		c.JSON(appErr.HTTPStatus, response)
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: "Internal server error",
		Code:    errors.ErrCodeInternalServer,
	})
}

// HTTP SUCCESS RESPONSES ===============
func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

func OKWithPagination(c *gin.Context, message string, data any, pagination *Pagination) {
	SuccessWithMeta(c, http.StatusOK, message, data, &Meta{
		Pagination: pagination,
	})
}

func OKWithPaginationAndPermissions(c *gin.Context, message string, data any, pagination *Pagination, permissions map[string]bool) {
	SuccessWithMeta(c, http.StatusOK, message, data, &Meta{
		Pagination:  pagination,
		Permissions: permissions,
	})
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta *Meta) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

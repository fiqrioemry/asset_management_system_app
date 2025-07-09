package utils

import (
	"net/http"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message"`
	Code      errors.ErrorCode `json:"code"`
	Timestamp time.Time        `json:"timestamp"`
	Path      string           `json:"path"`
	Context   map[string]any   `json:"context,omitempty"`
}

func HandleError(c *gin.Context, err error) {
	logger := GetLogger()

	logFields := []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	}

	if appErr, ok := errors.IsAppError(err); ok {
		logFields = append(logFields,
			zap.String("error_code", string(appErr.Code)),
			zap.String("error_message", appErr.Message),
		)

		if errors.IsServerError(err) {
			logFields = append(logFields, zap.Error(appErr.Err))
			logger.Error("Server error occurred", logFields...)
		} else {
			logger.Warn("Client error occurred", logFields...)
		}

		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Success:   false,
			Message:   appErr.Message,
			Code:      appErr.Code,
			Timestamp: time.Now(),
			Path:      c.Request.URL.Path,
			Context:   appErr.Context,
		})
		return
	}

	// Handle unknown errors
	logFields = append(logFields,
		zap.Error(err),
		zap.String("error_type", "unknown"),
	)
	logger.Error("Unknown error occurred", logFields...)

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success:   false,
		Message:   "Internal server error",
		Code:      errors.ErrCodeInternalServer,
		Timestamp: time.Now(),
		Path:      c.Request.URL.Path,
	})
}

package middlewares

import (
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		err := response.NewInternalServerError("Internal server error", recovered.(error))
		response.Error(c, err)
		c.Abort()
	})
}

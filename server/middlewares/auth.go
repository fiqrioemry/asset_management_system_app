package middlewares

import (
	"github.com/fiqrioemry/asset_management_system_app/server/utils"
	"github.com/fiqrioemry/go-api-toolkit/response"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("accessToken")
		if err != nil || tokenString == "" {
			response.Error(c, response.NewUnauthorized("Unauthorized!! Token missing"))
			c.Abort()
			return
		}

		claims, err := utils.DecodeAccessToken(tokenString)
		if err != nil {
			response.Error(c, response.NewUnauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}

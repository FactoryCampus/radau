package pkg

import (
	"github.com/gin-gonic/gin"
)

func HandleApiKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		if apiKey != validKey {
			c.AbortWithStatus(403)
		}

		c.Set("apiKey", apiKey)
	}
}

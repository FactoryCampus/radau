package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func IsValidApiKey(ctx *gin.Context, validKey string) error {
	apiKey := ctx.GetHeader("Authorization")
	if apiKey != validKey {
		return errors.New("Invalid API key")
	}

	ctx.Set("apiKey", validKey)

	return nil
}

func HandleApiKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := IsValidApiKey(c, validKey); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

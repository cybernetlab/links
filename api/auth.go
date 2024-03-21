package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func AuthMiddleware(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		respondWithError(c, 401, "Bearer authorization required")
		return
	}
	bearer := strings.SplitN(strings.TrimSpace(auth), " ", 2)
	if strings.ToLower(bearer[0]) != "bearer" {
		respondWithError(c, 401, "Bearer authorization required")
		return
	}
	c.Next()
}

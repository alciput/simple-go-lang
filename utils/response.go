package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONError sends a JSON error response with the given status and message.
// Keeps error responses consistent and avoids leaking internal details.
func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

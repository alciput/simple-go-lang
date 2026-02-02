package middleware

import (
	"net/http"
	"strings"

	"app/config"
	"app/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims holds JWT payload. Use minimal claims; add "sub" (user id) if needed later.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Auth validates Bearer JWT and sets email in context. Use after this middleware: c.GetString("email").
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}
		tokenStr := parts[1]
		var claims Claims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}

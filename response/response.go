package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Envelope is the consistent API response format. Use concrete Data types for Swagger.
type Envelope struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// RegisterData is the data payload for POST /register success.
type RegisterData struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// LoginData is the data payload for POST /login success.
type LoginData struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// UserPublic is user without password_hash for API responses (profile and login user).
type UserPublic struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterResponse is the full response for POST /register (for Swagger).
type RegisterResponse struct {
	Data  *RegisterData `json:"data,omitempty"`
	Error string        `json:"error,omitempty"`
}

// LoginResponse is the full response for POST /login (for Swagger).
type LoginResponse struct {
	Data  *LoginData `json:"data,omitempty"`
	Error string     `json:"error,omitempty"`
}

// ProfileResponse is the full response for GET /profile (for Swagger).
type ProfileResponse struct {
	Data  *UserPublic `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// ErrorResponse is the envelope for error responses (for Swagger).
type ErrorResponse struct {
	Data  interface{} `json:"data"`
	Error string       `json:"error"`
}

// OK sends 200 with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Data: data})
}

// Created sends 201 with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Data: data})
}

// Error sends status with error message. No sensitive details.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{Error: message})
}

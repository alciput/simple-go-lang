package handlers

import (
	"errors"
	"net/http"
	"time"

	"app/config"
	"app/models"
	"app/response"
	"app/store"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest is the JSON body for POST /register.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is the JSON body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Auth holds dependencies for register and login.
type Auth struct {
	Repo   store.UserRepository
	Config *config.Config
}

// Register creates a user. Password is hashed with bcrypt; only hash is stored.
//
//	@Summary		Register a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RegisterRequest	true	"Register payload"
//	@Success		201		{object}	response.RegisterResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		409		{object}	response.ErrorResponse
//	@Router			/register [post]
func (h *Auth) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to process password")
		return
	}
	id, err := h.Repo.CreateUser(c.Request.Context(), req.Email, string(hash))
	if err != nil {
		if errors.Is(err, store.ErrDuplicateEmail) {
			response.Error(c, http.StatusConflict, "email already registered")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create user")
		return
	}
	response.Created(c, response.RegisterData{ID: id, Email: req.Email})
}

// Login verifies password and returns a JWT. Same generic error for wrong email vs wrong password.
//
//	@Summary		Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login payload"
//	@Success		200		{object}	response.LoginResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Router			/login [post]
func (h *Auth) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input")
		return
	}
	u, err := h.Repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}
	if u == nil {
		response.Error(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"email": u.Email,
		"exp":   exp.Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create token")
		return
	}
	response.OK(c, response.LoginData{
		Token: tokenStr,
		User:  userToPublic(u),
	})
}

func userToPublic(u *models.User) response.UserPublic {
	return response.UserPublic{ID: u.ID, Email: u.Email, CreatedAt: u.CreatedAt}
}

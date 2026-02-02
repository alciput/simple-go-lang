package handlers

import (
	"net/http"

	"app/response"
	"app/store"

	"github.com/gin-gonic/gin"
)

// Profile holds dependencies for profile handler.
type Profile struct {
	Repo store.UserRepository
}

// Profile returns the current user's profile. Email comes from JWT middleware.
//
//	@Summary		Get current user profile
//	@Tags			auth
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	response.ProfileResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Router			/profile [get]
func (h *Profile) Profile(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := h.Repo.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}
	if u == nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}
	response.OK(c, userToPublic(u))
}

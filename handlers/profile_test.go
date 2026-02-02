package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/response"
	"app/store"

	"github.com/gin-gonic/gin"
)

func TestProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := store.NewMock()
	h := &Profile{Repo: repo}

	_, _ = repo.CreateUser(context.Background(), "me@b.com", "hash")

	t.Run("with email in context returns 200 and user", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
		c.Set("email", "me@b.com")

		h.Profile(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d; want 200", w.Code)
		}
		var env response.Envelope
		if err := json.NewDecoder(w.Body).Decode(&env); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if env.Error != "" {
			t.Errorf("error = %q; want empty", env.Error)
		}
		data, _ := env.Data.(map[string]interface{})
		if data["email"] != "me@b.com" {
			t.Errorf("data.email = %v; want me@b.com", data["email"])
		}
	})

	t.Run("without email in context returns 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
		// do not set email

		h.Profile(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d; want 401", w.Code)
		}
	})

	t.Run("email not in store returns 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
		c.Set("email", "missing@b.com")

		h.Profile(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d; want 404", w.Code)
		}
	})
}

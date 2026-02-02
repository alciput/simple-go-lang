package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/config"
	"app/response"
	"app/store"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := store.NewMock()
	cfg := &config.Config{JWTSecret: "test-secret"}
	h := &Auth{Repo: repo, Config: cfg}

	t.Run("valid body returns 201 and data", func(t *testing.T) {
		body := map[string]string{"email": "a@b.com", "password": "password123"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d; want 201", w.Code)
		}
		var env response.Envelope
		if err := json.NewDecoder(w.Body).Decode(&env); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if env.Error != "" {
			t.Errorf("error = %q; want empty", env.Error)
		}
		data, _ := env.Data.(map[string]interface{})
		if data["email"] != "a@b.com" {
			t.Errorf("data.email = %v; want a@b.com", data["email"])
		}
		if _, ok := data["id"]; !ok {
			t.Error("data.id missing")
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		_, _ = repo.CreateUser(context.Background(), "dup@b.com", "hash")
		body := map[string]string{"email": "dup@b.com", "password": "password123"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		if w.Code != http.StatusConflict {
			t.Errorf("status = %d; want 409", w.Code)
		}
		var env response.Envelope
		_ = json.NewDecoder(w.Body).Decode(&env)
		if env.Error == "" {
			t.Error("expected error message")
		}
	})

	t.Run("invalid body returns 400", func(t *testing.T) {
		body := map[string]string{"email": "bad"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d; want 400", w.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := store.NewMock()
	cfg := &config.Config{JWTSecret: "test-secret"}
	h := &Auth{Repo: repo, Config: cfg}

	realHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	_, _ = repo.CreateUser(context.Background(), "login@b.com", string(realHash))

	t.Run("valid credentials return 200 and token", func(t *testing.T) {
		body := map[string]string{"email": "login@b.com", "password": "password123"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

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
		if data["token"] == nil || data["token"] == "" {
			t.Error("data.token missing or empty")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		body := map[string]string{"email": "login@b.com", "password": "wrong"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d; want 401", w.Code)
		}
	})

	t.Run("unknown email returns 401", func(t *testing.T) {
		body := map[string]string{"email": "nobody@b.com", "password": "password123"}
		raw, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(raw))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d; want 401", w.Code)
		}
	})
}


package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"khalif-identify/internal/domain"
	"khalif-identify/internal/handler"
	"khalif-identify/internal/tests/mocks" // Import Mock yang kita buat

)

func TestLogin(t *testing.T) {
	// Setup Gin Mode Test (Supaya log bersih)
	gin.SetMode(gin.TestMode)

	t.Run("Success Login", func(t *testing.T) {
		// 1. Setup Mock (Panggil dari package mocks)
		mockUC := new(mocks.MockUserUseCase)

		// Data Dummy
		expectedToken := "token-rahasia-123"
		expectedUser := &domain.User{
			ID:    1,
			Email: "khalif@gmail.com",
			Name:  "Khalif",
		}

		// Ekspektasi: Jika login dipanggil dengan email benar -> Return Sukses
		mockUC.On("Login", "khalif@gmail.com", "password123").
			Return(expectedToken, expectedUser, nil)

		// 2. Init Handler (Inject Mock)
		h := handler.NewUserHandler(mockUC)

		// 3. Setup Router
		r := gin.Default()
		r.POST("/login", h.Login)

		// 4. Lakukan Request Pura-pura
		reqBody := []byte(`{"email":"khalif@gmail.com", "password":"password123"}`)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 5. Assert (Pengecekan Hasil)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Pastikan token & data user ada di response JSON
		assert.Equal(t, expectedToken, response["token"])
		assert.NotNil(t, response["data"]) // Data user tidak boleh null

		// Verifikasi mock terpanggil
		mockUC.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		// 1. Setup
		mockUC := new(mocks.MockUserUseCase)
		
		// Ekspektasi: Jika password salah -> Return Error
		mockUC.On("Login", "wrong@gmail.com", "wrongpass").
			Return("", nil, errors.New("invalid credentials"))

		h := handler.NewUserHandler(mockUC)
		r := gin.Default()
		r.POST("/login", h.Login)

		// 2. Request
		reqBody := []byte(`{"email":"wrong@gmail.com", "password":"wrongpass"}`)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUC.AssertExpectations(t)
	})
}
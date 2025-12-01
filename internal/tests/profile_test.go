package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"khalif-identify/internal/domain"
	"khalif-identify/internal/handler"
	"khalif-identify/internal/tests/mocks"

)

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Update Profile with Image", func(t *testing.T) {
		mockUC := new(mocks.MockUserUseCase)

		// Data Dummy
		userID := uint(1)
		updatedUser := &domain.User{ID: 1, Name: "Khalif Baru", PhoneNumber: "08999"}

		// Ekspektasi Mock:
		// Menggunakan mock.Anything untuk file karena pointer file sulit diprediksi di test
		mockUC.On("UpdateProfile", userID, "Khalif Baru", "08999", "", mock.Anything, mock.Anything).
			Return(updatedUser, nil)

		h := handler.NewUserHandler(mockUC)

		// Setup Router dengan Middleware Dummy untuk set User ID
		r := gin.Default()
		r.Use(func(c *gin.Context) {
			// Pura-pura AuthMiddleware sudah jalan dan set user_id (float64 dari JWT)
			c.Set("user_id", float64(userID)) 
			c.Next()
		})
		r.POST("/profile/update", h.UpdateProfile)

		// Buat Request Multipart (Form Data)
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		
		writer.WriteField("name", "Khalif Baru")
		writer.WriteField("phone", "08999")
		
		// Simulasi File Upload
		part, _ := writer.CreateFormFile("image", "avatar.jpg")
		part.Write([]byte("fake image content"))
		
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/profile/update", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Khalif Baru", data["name"])
		
		mockUC.AssertExpectations(t)
	})

	t.Run("Unauthorized (No Token)", func(t *testing.T) {
		mockUC := new(mocks.MockUserUseCase)
		h := handler.NewUserHandler(mockUC)

		r := gin.Default()
		// Tidak ada middleware yang set "user_id" disini
		r.POST("/profile/update", h.UpdateProfile)

		req, _ := http.NewRequest(http.MethodPost, "/profile/update", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Harusnya 401 karena c.Get("user_id") gagal
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Usecase Error", func(t *testing.T) {
		mockUC := new(mocks.MockUserUseCase)
		userID := uint(1)

		// Ekspektasi Error dari usecase
		mockUC.On("UpdateProfile", userID, "Khalif", "", "", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		h := handler.NewUserHandler(mockUC)

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", float64(userID))
		})
		r.POST("/profile/update", h.UpdateProfile)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("name", "Khalif")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/profile/update", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
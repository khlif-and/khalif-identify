package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"khalif-identify/internal/usecase"

)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(u usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: u}
}

func (h *UserHandler) Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	file, header, err := c.Request.FormFile("image")

	if err != nil && err != http.ErrMissingFile {
		log.Printf("[Register Failed] Upload Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload error"})
		return
	}

	user, err := h.useCase.Register(name, email, phone, password, file, header)
	if err != nil {
		log.Printf("[Register Failed] Usecase Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Register Success] Data: %+v", user)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) RegisterCustomer(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	file, header, err := c.Request.FormFile("image")

	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload error"})
		return
	}

	user, err := h.useCase.RegisterCustomer(name, email, phone, password, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Login Failed] Bind JSON Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, user, err := h.useCase.Login(input.Email, input.Password)
	if err != nil {
		log.Printf("[Login Failed] Auth Error (Email: %s): %v", input.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Login Success] User: %s", input.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"data":  user,
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.useCase.GetAllAdmins()
	if err != nil {
		log.Printf("[GetAll Failed] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[GetAll Success] Count: %d users retrieved", len(users))
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetCountryCodes(c *gin.Context) {
	countries := h.useCase.GetCountryCodes()
	c.JSON(http.StatusOK, gin.H{"data": countries})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	file, header, _ := c.Request.FormFile("image")

	updatedUser, err := h.useCase.UpdateProfile(userID, name, phone, password, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    updatedUser,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenSplit := strings.Split(authHeader, " ")
	if len(tokenSplit) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
		return
	}
	tokenString := tokenSplit[1]

	err := h.useCase.Logout(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	user, err := h.useCase.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
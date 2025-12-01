package domain

import (
	"mime/multipart" // <-- PENTING: Untuk handle file upload
	"time"

	"khalif-identify/pkg/utils" // <-- PENTING: Untuk struct Country

)

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	Email         string    `gorm:"uniqueIndex" json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Password      string    `json:"-"`
	ProfileImage  string    `json:"profile_image"`
	DominantColor string    `json:"dominant_color"`
	RoleID        uint      `json:"role_id"`
	Role          Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	FindAll() ([]User, error)
	Update(user *User) error
	CountByRoleID(roleID uint) (int64, error)
}

type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
}

// Interface Usecase ditaruh sini (Input Port)
type UserUseCase interface {
	// Perhatikan: Kita pakai *User (bukan *domain.User) karena sudah di dalam package domain
	Register(name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*User, error)
	Login(email, password string) (string, *User, error)
	Logout(tokenString string) error
	UpdateProfile(userID uint, name, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*User, error)
	GetAllAdmins() ([]User, error)
	GetCountryCodes() []utils.Country
}
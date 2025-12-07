package domain
import (
	"context" 
	"mime/multipart"
	"time"

	"khalif-identify/pkg/utils"

)
type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
type User struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	UUID          string    `gorm:"type:varchar(36);uniqueIndex" json:"user_id"`
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
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, page, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	CountByRoleID(ctx context.Context, roleID uint) (int64, error)
}
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
type UserUseCase interface {
	Register(ctx context.Context, name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*User, error)
	RegisterCustomer(ctx context.Context, name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*User, error)
	Login(ctx context.Context, email, password string) (string, *User, error)
	Logout(ctx context.Context, tokenString string) error
	UpdateProfile(ctx context.Context, userID uint, name, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*User, error)
	GetAllAdmins(ctx context.Context, page, limit int) ([]User, int64, error)
	GetCountryCodes() []utils.Country
	GetProfile(ctx context.Context, userID uint) (*User, error)
}
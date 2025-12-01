package repository

import (
	"gorm.io/gorm"

	"khalif-identify/internal/domain"

)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Role").Find(&users).Error
	return users, err
}

func (r *UserRepo) CountByRoleID(roleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

// --- UPDATE DISINI ---

func (r *UserRepo) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Role").First(&user, id).Error
	return &user, err
}

func (r *UserRepo) Update(user *domain.User) error {
	// Save akan mengupdate semua field struct ke database
	return r.db.Save(user).Error
}
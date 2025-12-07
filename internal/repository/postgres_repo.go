package repository
import (
	"context"

	"gorm.io/gorm"

	"khalif-identify/internal/domain"

)
type UserRepo struct {
	db *gorm.DB
}
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	return &user, err
}
func (r *UserRepo) FindAll(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).Preload("Role").
		Limit(limit).Offset(offset).
		Find(&users).Error
	return users, total, err
}
func (r *UserRepo) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}
func (r *UserRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").First(&user, id).Error
	return &user, err
}
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
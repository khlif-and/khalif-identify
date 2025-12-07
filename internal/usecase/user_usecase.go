package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"khalif-identify/internal/domain"
	"khalif-identify/pkg/utils"

)

type userUseCase struct {
	repo      domain.UserRepository
	cache     domain.CacheRepository
	uploader  *utils.AzureUploader
	jwtSecret string
}

func NewUserUseCase(repo domain.UserRepository, cache domain.CacheRepository, uploader *utils.AzureUploader, jwtSecret string) domain.UserUseCase {
	return &userUseCase{
		repo:      repo,
		cache:     cache,
		uploader:  uploader,
		jwtSecret: jwtSecret,
	}
}

func (u *userUseCase) GetCountryCodes() []utils.Country {
	return utils.GetCountryList()
}

func (u *userUseCase) Register(ctx context.Context, name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
	formattedPhone, err := utils.FormatPhoneNumber(phone, "ID")
	if err != nil {
		return nil, err
	}

	const MaxAdminCount = 3
	const TargetRoleID = 1
	const TargetRoleName = "Admin"

	currentCount, err := u.repo.CountByRoleID(ctx, TargetRoleID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengecek kuota admin: %v", err)
	}

	if currentCount >= MaxAdminCount {
		return nil, errors.New("kuota admin sudah penuh (maksimal 3 admin)")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	imgResult, err := utils.HandleProfileImageLogic(file, name, email)
	if err != nil {
		return nil, err
	}

	finalImageUrl := imgResult.AvatarURL

	if file != nil && fileHeader != nil {
		ext := filepath.Ext(fileHeader.Filename)
		filename := uuid.New().String() + ext

		uploadedUrl, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		finalImageUrl = uploadedUrl
	}

	user := &domain.User{
		UUID:          uuid.New().String(),
		Name:          name,
		Email:         email,
		PhoneNumber:   formattedPhone,
		Password:      hashedPassword,
		RoleID:        TargetRoleID,
		ProfileImage:  finalImageUrl,
		DominantColor: imgResult.DominantColor,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Role = domain.Role{ID: TargetRoleID, Name: TargetRoleName}
	u.cache.Del(ctx, "list_admins")
	return user, nil
}

func (u *userUseCase) RegisterCustomer(ctx context.Context, name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
	formattedPhone, err := utils.FormatPhoneNumber(phone, "ID")
	if err != nil {
		return nil, err
	}

	const CustomerRoleID = 3
	const CustomerRoleName = "Customer"

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	imgResult, err := utils.HandleProfileImageLogic(file, name, email)
	if err != nil {
		return nil, err
	}

	finalImageUrl := imgResult.AvatarURL

	if file != nil && fileHeader != nil {
		ext := filepath.Ext(fileHeader.Filename)
		filename := uuid.New().String() + ext
		uploadedUrl, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		finalImageUrl = uploadedUrl
	}

	user := &domain.User{
		UUID:          uuid.New().String(),
		Name:          name,
		Email:         email,
		PhoneNumber:   formattedPhone,
		Password:      hashedPassword,
		RoleID:        CustomerRoleID,
		ProfileImage:  finalImageUrl,
		DominantColor: imgResult.DominantColor,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Role = domain.Role{ID: CustomerRoleID, Name: CustomerRoleName}
	return user, nil
}

func (u *userUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	// PERBAIKAN: Gunakan user.UUID (string), bukan user.ID (uint)
	token, err := utils.GenerateToken(user.UUID, user.Role.Name, u.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *userUseCase) Logout(ctx context.Context, tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration in token")
	}

	expirationTime := time.Unix(int64(expFloat), 0)
	timeRemaining := time.Until(expirationTime)

	if timeRemaining < 0 {
		return nil
	}

	return u.cache.Set(ctx, "blacklist:"+tokenString, "revoked", timeRemaining)
}

func (u *userUseCase) GetAllAdmins(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	return u.repo.FindAll(ctx, page, limit)
}

func (u *userUseCase) UpdateProfile(ctx context.Context, userID uint, name, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}
	if phone != "" {
		formattedPhone, err := utils.FormatPhoneNumber(phone, "ID")
		if err != nil {
			return nil, err
		}
		user.PhoneNumber = formattedPhone
	}
	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if file != nil && fileHeader != nil {
		newColor, err := utils.ExtractDominantColor(file)
		if err == nil {
			user.DominantColor = newColor
		} else {
			user.DominantColor = "#000000"
		}

		file.Seek(0, io.SeekStart)
		ext := filepath.Ext(fileHeader.Filename)
		filename := uuid.New().String() + ext
		uploadedUrl, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		user.ProfileImage = uploadedUrl
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	u.cache.Del(ctx, "list_admins")
	return user, nil
}

func (u *userUseCase) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
	return u.repo.FindByID(ctx, userID)
}
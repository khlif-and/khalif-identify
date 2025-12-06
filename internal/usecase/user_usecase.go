package usecase

import (
	"encoding/json"
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

type UserUseCase interface {
	Register(name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error)
	RegisterCustomer(name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error)
	GetAllAdmins() ([]domain.User, error)
	GetCountryCodes() []utils.Country
	UpdateProfile(userID uint, name, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error)
	Logout(tokenString string) error
	GetProfile(userID uint) (*domain.User, error)
}

type userUseCase struct {
	repo      domain.UserRepository
	cache     domain.CacheRepository
	uploader  *utils.AzureUploader
	jwtSecret string
}

func NewUserUseCase(repo domain.UserRepository, cache domain.CacheRepository, uploader *utils.AzureUploader, jwtSecret string) UserUseCase {
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

func (u *userUseCase) Register(name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
	formattedPhone, err := utils.FormatPhoneNumber(phone, "ID")
	if err != nil {
		return nil, err
	}

	const MaxAdminCount = 3
	const TargetRoleID = 1
	const TargetRoleName = "Admin"

	currentCount, err := u.repo.CountByRoleID(TargetRoleID)
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
		Name:          name,
		Email:         email,
		PhoneNumber:   formattedPhone,
		Password:      hashedPassword,
		RoleID:        TargetRoleID,
		ProfileImage:  finalImageUrl,
		DominantColor: imgResult.DominantColor,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	user.Role = domain.Role{
		ID:   TargetRoleID,
		Name: TargetRoleName,
	}

	u.cache.Del("list_admins")
	return user, nil
}

func (u *userUseCase) RegisterCustomer(name, email, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
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
		Name:          name,
		Email:         email,
		PhoneNumber:   formattedPhone,
		Password:      hashedPassword,
		RoleID:        CustomerRoleID,
		ProfileImage:  finalImageUrl,
		DominantColor: imgResult.DominantColor,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	user.Role = domain.Role{
		ID:   CustomerRoleID,
		Name: CustomerRoleName,
	}

	return user, nil
}

func (u *userUseCase) Login(email, password string) (string, *domain.User, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Role.Name, u.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *userUseCase) Logout(tokenString string) error {
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

	err = u.cache.Set("blacklist:"+tokenString, "revoked", timeRemaining)
	return err
}

func (u *userUseCase) GetAllAdmins() ([]domain.User, error) {
	cachedData, err := u.cache.Get("list_admins")
	if err == nil {
		var users []domain.User
		json.Unmarshal([]byte(cachedData), &users)
		return users, nil
	}

	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(users)
	u.cache.Set("list_admins", jsonData, 5*time.Minute)

	return users, nil
}

func (u *userUseCase) UpdateProfile(userID uint, name, phone, password string, file multipart.File, fileHeader *multipart.FileHeader) (*domain.User, error) {
	user, err := u.repo.FindByID(userID)
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

	if err := u.repo.Update(user); err != nil {
		return nil, err
	}

	u.cache.Del("list_admins")

	return user, nil
}

func (u *userUseCase) GetProfile(userID uint) (*domain.User, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
package mocks

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"

	"khalif-identify/internal/domain"
	"khalif-identify/pkg/utils"

)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Login(email, password string) (string, *domain.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*domain.User), args.Error(2)
}

func (m *MockUserUseCase) Register(name, email, phone, password string, file multipart.File, fh *multipart.FileHeader) (*domain.User, error) {
	args := m.Called(name, email, phone, password, file, fh)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// --- UPDATE DISINI ---
func (m *MockUserUseCase) UpdateProfile(userID uint, name, phone, password string, file multipart.File, fh *multipart.FileHeader) (*domain.User, error) {
	// Kita gunakan mock.Called untuk merekam panggilan
	args := m.Called(userID, name, phone, password, file, fh)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) GetAllAdmins() ([]domain.User, error) {
	return nil, nil
}

func (m *MockUserUseCase) GetCountryCodes() []utils.Country {
	return nil
}
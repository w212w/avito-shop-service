package service_test

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	userData := args.Get(0)
	if userData == nil {
		return nil, args.Error(1)
	}
	return userData.(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateCoins(userID int, amount int) error {
	args := m.Called(userID, amount)
	return args.Error(0)
}

func TestRegisterSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo, "supersecretkey")

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	err := authService.Register("testuser", "password")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo, "supersecretkey")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	token, err := authService.Login("testuser", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo, "supersecretkey")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	token, err := authService.Login("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

package service

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	secretKey string
}

func NewAuthService(userRepo *repository.UserRepository, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, secretKey: secretKey}
}

// Register создает нового пользователя с хешированным паролем
func (s *AuthService) Register(username, password string) error {
	// Генерация хеша пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Создание пользователя с хешированным паролем
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword), // Используем PasswordHash
		Coins:        1000,                   // Начальный баланс, можно поменять на 0
	}

	// Сохраняем пользователя в базе данных
	return s.userRepo.CreateUser(user)
}

// Login выполняет проверку пользователя и создает JWT-токен
func (s *AuthService) Login(username, password string) (string, error) {
	// Получаем пользователя по имени
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Сравниваем хеш пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

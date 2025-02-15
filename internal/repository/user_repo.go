package repository

import (
	"avito-shop-service/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных
func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)",
		user.Username, user.PasswordHash, user.Coins)
	return err
}

// GetUserByUsername возвращает пользователя по логину
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, username, password_hash, coins, created_at FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UpdateCoins обновляет баланс пользователя
func (r *UserRepository) UpdateCoins(userID int, amount int) error {
	// Проверка на положительный баланс
	if amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}

	_, err := r.db.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, userID)
	return err
}

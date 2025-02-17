package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`     // Приватное поле для хранения хеша пароля
	Coins        int       `json:"coins"` // Баланс пользователя, при регистрации пользователь получает 1000 Coins
	CreatedAt    time.Time `json:"created_at"`
}

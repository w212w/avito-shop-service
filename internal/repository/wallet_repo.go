package repository

import (
	"avito-shop-service/internal/models"
	"database/sql"
	"errors"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Получение баланса пользователя (coins)
func (r *WalletRepository) GetBalance(userID int) (int, error) {
	var balance int
	err := r.db.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&balance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return balance, nil
}

// Обновление баланса пользователя (coins)
func (r *WalletRepository) UpdateBalance(userID int, amount int) error {
	_, err := r.db.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, userID)
	return err
}

// Перевод монет между пользователями
func (r *WalletRepository) Transfer(fromUserID, toUserID, amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Проверяем баланс отправителя
	var senderBalance int
	err = tx.QueryRow("SELECT coins FROM users WHERE id = $1", fromUserID).Scan(&senderBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	if senderBalance < amount {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	// Вычитаем монеты у отправителя
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, fromUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавляем монеты получателю
	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, toUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Записываем транзакцию в таблицу transactions
	_, err = tx.Exec(
		"INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)",
		fromUserID, toUserID, amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *WalletRepository) GetTransactions(userID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT id, from_user_id, to_user_id, amount, created_at 
		FROM transactions 
		WHERE from_user_id = $1 OR to_user_id = $1 
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

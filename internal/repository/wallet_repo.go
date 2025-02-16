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

// Получение цены товара из базы данных
func (r *WalletRepository) GetItemPrice(itemName string) (int, error) {
	var price int
	err := r.db.QueryRow("SELECT price FROM shop WHERE item = $1", itemName).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("item not found")
		}
		return 0, err
	}
	return price, nil
}

// Покупка товара
func (r *WalletRepository) PurchaseItem(userID int, itemName string, price int, quantity int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Проверяем баланс пользователя
	var userBalance int
	err = tx.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&userBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	totalPrice := price * quantity
	if userBalance < totalPrice {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	// Обновляем баланс пользователя
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", totalPrice, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Записываем покупку в таблицу purchases
	_, err = tx.Exec("INSERT INTO purchases (user_id, item, price, quantity) VALUES ($1, $2, $3, $4)", userID, itemName, price, quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершаем транзакцию
	return tx.Commit()
}

// Получение инвентаря пользователя
func (r *WalletRepository) GetInventory(userID int) ([]models.Item, error) {
	rows, err := r.db.Query("SELECT item, price, quantity FROM purchases WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Name, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

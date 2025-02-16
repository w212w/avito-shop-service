package models

import "time"

type Transaction struct {
	ID         int       `json:"id"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// InfoResponse - структура для ответа на запрос /api/info
type InfoResponse struct {
	Balance      int           `json:"balance"`
	Inventory    []Item        `json:"inventory"`    // Предполагаем, что Item - это структура для товара
	Transactions []Transaction `json:"transactions"` // Структура для транзакций
}

// Item - структура для представления предмета в инвентаре
type Item struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

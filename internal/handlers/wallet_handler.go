package handlers

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// Получение баланса пользователя
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.GetBalance(req.UserID)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"balance": balance})
}

// Пополнение баланса
func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Invalid deposit amount", http.StatusBadRequest)
		return
	}

	if err := h.walletService.Deposit(req.UserID, req.Amount); err != nil {
		http.Error(w, "Deposit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Перевод монет между пользователями
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из заголовков (устанавливается в middleware)
	userID := r.Header.Get("UserID")

	// Преобразуем userID в int
	fromUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим тело запроса
	var req struct {
		ToUserID int `json:"to_user_id"`
		Amount   int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Проверяем, что пользователь не переводит монеты самому себе
	if fromUserID == req.ToUserID {
		http.Error(w, "Cannot transfer to yourself", http.StatusBadRequest)
		return
	}

	// Выполняем перевод через сервис
	if err := h.walletService.Transfer(fromUserID, req.ToUserID, req.Amount); err != nil {
		http.Error(w, "Transfer failed", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
}

// Получение информации о истории транзакций
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из заголовка (должен быть установлен в middleware)
	userIDStr := r.Header.Get("UserID")

	// Преобразуем user_id в int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Получаем транзакции пользователя
	transactions, err := h.walletService.GetTransactions(userID)
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// BuyItem обрабатывает покупку товара
func (h *WalletHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	// Извлекаем item из параметров пути
	vars := mux.Vars(r)
	itemName := vars["item"]

	// Извлекаем количество товара из query параметров
	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		quantityStr = "1" // Если количество не указано, по умолчанию 1
	}

	// Преобразуем количество в int
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Получаем user_id из заголовков (должен быть установлен в middleware)
	userID := r.Header.Get("UserID")

	// Преобразуем user_id в int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Получаем цену товара из базы через WalletService
	itemPrice, err := h.walletService.GetItemPrice(itemName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Вызываем метод для выполнения покупки
	err = h.walletService.PurchaseItem(userIDInt, itemName, itemPrice, quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item purchased successfully",
	})
}

// Получение информации о монетах, инвентаре и истории транзакций
func (h *WalletHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из заголовков (должен быть установлен в middleware)
	userID := r.Header.Get("UserID")

	// Преобразуем user_id в int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Получаем баланс пользователя
	balance, err := h.walletService.GetBalance(userIDInt)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	// Получаем инвентарь пользователя
	inventory, err := h.walletService.GetInventory(userIDInt)
	if err != nil {
		http.Error(w, "Failed to get inventory", http.StatusInternalServerError)
		return
	}

	// Получаем транзакции пользователя
	transactions, err := h.walletService.GetTransactions(userIDInt)
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	infoResponse := models.InfoResponse{
		Balance:      balance,
		Inventory:    inventory,
		Transactions: transactions,
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(infoResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

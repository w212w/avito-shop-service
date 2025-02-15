package handlers

import (
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
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.GetBalance(userID)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"balance": balance})
}

// Пополнение баланса
func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.walletService.Deposit(userID, req.Amount); err != nil {
		http.Error(w, "Deposit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Перевод монет между пользователями
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID int `json:"from_user_id"` // Добавляем fromUserID в тело запроса
		ToUserID   int `json:"to_user_id"`
		Amount     int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.walletService.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		http.Error(w, "Transfer failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение информации о истории транзакций
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
	}

	// Декодируем `user_id` из body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем транзакции
	transactions, err := h.walletService.GetTransactions(req.UserID)
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

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

// Перевод монет между пользователями
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из заголовков (устанавливается в middleware)
	userID := r.Header.Get("UserID")

	fromUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	// Выполняем перевод
	if err := h.walletService.Transfer(fromUserID, req.ToUserID, req.Amount); err != nil {
		http.Error(w, "Transfer failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение информации о истории транзакций
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Header.Get("UserID")

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// BuyItem обрабатывает покупку товара
func (h *WalletHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemName := vars["item"]

	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		quantityStr = "1"
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	itemPrice, err := h.walletService.GetItemPrice(itemName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.walletService.PurchaseItem(userIDInt, itemName, itemPrice, quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Item purchased successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

// Получение информации о монетах, инвентаре и истории транзакций
func (h *WalletHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")

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

	infoResponse := models.InfoResponse{
		Balance:      balance,
		Inventory:    inventory,
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(infoResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

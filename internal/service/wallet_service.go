package service

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/repository"
	"errors"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

// Получение баланса пользователя
func (s *WalletService) GetBalance(userID int) (int, error) {
	return s.walletRepo.GetBalance(userID)
}

// Пополнение баланса (добавление монет)
func (s *WalletService) Deposit(userID int, amount int) error {
	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}
	return s.walletRepo.UpdateBalance(userID, amount)
}

// Перевод монет между пользователями
func (s *WalletService) Transfer(fromUserID, toUserID, amount int) error {
	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}
	return s.walletRepo.Transfer(fromUserID, toUserID, amount)
}

// Получение истории транзакций
func (s *WalletService) GetTransactions(userID int) ([]models.Transaction, error) {
	return s.walletRepo.GetTransactions(userID)
}

// Покупка товара
func (s *WalletService) PurchaseItem(userID int, itemName string, price int, quantity int) error {
	// Проверка на достаточность средств
	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return err
	}

	totalCost := price * quantity
	if balance < totalCost {
		return errors.New("insufficient funds")
	}

	// Обновление баланса и добавление записи о покупке
	return s.walletRepo.PurchaseItem(userID, itemName, price, quantity)
}

// Получение цены товара
func (s *WalletService) GetItemPrice(itemName string) (int, error) {
	return s.walletRepo.GetItemPrice(itemName)
}

// Получение инвентаря пользователя
func (s *WalletService) GetInventory(userID int) ([]models.Item, error) {
	items, err := s.walletRepo.GetInventory(userID)
	return items, err
}

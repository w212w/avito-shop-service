package service

import (
	"avito-shop-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock WalletRepository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetBalance(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockWalletRepository) Transfer(fromUserID, toUserID, amount int) error {
	args := m.Called(fromUserID, toUserID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) GetTransactions(userID int) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockWalletRepository) PurchaseItem(userID int, itemName string, price int, quantity int) error {
	args := m.Called(userID, itemName, price, quantity)
	return args.Error(0)
}

func (m *MockWalletRepository) GetItemPrice(itemName string) (int, error) {
	args := m.Called(itemName)
	return args.Int(0), args.Error(1)
}

func (m *MockWalletRepository) GetInventory(userID int) ([]models.Item, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Item), args.Error(1)
}

// Получение баланса
func TestGetBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	mockRepo.On("GetBalance", 10).Return(1000, nil)

	balance, err := service.GetBalance(10)

	assert.NoError(t, err)
	assert.Equal(t, 1000, balance)

	mockRepo.AssertExpectations(t)
}

// Перевод монет
func TestTransfer(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	mockRepo.On("Transfer", 1, 2, 300).Return(nil)

	err := service.Transfer(1, 2, 300)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ПокупкА товара
func TestPurchaseItem(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	mockRepo.On("GetBalance", 1).Return(1000, nil)
	mockRepo.On("PurchaseItem", 1, "T-Shirt", 200, 2).Return(nil)

	err := service.PurchaseItem(1, "T-Shirt", 200, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetInventory(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	expectedInventory := []models.Item{
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
	}

	mockRepo.On("GetInventory", 1).Return(expectedInventory, nil)

	inventory, err := service.GetInventory(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedInventory, inventory)

	mockRepo.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	expectedTransactions := []models.Transaction{
		{FromUserID: 1, ToUserID: 2, Amount: 200},
		{FromUserID: 2, ToUserID: 1, Amount: 500},
	}

	mockRepo.On("GetTransactions", 1).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactions(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	mockRepo.AssertExpectations(t)
}

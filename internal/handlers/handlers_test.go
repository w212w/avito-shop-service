package handlers

import (
	"avito-shop-service/config"
	"avito-shop-service/internal/middleware"
	"avito-shop-service/internal/repository"
	"avito-shop-service/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {

	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "shop",
		JWTSecret:  "supersecretkey",
	}
	db := repository.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewPostgresWalletRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	walletService := service.NewWalletService(walletRepo)

	authHandler := NewAuthHandler(authService)
	walletHandler := NewWalletHandler(walletService)

	router := mux.NewRouter()

	router.HandleFunc("/api/auth", authHandler.Auth).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("/info", walletHandler.GetInfo).Methods("GET")
	protected.HandleFunc("/sendCoin", walletHandler.Transfer).Methods("POST")
	protected.HandleFunc("/buy/{item}", walletHandler.BuyItem).Methods("POST")

	return router
}

// Получение токена
func getValidToken() string {
	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password",
	})

	req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(w, req)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		return ""
	}

	return resp["token"]
}

// Покупка товара
func TestPurchase(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"item":     "book",
		"price":    50,
		"quantity": 1,
	})
	req := httptest.NewRequest("POST", "/api/buy/book", bytes.NewBuffer(reqBody))
	validToken := getValidToken()
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Перевод монет между пользователями
func TestTransfer(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]int{
		"to_user_id": 2,
		"amount":     10,
	})
	req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(reqBody))
	validToken := getValidToken()
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Запрос информации о пользователе
func TestGetInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/info", nil)
	validToken := getValidToken()
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

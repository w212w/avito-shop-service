package main

import (
	"avito-shop-service/config"
	"avito-shop-service/internal/handlers"
	"avito-shop-service/internal/middleware"
	"avito-shop-service/internal/repository"
	"avito-shop-service/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := repository.ConnectDB(cfg)
	defer db.Close()

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewPostgresWalletRepository(db)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	walletService := service.NewWalletService(walletRepo)

	// Инициализируем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	walletHandler := handlers.NewWalletHandler(walletService)

	router := mux.NewRouter()

	// Применяем middleware для защищенных маршрутов
	// Регистрация и вход не требуют аутентификации
	router.HandleFunc("/api/auth", authHandler.Auth).Methods("POST")

	// Защищенные маршруты с middleware
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	// Роуты, которые требуют аутентификации
	protected.HandleFunc("/info", walletHandler.GetInfo).Methods("GET")
	protected.HandleFunc("/sendCoin", walletHandler.Transfer).Methods("POST")
	protected.HandleFunc("/buy/{item}", walletHandler.BuyItem).Methods("POST")

	log.Println("Server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

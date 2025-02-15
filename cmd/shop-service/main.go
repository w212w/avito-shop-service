package main

import (
	"avito-shop-service/config"
	"avito-shop-service/internal/handlers"
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

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	router := mux.NewRouter()

	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/wallet/{user_id}/balance", walletHandler.GetBalance).Methods("GET")
	router.HandleFunc("/api/wallet/{user_id}/deposit", walletHandler.Deposit).Methods("POST")
	router.HandleFunc("/api/wallet/transfer", walletHandler.Transfer).Methods("POST")
	router.HandleFunc("/api/transactions", walletHandler.GetTransactions).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

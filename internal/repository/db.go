package repository

import (
	"avito-shop-service/config"
	"database/sql"
	"fmt"
	"log"

	// Импортируем PostgreSQL
	_ "github.com/lib/pq"
)

func ConnectDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB is not reachable:", err)
	}

	log.Println("Connected to the database")
	return db
}

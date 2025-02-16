package middleware

import (
	"avito-shop-service/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func AuthMiddleware(authService *service.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем токен из заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Убираем префикс "Bearer "
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Парсим токен и извлекаем user_id
			userID, err := authService.ParseToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем user_id в заголовки запроса
			r.Header.Set("UserID", strconv.Itoa(userID))

			// Переходим к следующему обработчику
			next.ServeHTTP(w, r)
		})
	}
}

package handlers

import (
	"avito-shop-service/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Проверка на пустое имя пользователя и пароль
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password cannot be empty", http.StatusBadRequest)
		return
	}

	// Попытка авторизации
	token, err := h.authService.Login(req.Username, req.Password)
	if err == nil {
		// Если авторизация успешна, возвращаем токен
		if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		return
	}

	// Если ошибка (пользователя не сущетвует - первая аутентификация), регистрируем нового пользователя
	if err := h.authService.Register(req.Username, req.Password); err != nil {
		http.Error(w, "User registration failed", http.StatusConflict)
		return
	}

	// После успешной регистрации авторизуем пользователя и возвращаем токен
	token, err = h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Error during login after registration", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

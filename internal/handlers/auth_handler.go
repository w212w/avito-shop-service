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

	// Сначала пробуем залогиниться
	token, err := h.authService.Login(req.Username, req.Password)
	if err == nil {
		// Если аутентификация успешна, возвращаем токен
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}

	// Если ошибка (пользователя не сущетвует - первая аутентификация), регистрируем нового пользователя
	if err := h.authService.Register(req.Username, req.Password); err != nil {
		http.Error(w, "User registration failed", http.StatusConflict)
		return
	}

	// После успешной регистрации аутентифицируем пользователя и возвращаем токен
	token, err = h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Error during login after registration", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

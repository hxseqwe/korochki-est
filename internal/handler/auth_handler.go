package handler

import (
	"encoding/json"
	"github.com/hxseqwe/korochki-est/internal/model"
	"github.com/hxseqwe/korochki-est/internal/service"
	"net/http"
	"regexp"
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Phone = strings.ReplaceAll(req.Phone, "(", "")
	req.Phone = strings.ReplaceAll(req.Phone, ")", "")
	req.Phone = strings.ReplaceAll(req.Phone, "-", "")

	user, err := h.authService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authService.SetSession(w, r, user); err != nil {
		http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"login":    user.Login,
		"fullName": user.FullName,
		"isAdmin":  user.IsAdmin,
	})
}

func validateRegisterRequest(req *model.RegisterRequest) error {
	if len(req.Login) < 6 {
		return httpError("Логин должен содержать минимум 6 символов")
	}
	if matched, _ := regexp.MatchString("^[A-Za-z0-9]+$", req.Login); !matched {
		return httpError("Логин может содержать только латиницу и цифры")
	}

	if len(req.Password) < 8 {
		return httpError("Пароль должен содержать минимум 8 символов")
	}
	if matched, _ := regexp.MatchString("[A-Z]", req.Password); !matched {
		return httpError("Пароль должен содержать хотя бы одну заглавную букву")
	}
	if matched, _ := regexp.MatchString("[a-z]", req.Password); !matched {
		return httpError("Пароль должен содержать хотя бы одну строчную букву")
	}
	if matched, _ := regexp.MatchString("[0-9]", req.Password); !matched {
		return httpError("Пароль должен содержать хотя бы одну цифру")
	}
	if matched, _ := regexp.MatchString("[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>\\/?]", req.Password); !matched {
		return httpError("Пароль должен содержать хотя бы один специальный символ")
	}

	if req.FullName == "" {
		return httpError("ФИО обязательно для заполнения")
	}
	if matched, _ := regexp.MatchString("^[А-Яа-яЁё\\s]+$", req.FullName); !matched {
		return httpError("ФИО может содержать только кириллицу и пробелы")
	}

	if req.Phone == "" {
		return httpError("Телефон обязателен для заполнения")
	}
	cleanPhone := strings.ReplaceAll(req.Phone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")

	if matched, _ := regexp.MatchString("^8[0-9]{10}$", cleanPhone); !matched {
		return httpError("Телефон должен быть в формате 8(XXX)XXX-XX-XX")
	}

	if req.Email == "" {
		return httpError("Email обязателен для заполнения")
	}
	if matched, _ := regexp.MatchString("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$", req.Email); !matched {
		return httpError("Введите корректный email адрес")
	}

	return nil
}

func httpError(msg string) error {
	return &httpErr{msg: msg}
}

type httpErr struct {
	msg string
}

func (e *httpErr) Error() string {
	return e.msg
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(&req)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if err := h.authService.SetSession(w, r, user); err != nil {
		http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"login":    user.Login,
		"fullName": user.FullName,
		"isAdmin":  user.IsAdmin,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.authService.ClearSession(w, r)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := h.authService.GetCurrentUser(r)
		if err != nil {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *AuthHandler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.authService.IsAdmin(r) {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

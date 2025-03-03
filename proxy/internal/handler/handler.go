// @title DD API
// @version 1.0
// @description API сервиса DD
// @host localhost:8080
// @BasePath /api
// @schemes http https
// @contact.name API Support
// @contact.email support@example.com

package handler

import (
	"context"
	pb_auth "dd/pkg/auth"
	pb_geo "dd/pkg/geo"
	pb_user "dd/pkg/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"log"

	"google.golang.org/grpc/metadata"
)

// RegisterRequest Запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// LoginRequest Запрос на вход
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// AuthResponse Ответ с токеном
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// Структуры для запросов geo сервиса
// SearchAddressRequest Запрос на поиск адреса
type SearchAddressRequest struct {
	Query string `json:"query" example:"Москва, Тверская улица"`
}

// GeocodeRequest Запрос на геокодирование
type GeocodeRequest struct {
	Lat string `json:"lat" example:"55.7558"`
	Lng string `json:"lng" example:"37.6173"`
}

// Address Адрес
type Address struct {
	City   string `json:"city" example:"Москва"`
	Street string `json:"street" example:"Тверская"`
	House  string `json:"house" example:"1"`
	Lat    string `json:"lat" example:"55.7558"`
	Lon    string `json:"lon" example:"37.6173"`
}

// SearchAddressResponse Ответ на поиск адреса
type SearchAddressResponse struct {
	Addresses []*Address `json:"addresses"`
}

// GeocodeResponse Ответ на геокодирование
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

// Структуры для ответов user сервиса
// ProfileResponse Ответ с данными пользователя
type ProfileResponse struct {
	ID        string `json:"id" example:"123"`
	Email     string `json:"email" example:"user@example.com"`
	CreatedAt string `json:"created_at" example:"2024-03-03T12:00:00Z"`
}

// ListUsersResponse Ответ со списком пользователей
type ListUsersResponse struct {
	Users []ProfileResponse `json:"users"`
	Total int32             `json:"total" example:"42"`
}

type Handler struct {
	authClient pb_auth.AuthServiceClient
	geoClient  pb_geo.GeoServiceClient
	userClient pb_user.UserServiceClient
}

func New(authClient pb_auth.AuthServiceClient, geoClient pb_geo.GeoServiceClient, userClient pb_user.UserServiceClient) *Handler {
	return &Handler{
		authClient: authClient,
		geoClient:  geoClient,
		userClient: userClient,
	}
}

// Auth endpoints

// @Summary Register user
// @Description Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 200 {object} AuthResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в auth сервис
	resp, err := h.authClient.Register(context.Background(), &pb_auth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Registration failed: %v", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil || resp.Token == "" {
		log.Printf("Empty response from auth service")
		http.Error(w, "Registration failed: empty response", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(AuthResponse{
		Token: resp.Token,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// @Summary Login user
// @Description Аутентификация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} AuthResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в auth сервис
	resp, err := h.authClient.Login(context.Background(), &pb_auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: resp.Token,
	})
}

// Geo endpoints

// @Summary Search address
// @Description Поиск адреса по строке
// @Tags geo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body SearchAddressRequest true "Параметры поиска"
// @Success 200 {object} SearchAddressResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Router /address/search [post]
func (h *Handler) SearchAddress(w http.ResponseWriter, r *http.Request) {
	var req SearchAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"authorization", token,
	))

	resp, err := h.geoClient.SearchAddress(ctx, &pb_geo.SearchAddressRequest{
		Query: req.Query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Geocode coordinates
// @Description Получение адреса по координатам
// @Tags geo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body GeocodeRequest true "Координаты"
// @Success 200 {object} GeocodeResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Router /address/geocode [post]
func (h *Handler) Geocode(w http.ResponseWriter, r *http.Request) {
	var req GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"authorization", token,
	))

	resp, err := h.geoClient.Geocode(ctx, &pb_geo.GeocodeRequest{
		Address: fmt.Sprintf("%s,%s", req.Lat, req.Lng),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// User endpoints

// @Summary Get user profile
// @Description Получение профиля пользователя
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param email query string true "Email пользователя"
// @Success 200 {object} ProfileResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /user/profile [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Создаем контекст с метаданными
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"authorization", token,
	))

	resp, err := h.userClient.GetProfile(ctx, &pb_user.GetProfileRequest{
		Email: r.URL.Query().Get("email"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := ProfileResponse{
		ID:        resp.User.Id,
		Email:     resp.User.Email,
		CreatedAt: resp.User.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// @Summary List users
// @Description Получение списка пользователей
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query integer false "Номер страницы" default(1)
// @Param per_page query integer false "Количество записей на странице" default(10)
// @Success 200 {object} ListUsersResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /user/list [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем параметры пагинации
	page := 1
	perPage := 10
	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if val, err := strconv.Atoi(pp); err == nil && val > 0 {
			perPage = val
		}
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"authorization", token,
	))

	resp, err := h.userClient.ListUsers(ctx, &pb_user.ListUsersRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []ProfileResponse
	for _, u := range resp.Users {
		users = append(users, ProfileResponse{
			ID:        u.Id,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	response := ListUsersResponse{
		Users: users,
		Total: resp.Total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

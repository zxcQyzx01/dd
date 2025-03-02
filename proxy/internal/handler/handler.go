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

// Структуры для запросов
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Структура для ответа
type AuthResponse struct {
	Token string `json:"token"`
}

// Структуры для запросов geo сервиса
type SearchAddressRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// Структуры для ответов user сервиса
type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ListUsersResponse struct {
	Users []ProfileResponse `json:"users"`
	Total int32             `json:"total"`
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

// Register обрабатывает регистрацию нового пользователя
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

// Login обрабатывает вход пользователя
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

// Geo handlers
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

// User handlers
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

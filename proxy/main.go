package main

import (
	pb_auth "dd/pkg/auth"
	pb_geo "dd/pkg/geo"
	pb_user "dd/pkg/user"
	"dd/proxy/internal/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	_ "dd/docs" // импорт сгенерированной документации

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Подключаемся к микросервисам
	authConn, err := grpc.Dial("auth:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	geoConn, err := grpc.Dial("geo:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to geo service: %v", err)
	}
	defer geoConn.Close()

	userConn, err := grpc.Dial("user:50053", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	// Создаем клиентов для каждого сервиса
	authClient := pb_auth.NewAuthServiceClient(authConn)
	geoClient := pb_geo.NewGeoServiceClient(geoConn)
	userClient := pb_user.NewUserServiceClient(userConn)

	// Инициализируем handler
	h := handler.New(authClient, geoClient, userClient)

	// Настраиваем маршруты
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/auth/register", h.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", h.Login).Methods("POST")

	// Geo routes
	r.HandleFunc("/api/address/search", h.SearchAddress).Methods("POST")
	r.HandleFunc("/api/address/geocode", h.Geocode).Methods("POST")

	// User routes
	r.HandleFunc("/api/user/profile", h.GetProfile).Methods("GET")
	r.HandleFunc("/api/user/list", h.ListUsers).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Запускаем сервер
	log.Printf("Starting proxy server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

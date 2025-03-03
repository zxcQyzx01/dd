package main

import (
	"log"
	"net"
	"os"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"dd/geo/internal/config"
	"dd/geo/internal/service"
	pb "dd/pkg/geo"
)

func main() {
	cfg := config.New()

	// Получаем адрес auth сервиса из переменной окружения
	authHost := os.Getenv("AUTH_SERVICE")
	if authHost == "" {
		authHost = "auth1:50051" // значение по умолчанию
	}

	// Подключаемся к auth сервису
	authConn, err := grpc.Dial(authHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Получаем адрес Redis из переменной окружения
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	// Подключаемся к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	geoService := service.New(authConn, redisClient)
	pb.RegisterGeoServiceServer(grpcServer, geoService)

	log.Printf("Starting Geo service on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    "dd/auth/internal/config"
    "dd/auth/internal/service"
    pb "dd/pkg/auth"
)

func main() {
    cfg := config.New()

    // Подключаемся к сервису пользователей
    userConn, err := grpc.Dial(cfg.UserService, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to connect to user service: %v", err)
    }
    defer userConn.Close()

    

    // Создаем gRPC сервер
    lis, err := net.Listen("tcp", cfg.GRPCPort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    authService := service.New(userConn, cfg.JWTSecret)
    pb.RegisterAuthServiceServer(grpcServer, authService)

    log.Printf("Starting Auth service on port %s", cfg.GRPCPort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
} 
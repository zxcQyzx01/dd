package service

import (
	"context"
	pb "dd/pkg/auth"
	userpb "dd/pkg/user"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"

	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userClient userpb.UserServiceClient
	jwtSecret  string
}

func New(userConn *grpc.ClientConn, jwtSecret string) *AuthService {
	return &AuthService{
		userClient: userpb.NewUserServiceClient(userConn),
		jwtSecret:  jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Проверяем существование пользователя
	_, err := s.userClient.GetProfile(ctx, &userpb.GetProfileRequest{
		Email: req.Email,
	})
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	// Создаем пользователя
	userResp, err := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}

	// Генерируем токен
	token, err := s.generateToken(userResp.User.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.RegisterResponse{
		Token: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Проверяем пользователя через user service
	userResp, err := s.userClient.GetProfile(ctx, &userpb.GetProfileRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return nil, fmt.Errorf("authentication failed")
	}

	if userResp.User == nil {
		return nil, fmt.Errorf("user not found")
	}

	token, err := s.generateToken(userResp.User.Id)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	log.Printf("User %s logged in successfully", req.Email)

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims["user_id"].(string),
	}, nil
}

func (s *AuthService) generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"user_id": "user-id", // В реальном приложении здесь будет ID из базы
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

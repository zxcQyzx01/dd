package service

import (
	"context"
	"database/sql"
	pb "dd/pkg/user"
	"dd/user/internal/model"
	"fmt"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func New(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	var user model.User
	err = s.db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash) 
         VALUES ($1, $2) 
         RETURNING id, email, password_hash, created_at`,
		req.Email, string(hashedPassword)).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		},
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var user model.User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, created_at FROM users WHERE email = $1`,
		req.Email).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.GetProfileResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		},
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := (req.Page - 1) * req.PerPage

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, email, created_at FROM users 
         ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		req.PerPage, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan user: %v", err)
		}
		users = append(users, &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		})
	}

	var total int32
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count users: %v", err)
	}

	return &pb.ListUsersResponse{
		Users: users,
		Total: total,
	}, nil
}

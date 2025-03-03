package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb_auth "dd/pkg/auth"
	pb_user "dd/pkg/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Мок для клиента user сервиса
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) CreateUser(ctx context.Context, req *pb_user.CreateUserRequest, opts ...grpc.CallOption) (*pb_user.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_user.CreateUserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserClient) GetProfile(ctx context.Context, req *pb_user.GetProfileRequest, opts ...grpc.CallOption) (*pb_user.GetProfileResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_user.GetProfileResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserClient) ListUsers(ctx context.Context, req *pb_user.ListUsersRequest, opts ...grpc.CallOption) (*pb_user.ListUsersResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_user.ListUsersResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupTest(t *testing.T) (*AuthService, *MockUserClient) {
	mockUser := &MockUserClient{}
	service := &AuthService{
		userClient: mockUser,
		jwtSecret:  "test-secret",
	}
	return service, mockUser
}

func TestAuthService_Register(t *testing.T) {
	service, mockUser := setupTest(t)

	t.Run("successful registration", func(t *testing.T) {
		mockUser.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *pb_user.CreateUserRequest) bool {
			return req.Email == "test@example.com"
		})).Return(&pb_user.CreateUserResponse{
			User: &pb_user.User{
				Id:        "123",
				Email:     "test@example.com",
				CreatedAt: time.Now().String(),
			},
		}, nil)

		resp, err := service.Register(context.Background(), &pb_auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUser.On("CreateUser", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("failed to create user"))

		resp, err := service.Register(context.Background(), &pb_auth.RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to create user")
	})
}

func TestAuthService_Login(t *testing.T) {
	service, mockUser := setupTest(t)

	t.Run("successful login", func(t *testing.T) {
		mockUser.On("GetProfile", mock.Anything, mock.MatchedBy(func(req *pb_user.GetProfileRequest) bool {
			return req.Email == "test@example.com" && req.Password == "password123"
		})).Return(&pb_user.GetProfileResponse{
			User: &pb_user.User{
				Id:        "123",
				Email:     "test@example.com",
				CreatedAt: time.Now().String(),
			},
		}, nil)

		resp, err := service.Login(context.Background(), &pb_auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUser.On("GetProfile", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("authentication failed"))

		resp, err := service.Login(context.Background(), &pb_auth.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	service, _ := setupTest(t)

	t.Run("valid token", func(t *testing.T) {
		// Создаем валидный токен
		token, err := service.generateToken("test@example.com")
		assert.NoError(t, err)

		resp, err := service.ValidateToken(context.Background(), &pb_auth.ValidateTokenRequest{
			Token: token,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Valid)
		assert.Equal(t, "user-id", resp.UserId)
	})

	t.Run("invalid token", func(t *testing.T) {
		resp, err := service.ValidateToken(context.Background(), &pb_auth.ValidateTokenRequest{
			Token: "invalid-token",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Valid)
		assert.Empty(t, resp.UserId)
	})
}

func TestNew(t *testing.T) {
	conn, err := grpc.Dial("dummy", grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	service := New(conn, "test-secret")

	assert.NotNil(t, service)
	assert.NotNil(t, service.userClient)
	assert.Equal(t, "test-secret", service.jwtSecret)
}

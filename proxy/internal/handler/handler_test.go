package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb_auth "dd/pkg/auth"
	pb_geo "dd/pkg/geo"
	pb_user "dd/pkg/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Моки клиентов
type MockAuthClient struct{ mock.Mock }
type MockGeoClient struct{ mock.Mock }
type MockUserClient struct{ mock.Mock }

// Auth методы
func (m *MockAuthClient) Register(ctx context.Context, req *pb_auth.RegisterRequest, opts ...grpc.CallOption) (*pb_auth.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_auth.RegisterResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthClient) Login(ctx context.Context, req *pb_auth.LoginRequest, opts ...grpc.CallOption) (*pb_auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_auth.LoginResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthClient) ValidateToken(ctx context.Context, req *pb_auth.ValidateTokenRequest, opts ...grpc.CallOption) (*pb_auth.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_auth.ValidateTokenResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Geo методы
func (m *MockGeoClient) SearchAddress(ctx context.Context, req *pb_geo.SearchAddressRequest, opts ...grpc.CallOption) (*pb_geo.SearchAddressResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_geo.SearchAddressResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGeoClient) Geocode(ctx context.Context, req *pb_geo.GeocodeRequest, opts ...grpc.CallOption) (*pb_geo.GeocodeResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_geo.GeocodeResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// User методы
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

func (m *MockUserClient) CreateUser(ctx context.Context, req *pb_user.CreateUserRequest, opts ...grpc.CallOption) (*pb_user.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb_user.CreateUserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupTest() (*Handler, *MockAuthClient, *MockGeoClient, *MockUserClient) {
	mockAuth := &MockAuthClient{}
	mockGeo := &MockGeoClient{}
	mockUser := &MockUserClient{}
	h := New(mockAuth, mockGeo, mockUser)
	return h, mockAuth, mockGeo, mockUser
}

func TestHandler_Register(t *testing.T) {
	h, mockAuth, _, _ := setupTest()

	t.Run("successful registration", func(t *testing.T) {
		mockAuth.On("Register", mock.Anything, &pb_auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}).Return(&pb_auth.RegisterResponse{
			Token: "test-token",
		}, nil)

		body := bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"password": "password123"
		}`))

		req := httptest.NewRequest("POST", "/api/auth/register", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-token", response.Token)
	})

	t.Run("invalid request body", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`invalid json`))
		req := httptest.NewRequest("POST", "/api/auth/register", body)
		w := httptest.NewRecorder()

		h.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_Login(t *testing.T) {
	h, mockAuth, _, _ := setupTest()

	t.Run("successful login", func(t *testing.T) {
		mockAuth.On("Login", mock.Anything, &pb_auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}).Return(&pb_auth.LoginResponse{
			Token: "test-token",
		}, nil)

		body := bytes.NewBuffer([]byte(`{
			"email": "test@example.com",
			"password": "password123"
		}`))

		req := httptest.NewRequest("POST", "/api/auth/login", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-token", response.Token)
	})
}

func TestHandler_SearchAddress(t *testing.T) {
	h, _, mockGeo, _ := setupTest()

	t.Run("successful search", func(t *testing.T) {
		mockGeo.On("SearchAddress", mock.Anything, &pb_geo.SearchAddressRequest{
			Query: "test address",
		}).Return(&pb_geo.SearchAddressResponse{
			Addresses: []*pb_geo.Address{
				{
					City:   "Moscow",
					Street: "Test",
					House:  "1",
				},
			},
		}, nil)

		body := bytes.NewBuffer([]byte(`{
			"query": "test address"
		}`))

		req := httptest.NewRequest("POST", "/api/address/search", body)
		req.Header.Set("Authorization", "test-token")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.SearchAddress(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{
			"query": "test address"
		}`))
		req := httptest.NewRequest("POST", "/api/address/search", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.SearchAddress(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_GetProfile(t *testing.T) {
	h, _, _, mockUser := setupTest()

	t.Run("successful profile retrieval", func(t *testing.T) {
		mockUser.On("GetProfile", mock.Anything, &pb_user.GetProfileRequest{
			Email: "test@example.com",
		}).Return(&pb_user.GetProfileResponse{
			User: &pb_user.User{
				Id:        "123",
				Email:     "test@example.com",
				CreatedAt: "2024-01-01",
			},
		}, nil)

		req := httptest.NewRequest("GET", "/api/user/profile?email=test@example.com", nil)
		req.Header.Set("Authorization", "test-token")
		w := httptest.NewRecorder()

		h.GetProfile(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response ProfileResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", response.Email)
	})
}

func TestHandler_Geocode(t *testing.T) {
	h, _, mockGeo, _ := setupTest()

	t.Run("successful geocode", func(t *testing.T) {
		mockGeo.On("Geocode", mock.Anything, &pb_geo.GeocodeRequest{
			Address: "55.7558,37.6173",
		}).Return(&pb_geo.GeocodeResponse{
			Addresses: []*pb_geo.Address{
				{
					City:   "Moscow",
					Street: "Test",
					House:  "1",
					Lat:    "55.7558",
					Lon:    "37.6173",
				},
			},
		}, nil)

		body := bytes.NewBuffer([]byte(`{
			"lat": "55.7558",
			"lng": "37.6173"
		}`))

		req := httptest.NewRequest("POST", "/api/address/geocode", body)
		req.Header.Set("Authorization", "test-token")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Geocode(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_ListUsers(t *testing.T) {
	h, _, _, mockUser := setupTest()

	t.Run("successful listing", func(t *testing.T) {
		mockUser.On("ListUsers", mock.Anything, &pb_user.ListUsersRequest{
			Page:    1,
			PerPage: 10,
		}).Return(&pb_user.ListUsersResponse{
			Users: []*pb_user.User{
				{
					Id:        "1",
					Email:     "user1@example.com",
					CreatedAt: "2024-01-01",
				},
				{
					Id:        "2",
					Email:     "user2@example.com",
					CreatedAt: "2024-01-01",
				},
			},
			Total: 2,
		}, nil)

		req := httptest.NewRequest("GET", "/api/user/list?page=1&per_page=10", nil)
		req.Header.Set("Authorization", "test-token")
		w := httptest.NewRecorder()

		h.ListUsers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response ListUsersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Users, 2)
		assert.Equal(t, int32(2), response.Total)
	})
}

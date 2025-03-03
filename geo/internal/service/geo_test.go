package service

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"dd/geo/internal/dadata"
	"dd/geo/internal/domain"
	pb_auth "dd/pkg/auth"
	pb "dd/pkg/geo"
)

// MockAuthClient - мок для auth клиента
type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) Register(ctx context.Context, req *pb_auth.RegisterRequest, opts ...grpc.CallOption) (*pb_auth.RegisterResponse, error) {
	return nil, nil
}

func (m *MockAuthClient) Login(ctx context.Context, req *pb_auth.LoginRequest, opts ...grpc.CallOption) (*pb_auth.LoginResponse, error) {
	return nil, nil
}

func (m *MockAuthClient) ValidateToken(ctx context.Context, req *pb_auth.ValidateTokenRequest, opts ...grpc.CallOption) (*pb_auth.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb_auth.ValidateTokenResponse), args.Error(1)
}

// MockDaDataProvider - мок для DaData провайдера
type MockDaDataProvider struct {
	mock.Mock
}

func (m *MockDaDataProvider) AddressSearch(input string) ([]*domain.Address, error) {
	args := m.Called(input)
	return args.Get(0).([]*domain.Address), args.Error(1)
}

func (m *MockDaDataProvider) GeoCode(lat, lng string) ([]*domain.Address, error) {
	args := m.Called(lat, lng)
	return args.Get(0).([]*domain.Address), args.Error(1)
}

func setupTest(t *testing.T) (*GeoService, *MockAuthClient, *MockDaDataProvider, func()) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	mockAuth := &MockAuthClient{}
	mockDaData := &MockDaDataProvider{}

	service := &GeoService{
		authClient:  mockAuth,
		redisClient: redisClient,
		geoProvider: mockDaData,
	}

	cleanup := func() {
		mr.Close()
		redisClient.Close()
		mockAuth.ExpectedCalls = nil
		mockDaData.ExpectedCalls = nil
		mr.FlushAll()
	}

	return service, mockAuth, mockDaData, cleanup
}

func TestGeoService_SearchAddress(t *testing.T) {
	service, mockAuth, mockDaData, cleanup := setupTest(t)
	defer cleanup()

	t.Run("successful search", func(t *testing.T) {
		service.redisClient.FlushAll(context.Background())
		mockAuth.ExpectedCalls = nil
		mockDaData.ExpectedCalls = nil
		// Настраиваем мок auth клиента
		mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&pb_auth.ValidateTokenResponse{
			Valid: true,
		}, nil)

		// Настраиваем мок DaData
		mockAddresses := []*domain.Address{
			{
				City:   "Москва",
				Street: "Сухаревская",
				House:  "11",
				Lat:    "55.77412",
				Lon:    "37.624065",
			},
		}
		mockDaData.On("AddressSearch", "сухаревская 11").Return(mockAddresses, nil)

		// Выполняем запрос
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"authorization", "test-token",
		))
		resp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "сухаревская 11",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Addresses, 1)
		assert.Equal(t, "Москва", resp.Addresses[0].City)
		assert.Equal(t, "Сухаревская", resp.Addresses[0].Street)
		assert.Equal(t, "11", resp.Addresses[0].House)
	})

	t.Run("invalid token", func(t *testing.T) {
		service.redisClient.FlushAll(context.Background())
		mockAuth.ExpectedCalls = nil
		mockDaData.ExpectedCalls = nil
		mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&pb_auth.ValidateTokenResponse{
			Valid: false,
		}, nil)

		ctx := context.Background() // Используем пустой контекст для теста невалидного токена

		resp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "test",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		// Проверяем, что AddressSearch не был вызван
		mockDaData.AssertNotCalled(t, "AddressSearch")
	})
}

func TestGeoService_Geocode(t *testing.T) {
	service, mockAuth, mockDaData, cleanup := setupTest(t)
	defer cleanup()

	t.Run("successful geocode", func(t *testing.T) {
		mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&pb_auth.ValidateTokenResponse{
			Valid: true,
		}, nil)

		mockAddresses := []*domain.Address{
			{
				City:   "Москва",
				Street: "Сухаревская",
				House:  "11",
				Lat:    "55.77412",
				Lon:    "37.624065",
			},
		}
		mockDaData.On("GeoCode", "55.77412", "37.624065").Return(mockAddresses, nil)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"authorization", "test-token",
		))
		resp, err := service.Geocode(ctx, &pb.GeocodeRequest{
			Address: "55.77412,37.624065",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Addresses, 1)
		assert.Equal(t, "Москва", resp.Addresses[0].City)
	})
}

func TestGeoService_SearchAddress_CachedResponse(t *testing.T) {
	service, mockAuth, mockDaData, cleanup := setupTest(t)
	defer cleanup()

	t.Run("cached response", func(t *testing.T) {
		service.redisClient.FlushAll(context.Background())
		mockAuth.ExpectedCalls = nil
		mockDaData.ExpectedCalls = nil

		mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&pb_auth.ValidateTokenResponse{
			Valid: true,
		}, nil).Times(2)

		mockAddresses := []*domain.Address{
			{
				City:   "Москва",
				Street: "Тестовая",
				House:  "1",
				Lat:    "55.7558",
				Lon:    "37.6173",
			},
		}
		mockDaData.On("AddressSearch", "test cache").Return(mockAddresses, nil).Once()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"authorization", "test-token",
		))

		// Первый запрос (сохранение в кэш)
		firstResp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "test cache",
		})
		assert.NoError(t, err)

		// Проверяем что значение сохранилось в кэш
		cacheKey := "search:test cache"
		_, err = service.redisClient.Get(ctx, cacheKey).Result()
		assert.NoError(t, err, "Значение должно быть в кэше")

		// Второй запрос (должен взяться из кэша)
		secondResp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "test cache",
		})

		assert.NoError(t, err)
		assert.NotNil(t, secondResp)
		assert.Equal(t, firstResp, secondResp)
		mockDaData.AssertNumberOfCalls(t, "AddressSearch", 1)
	})

	t.Run("cached response", func(t *testing.T) {
		service.redisClient.FlushAll(context.Background())
		mockAuth.ExpectedCalls = nil
		mockDaData.ExpectedCalls = nil

		mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&pb_auth.ValidateTokenResponse{
			Valid: true,
		}, nil).Times(2)

		mockAddresses := []*domain.Address{
			{
				City:   "Москва",
				Street: "Тестовая",
				House:  "1",
				Lat:    "55.7558",
				Lon:    "37.6173",
			},
		}
		mockDaData.On("AddressSearch", "test cache").Return(mockAddresses, nil).Once()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"authorization", "test-token",
		))

		// Первый запрос (сохранение в кэш)
		firstResp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "test cache",
		})
		assert.NoError(t, err)

		// Проверяем что значение сохранилось в кэш
		cacheKey := "search:test cache"
		_, err = service.redisClient.Get(ctx, cacheKey).Result()
		assert.NoError(t, err, "Значение должно быть в кэше")

		// Второй запрос (должен взяться из кэша)
		secondResp, err := service.SearchAddress(ctx, &pb.SearchAddressRequest{
			Query: "test cache",
		})

		assert.NoError(t, err)
		assert.NotNil(t, secondResp)
		assert.Equal(t, firstResp, secondResp)
		mockDaData.AssertExpectations(t)
	})
}

func TestNew(t *testing.T) {
	// Создаем тестовый Redis клиент
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	// Создаем тестовое gRPC подключение
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	// Создаем сервис
	service := New(conn, redisClient)

	// Проверяем что сервис создан корректно
	assert.NotNil(t, service)
	assert.NotNil(t, service.authClient)
	assert.NotNil(t, service.redisClient)
	assert.NotNil(t, service.geoProvider)

	// Проверяем что все зависимости установлены
	assert.Equal(t, redisClient, service.redisClient)
	assert.IsType(t, &dadata.Provider{}, service.geoProvider)
}

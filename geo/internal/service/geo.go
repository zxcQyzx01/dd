package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"dd/geo/internal/dadata"
	"dd/geo/internal/domain"
	pb_auth "dd/pkg/auth"
	pb "dd/pkg/geo"
)

const (
	searchCachePrefix = "search:"
	geoCachePrefix    = "geo:"
	cacheDuration     = 24 * time.Hour
)

type GeoService struct {
	pb.UnimplementedGeoServiceServer
	authClient  pb_auth.AuthServiceClient
	redisClient *redis.Client
	geoProvider domain.GeoProvider
}

func New(authConn *grpc.ClientConn, redisClient *redis.Client) *GeoService {
	provider := dadata.NewProvider(
		"627de73a10855ebb80eb0191f2bbb55cc72eef89",
		"7886bc85cac2562af90304564e7f04078d18dc4b",
	)

	return &GeoService{
		authClient:  pb_auth.NewAuthServiceClient(authConn),
		redisClient: redisClient,
		geoProvider: provider,
	}
}

func (s *GeoService) validateToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "no metadata provided")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return status.Error(codes.Unauthenticated, "no token provided")
	}

	resp, err := s.authClient.ValidateToken(ctx, &pb_auth.ValidateTokenRequest{
		Token: tokens[0],
	})
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	if !resp.Valid {
		return status.Error(codes.Unauthenticated, "token is not valid")
	}

	return nil
}

func (s *GeoService) SearchAddress(ctx context.Context, req *pb.SearchAddressRequest) (*pb.SearchAddressResponse, error) {
	if err := s.validateToken(ctx); err != nil {
		return nil, err
	}

	// Проверяем кэш
	cacheKey := searchCachePrefix + req.Query
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var response pb.SearchAddressResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	addresses, err := s.geoProvider.AddressSearch(req.Query)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to search address: %v", err))
	}

	var pbAddresses []*pb.Address
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		})
	}

	response := &pb.SearchAddressResponse{
		Addresses: pbAddresses,
	}

	// Сохраняем в кэш
	if cached, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, string(cached), cacheDuration)
	}

	return response, nil
}

func (s *GeoService) Geocode(ctx context.Context, req *pb.GeocodeRequest) (*pb.GeocodeResponse, error) {
	if err := s.validateToken(ctx); err != nil {
		return nil, err
	}

	parts := strings.Split(req.Address, ",")
	if len(parts) != 2 {
		return nil, status.Error(codes.InvalidArgument, "invalid coordinates format")
	}

	lat := strings.TrimSpace(parts[0])
	lon := strings.TrimSpace(parts[1])

	// Проверяем кэш
	cacheKey := fmt.Sprintf("%s%s,%s", geoCachePrefix, lat, lon)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var response pb.GeocodeResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	addresses, err := s.geoProvider.GeoCode(lat, lon)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to geocode: %v", err))
	}

	if len(addresses) == 0 {
		return nil, status.Error(codes.NotFound, "address not found")
	}

	var pbAddresses []*pb.Address
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		})
	}

	response := &pb.GeocodeResponse{
		Addresses: pbAddresses,
	}

	// Сохраняем в кэш
	if cached, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, string(cached), cacheDuration)
	}

	return response, nil
}

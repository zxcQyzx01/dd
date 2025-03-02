package dadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"dd/geo/internal/domain"
	"dd/geo/internal/models"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

// Provider реализует взаимодействие с API DaData
type Provider struct {
	api       *suggest.Api // Клиент API DaData
	apiKey    string       // Ключ API
	secretKey string       // Секретный ключ
}

// NewProvider создает новый экземпляр провайдера DaData
func NewProvider(apiKey, secretKey string) *Provider {
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	return &Provider{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

// AddressSearch выполняет поиск адресов через API DaData
func (p *Provider) AddressSearch(input string) ([]*domain.Address, error) {
	var res []*domain.Address
	rawRes, err := p.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &domain.Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}
	return res, nil
}

// GeoCode выполняет геокодирование через API DaData
func (p *Provider) GeoCode(lat, lng string) ([]*domain.Address, error) {
	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", p.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var geoCode models.GeoCode
	if err := json.NewDecoder(resp.Body).Decode(&geoCode); err != nil {
		return nil, err
	}

	var res []*domain.Address
	for _, r := range geoCode.Suggestions {
		res = append(res, &domain.Address{
			City:   string(r.Data.City),
			Street: string(r.Data.Street),
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}
	return res, nil
}

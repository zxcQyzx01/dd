package domain

// Address представляет структуру адреса с основными полями
type Address struct {
	City   string `json:"city"`   // Город
	Street string `json:"street"` // Улица
	House  string `json:"house"`  // Номер дома
	Lat    string `json:"lat"`    // Широта
	Lon    string `json:"lon"`    // Долгота
}

// GeoProvider определяет интерфейс для работы с геоданными
type GeoProvider interface {
	// AddressSearch ищет адреса по строке запроса
	AddressSearch(input string) ([]*Address, error)
	// GeoCode получает адрес по координатам
	GeoCode(lat, lng string) ([]*Address, error)
}

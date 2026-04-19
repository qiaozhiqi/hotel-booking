package suppliers

import (
	"hotel-booking/models"
	"sync"
)

type SupplierHotelData struct {
	SupplierHotelID string
	Name            string
	City            string
	Address         string
	Description     string
	Rating          float64
	ImageURL        string
	PriceRange      string
	Rooms           []SupplierRoomData
}

type SupplierRoomData struct {
	SupplierRoomID  string
	Name            string
	Description     string
	Price           float64
	Capacity        int
	Area            int
	BedType         string
	Amenities       string
	ImageURL        string
	TotalCount      int
	AvailableCount  int
}

type SupplierAdapter interface {
	GetCode() string
	GetName() string
	GetDescription() string
	GetAPIURL() string
	
	FetchHotels() ([]SupplierHotelData, error)
	FetchHotelDetail(hotelID string) (*SupplierHotelData, error)
}

var (
	adapters = make(map[string]SupplierAdapter)
	mutex    sync.RWMutex
)

func RegisterAdapter(adapter SupplierAdapter) {
	mutex.Lock()
	defer mutex.Unlock()
	adapters[adapter.GetCode()] = adapter
}

func GetAdapter(code string) SupplierAdapter {
	mutex.RLock()
	defer mutex.RUnlock()
	return adapters[code]
}

func GetAllAdapters() []SupplierAdapter {
	mutex.RLock()
	defer mutex.RUnlock()
	result := make([]SupplierAdapter, 0, len(adapters))
	for _, adapter := range adapters {
		result = append(result, adapter)
	}
	return result
}

func ToSupplierModel(adapter SupplierAdapter) models.Supplier {
	return models.Supplier{
		Code:        adapter.GetCode(),
		Name:        adapter.GetName(),
		Description: adapter.GetDescription(),
		APIURL:      adapter.GetAPIURL(),
		Status:      "active",
	}
}

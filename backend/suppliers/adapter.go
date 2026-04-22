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

type SupplierHotelStaticData struct {
	SupplierHotelID string
	Name            string
	City            string
	Address         string
	Description     string
	Rating          float64
	ImageURL        string
	PriceRange      string
	RoomTypes       []SupplierRoomTypeData
}

type SupplierRoomTypeData struct {
	SupplierRoomID string
	Name           string
	Description    string
	Capacity       int
	Area           int
	BedType        string
	Amenities      string
	ImageURL       string
	TotalCount     int
}

type SupplierPriceInventoryData struct {
	SupplierHotelID string
	SupplierRoomID  string
	Date            string
	Price           float64
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

type PriceInventoryFetcher interface {
	FetchPriceInventory(hotelID string, checkInDate, checkOutDate string) ([]SupplierPriceInventoryData, error)
}

type HotelStaticFetcher interface {
	FetchHotelStaticList() ([]SupplierHotelStaticData, error)
	FetchHotelStaticDetail(hotelID string) (*SupplierHotelStaticData, error)
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

func ConvertToStaticData(hotel SupplierHotelData) SupplierHotelStaticData {
	roomTypes := make([]SupplierRoomTypeData, len(hotel.Rooms))
	for i, room := range hotel.Rooms {
		roomTypes[i] = SupplierRoomTypeData{
			SupplierRoomID: room.SupplierRoomID,
			Name:           room.Name,
			Description:    room.Description,
			Capacity:       room.Capacity,
			Area:           room.Area,
			BedType:        room.BedType,
			Amenities:      room.Amenities,
			ImageURL:       room.ImageURL,
			TotalCount:     room.TotalCount,
		}
	}
	
	return SupplierHotelStaticData{
		SupplierHotelID: hotel.SupplierHotelID,
		Name:            hotel.Name,
		City:            hotel.City,
		Address:         hotel.Address,
		Description:     hotel.Description,
		Rating:          hotel.Rating,
		ImageURL:        hotel.ImageURL,
		PriceRange:      hotel.PriceRange,
		RoomTypes:       roomTypes,
	}
}

func ConvertFromStaticData(hotel SupplierHotelStaticData, priceInventory []SupplierPriceInventoryData) SupplierHotelData {
	rooms := make([]SupplierRoomData, len(hotel.RoomTypes))
	
	priceMap := make(map[string]SupplierPriceInventoryData)
	for _, pi := range priceInventory {
		priceMap[pi.SupplierRoomID] = pi
	}
	
	for i, roomType := range hotel.RoomTypes {
		price := 0.0
		available := 0
		
		if pi, ok := priceMap[roomType.SupplierRoomID]; ok {
			price = pi.Price
			available = pi.AvailableCount
		}
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: roomType.SupplierRoomID,
			Name:           roomType.Name,
			Description:    roomType.Description,
			Price:          price,
			Capacity:       roomType.Capacity,
			Area:           roomType.Area,
			BedType:        roomType.BedType,
			Amenities:      roomType.Amenities,
			ImageURL:       roomType.ImageURL,
			TotalCount:     roomType.TotalCount,
			AvailableCount: available,
		}
	}
	
	return SupplierHotelData{
		SupplierHotelID: hotel.SupplierHotelID,
		Name:            hotel.Name,
		City:            hotel.City,
		Address:         hotel.Address,
		Description:     hotel.Description,
		Rating:          hotel.Rating,
		ImageURL:        hotel.ImageURL,
		PriceRange:      hotel.PriceRange,
		Rooms:           rooms,
	}
}

package models

import "time"

type Supplier struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	APIURL       string    `json:"api_url"`
	APIKey       string    `json:"-"`
	Status       string    `json:"status"`
	Priority     int       `json:"priority"`
	PriceControl float64   `json:"price_control"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SupplierHotel struct {
	ID            int       `json:"id"`
	SupplierID    int       `json:"supplier_id"`
	SupplierHotelID string  `json:"supplier_hotel_id"`
	LocalHotelID  int       `json:"local_hotel_id"`
	HotelName     string    `json:"hotel_name"`
	City          string    `json:"city"`
	Address       string    `json:"address"`
	Rating        float64   `json:"rating"`
	ImageURL      string    `json:"image_url"`
	PriceRange    string    `json:"price_range"`
	RawData       string    `json:"-"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SupplierRoom struct {
	ID            int       `json:"id"`
	SupplierID    int       `json:"supplier_id"`
	SupplierHotelID string  `json:"supplier_hotel_id"`
	SupplierRoomID string   `json:"supplier_room_id"`
	LocalRoomID   int       `json:"local_room_id"`
	RoomName      string    `json:"room_name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Capacity      int       `json:"capacity"`
	Area          int       `json:"area"`
	BedType       string    `json:"bed_type"`
	Amenities     string    `json:"amenities"`
	ImageURL      string    `json:"image_url"`
	TotalCount    int       `json:"total_count"`
	AvailableCount int      `json:"available_count"`
	RawData       string    `json:"-"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Hotel struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	ImageURL    string    `json:"image_url"`
	PriceRange  string    `json:"price_range"`
	SupplierID  int       `json:"supplier_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	ID             int       `json:"id"`
	HotelID        int       `json:"hotel_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Capacity       int       `json:"capacity"`
	Area           int       `json:"area"`
	BedType        string    `json:"bed_type"`
	Amenities      string    `json:"amenities"`
	ImageURL       string    `json:"image_url"`
	TotalCount     int       `json:"total_count"`
	AvailableCount int       `json:"available_count"`
	SupplierID     int       `json:"supplier_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Order struct {
	ID          int       `json:"id"`
	OrderNo     string    `json:"order_no"`
	UserID      int       `json:"user_id"`
	HotelID     int       `json:"hotel_id"`
	RoomID      int       `json:"room_id"`
	CheckIn     string    `json:"check_in"`
	CheckOut    string    `json:"check_out"`
	GuestName   string    `json:"guest_name"`
	GuestPhone  string    `json:"guest_phone"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	HotelName   string    `json:"hotel_name"`
	RoomName    string    `json:"room_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CreateOrderRequest struct {
	HotelID    int    `json:"hotel_id" binding:"required"`
	RoomID     int    `json:"room_id" binding:"required"`
	CheckIn    string `json:"check_in" binding:"required"`
	CheckOut   string `json:"check_out" binding:"required"`
	GuestName  string `json:"guest_name" binding:"required"`
	GuestPhone string `json:"guest_phone" binding:"required"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SupplierInfo struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	Priority     int     `json:"priority"`
	PriceControl float64 `json:"price_control"`
}

type HotelWithSupplier struct {
	Hotel
	Supplier *SupplierInfo `json:"supplier,omitempty"`
}

type RoomWithSupplier struct {
	Room
	Supplier *SupplierInfo `json:"supplier,omitempty"`
}

type ChannelPrice struct {
	SupplierID     int     `json:"supplier_id"`
	SupplierName   string  `json:"supplier_name"`
	SupplierCode   string  `json:"supplier_code"`
	Price          float64 `json:"price"`
	OriginalPrice  float64 `json:"original_price"`
	AvailableCount int     `json:"available_count"`
	Priority       int     `json:"priority"`
	IsBestPrice    bool    `json:"is_best_price"`
}

type RoomWithChannelPrices struct {
	ID             int            `json:"id"`
	HotelID        int            `json:"hotel_id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Price          float64        `json:"price"`
	Capacity       int            `json:"capacity"`
	Area           int            `json:"area"`
	BedType        string         `json:"bed_type"`
	Amenities      string         `json:"amenities"`
	ImageURL       string         `json:"image_url"`
	TotalCount     int            `json:"total_count"`
	AvailableCount int            `json:"available_count"`
	BestPrice      float64        `json:"best_price"`
	ChannelPrices  []ChannelPrice `json:"channel_prices"`
}

type PriceInventory struct {
	ID             int       `json:"id"`
	SupplierID     int       `json:"supplier_id"`
	SupplierHotelID string   `json:"supplier_hotel_id"`
	SupplierRoomID string    `json:"supplier_room_id"`
	Date           string    `json:"date"`
	Price          float64   `json:"price"`
	OriginalPrice  float64   `json:"original_price"`
	AvailableCount int       `json:"available_count"`
	TotalCount     int       `json:"total_count"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type QiuguoPushHotelData struct {
	SupplierHotelID string                    `json:"hotel_id" binding:"required"`
	Name            string                    `json:"name" binding:"required"`
	City            string                    `json:"city"`
	Address         string                    `json:"address"`
	Description     string                    `json:"description"`
	Rating          float64                   `json:"rating"`
	ImageURL        string                    `json:"image_url"`
	PriceRange      string                    `json:"price_range"`
	Rooms           []QiuguoPushRoomData      `json:"rooms"`
}

type QiuguoPushRoomData struct {
	SupplierRoomID string `json:"room_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Capacity       int    `json:"capacity"`
	Area           int    `json:"area"`
	BedType        string `json:"bed_type"`
	Amenities      string `json:"amenities"`
	ImageURL       string `json:"image_url"`
	TotalCount     int    `json:"total_count"`
}

type QiuguoPushPriceInventoryData struct {
	SupplierHotelID string  `json:"hotel_id" binding:"required"`
	SupplierRoomID  string  `json:"room_id" binding:"required"`
	Date            string  `json:"date" binding:"required"`
	Price           float64 `json:"price" binding:"required"`
	OriginalPrice   float64 `json:"original_price"`
	AvailableCount  int     `json:"available_count" binding:"required"`
	TotalCount      int     `json:"total_count"`
}

type QiuguoPushRequest struct {
	RequestID       string                       `json:"request_id" binding:"required"`
	PushType        string                       `json:"push_type"`
	Hotels          []QiuguoPushHotelData        `json:"hotels"`
	PriceInventories []QiuguoPushPriceInventoryData `json:"price_inventories"`
	StartDate       string                       `json:"start_date"`
	EndDate         string                       `json:"end_date"`
}

type QiuguoPushResponse struct {
	RequestID    string `json:"request_id"`
	Success      bool   `json:"success"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	ProcessedCount int  `json:"processed_count,omitempty"`
	FailedCount  int    `json:"failed_count,omitempty"`
	ErrorDetails []string `json:"error_details,omitempty"`
}

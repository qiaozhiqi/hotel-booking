package models

import "time"

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

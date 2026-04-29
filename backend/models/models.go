package models

import (
	"fmt"
	"time"
)

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

type Guest struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	IDType      string    `json:"id_type"`
	IDNumber    string    `json:"id_number"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	ID                 int                 `json:"id"`
	OrderNo            string              `json:"order_no"`
	UserID             int                 `json:"user_id"`
	HotelID            int                 `json:"hotel_id"`
	RoomID             int                 `json:"room_id"`
	CheckIn            string              `json:"check_in"`
	CheckOut           string              `json:"check_out"`
	GuestName          string              `json:"guest_name"`
	GuestPhone         string              `json:"guest_phone"`
	TotalAmount        float64             `json:"total_amount"`
	Status             string              `json:"status"`
	HotelName          string              `json:"hotel_name"`
	RoomName           string              `json:"room_name"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	CancellationPolicy *CancellationPolicy `json:"cancellation_policy"`
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
	GuestName  string `json:"guest_name"`
	GuestPhone string `json:"guest_phone"`
	GuestID    int    `json:"guest_id"`
}

type CreateGuestRequest struct {
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	IDType    string `json:"id_type"`
	IDNumber  string `json:"id_number"`
	IsDefault bool   `json:"is_default"`
}

type UpdateGuestRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	IDType    string `json:"id_type"`
	IDNumber  string `json:"id_number"`
	IsDefault *bool  `json:"is_default"`
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
	SupplierID         int                   `json:"supplier_id"`
	SupplierName       string                `json:"supplier_name"`
	SupplierCode       string                `json:"supplier_code"`
	Price              float64               `json:"price"`
	OriginalPrice      float64               `json:"original_price"`
	AvailableCount     int                   `json:"available_count"`
	Priority           int                   `json:"priority"`
	IsBestPrice        bool                  `json:"is_best_price"`
	CancellationPolicy *CancellationPolicy   `json:"cancellation_policy,omitempty"`
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

type RoomPriceSummary struct {
	ID             int       `json:"id"`
	SupplierID     int       `json:"supplier_id"`
	SupplierHotelID string   `json:"supplier_hotel_id"`
	SupplierRoomID string    `json:"supplier_room_id"`
	MinPrice       float64   `json:"min_price"`
	MaxPrice       float64   `json:"max_price"`
	AvgPrice       float64   `json:"avg_price"`
	PriceRange     string    `json:"price_range"`
	HasInventory   bool      `json:"has_inventory"`
	TotalCount     int       `json:"total_count"`
	DateRangeStart string    `json:"date_range_start"`
	DateRangeEnd   string    `json:"date_range_end"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type HotelPriceSummary struct {
	ID             int       `json:"id"`
	SupplierID     int       `json:"supplier_id"`
	SupplierHotelID string   `json:"supplier_hotel_id"`
	MinPrice       float64   `json:"min_price"`
	MaxPrice       float64   `json:"max_price"`
	AvgPrice       float64   `json:"avg_price"`
	PriceRange     string    `json:"price_range"`
	TotalRooms     int       `json:"total_rooms"`
	RoomsWithInventory int   `json:"rooms_with_inventory"`
	DateRangeStart string    `json:"date_range_start"`
	DateRangeEnd   string    `json:"date_range_end"`
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

type CancellationPolicyType string

const (
	CancellationPolicyFree       CancellationPolicyType = "free"
	CancellationPolicyNonFree    CancellationPolicyType = "non_free"
	CancellationPolicyNonCancellable CancellationPolicyType = "non_cancellable"
)

type CancellationPolicy struct {
	Type             CancellationPolicyType `json:"type"`
	TypeName         string                 `json:"type_name"`
	FreeCancelBefore int                    `json:"free_cancel_before,omitempty"`
	FreeCancelUnit   string                 `json:"free_cancel_unit,omitempty"`
	CancelFeePercent float64                `json:"cancel_fee_percent,omitempty"`
	CancelFeeAmount  float64                `json:"cancel_fee_amount,omitempty"`
	Description      string                 `json:"description"`
}

func (t CancellationPolicyType) Name() string {
	switch t {
	case CancellationPolicyFree:
		return "免费取消"
	case CancellationPolicyNonFree:
		return "不可免费取消"
	case CancellationPolicyNonCancellable:
		return "不可取消"
	default:
		return "未知"
	}
}

func (t CancellationPolicyType) IsValid() bool {
	switch t {
	case CancellationPolicyFree, CancellationPolicyNonFree, CancellationPolicyNonCancellable:
		return true
	default:
		return false
	}
}

func GenerateCancellationPolicy(supplierCode string, hotelIndex int) CancellationPolicy {
	var policy CancellationPolicy
	
	policyTypes := []CancellationPolicyType{
		CancellationPolicyFree,
		CancellationPolicyFree,
		CancellationPolicyNonFree,
		CancellationPolicyNonCancellable,
	}
	
	policyIndex := (len(supplierCode) + hotelIndex) % len(policyTypes)
	policy.Type = policyTypes[policyIndex]
	policy.TypeName = policy.Type.Name()
	
	switch policy.Type {
	case CancellationPolicyFree:
		freeBeforeOptions := []struct {
			before int
			unit   string
		}{
			{24, "小时"},
			{48, "小时"},
			{72, "小时"},
			{24, "小时"},
		}
		optIndex := (hotelIndex + len(supplierCode)) % len(freeBeforeOptions)
		policy.FreeCancelBefore = freeBeforeOptions[optIndex].before
		policy.FreeCancelUnit = freeBeforeOptions[optIndex].unit
		policy.Description = fmt.Sprintf("入住前%d%s可免费取消", policy.FreeCancelBefore, policy.FreeCancelUnit)
		
	case CancellationPolicyNonFree:
		feeOptions := []struct {
			percent float64
			before  int
			unit    string
		}{
			{20.0, 24, "小时"},
			{30.0, 48, "小时"},
			{15.0, 12, "小时"},
		}
		optIndex := (hotelIndex * 2) % len(feeOptions)
		opt := feeOptions[optIndex]
		policy.FreeCancelBefore = opt.before
		policy.FreeCancelUnit = opt.unit
		policy.CancelFeePercent = opt.percent
		policy.Description = fmt.Sprintf("入住前%d%s内取消，收取订单金额%.0f%%的违约金", opt.before, opt.unit, opt.percent)
		
	case CancellationPolicyNonCancellable:
		policy.Description = "订单一经确认，不可取消和修改"
	}
	
	return policy
}

func GenerateCancellationPolicyBySupplier(supplierCode string) CancellationPolicy {
	var policy CancellationPolicy
	
	switch supplierCode {
	case "huazhu":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 24
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前24小时可免费取消，24小时内取消收取首晚房费的20%"
		
	case "jinjiang":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 48
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前48小时可免费取消，48小时内取消收取订单金额的15%"
		
	case "rujia":
		policy.Type = CancellationPolicyNonFree
		policy.TypeName = "不可免费取消"
		policy.FreeCancelBefore = 24
		policy.FreeCancelUnit = "小时"
		policy.CancelFeePercent = 20.0
		policy.Description = "入住前24小时内取消，收取订单金额20%的违约金"
		
	case "shiji_marriott":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 72
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前72小时可免费取消，72小时内取消收取首晚房费"
		
	case "shiji_hilton":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 48
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前48小时可免费取消，48小时内取消收取订单金额的25%"
		
	case "shiji_ihg":
		policy.Type = CancellationPolicyNonFree
		policy.TypeName = "不可免费取消"
		policy.FreeCancelBefore = 24
		policy.FreeCancelUnit = "小时"
		policy.CancelFeePercent = 30.0
		policy.Description = "入住前24小时内取消，收取订单金额30%的违约金"
		
	case "shiji_kaiyuan":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 24
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前24小时可免费取消，24小时内取消收取订单金额的10%"
		
	case "shiji_wanda":
		policy.Type = CancellationPolicyNonCancellable
		policy.TypeName = "不可取消"
		policy.Description = "订单一经确认，不可取消和修改，如需取消请联系客服"
		
	case "shiji_lvdi":
		policy.Type = CancellationPolicyFree
		policy.TypeName = "免费取消"
		policy.FreeCancelBefore = 12
		policy.FreeCancelUnit = "小时"
		policy.Description = "入住前12小时可免费取消，12小时内取消收取首晚房费的15%"
		
	case "shiji_qiuguo":
		policy.Type = CancellationPolicyNonFree
		policy.TypeName = "不可免费取消"
		policy.FreeCancelBefore = 6
		policy.FreeCancelUnit = "小时"
		policy.CancelFeePercent = 25.0
		policy.Description = "入住前6小时内取消，收取订单金额25%的违约金"
		
	default:
		policy = GenerateCancellationPolicy(supplierCode, 0)
	}
	
	return policy
}

type Favorite struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	HotelID    int       `json:"hotel_id"`
	HotelName  string    `json:"hotel_name"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
	Rating     float64   `json:"rating"`
	ImageURL   string    `json:"image_url"`
	PriceRange string    `json:"price_range"`
	CreatedAt  time.Time `json:"created_at"`
}

type Invoice struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	OrderID      int       `json:"order_id"`
	OrderNo      string    `json:"order_no"`
	InvoiceType  string    `json:"invoice_type"`
	InvoiceTitle string    `json:"invoice_title"`
	TaxNumber    string    `json:"tax_number"`
	BankName     string    `json:"bank_name"`
	BankAccount  string    `json:"bank_account"`
	Address      string    `json:"address"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	InvoiceNo    string    `json:"invoice_no"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateFavoriteRequest struct {
	HotelID int `json:"hotel_id" binding:"required"`
}

type CreateInvoiceRequest struct {
	OrderID      int    `json:"order_id" binding:"required"`
	InvoiceType  string `json:"invoice_type" binding:"required"`
	InvoiceTitle string `json:"invoice_title" binding:"required"`
	TaxNumber    string `json:"tax_number"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

type FavoriteWithHotel struct {
	Favorite
	Hotel *Hotel `json:"hotel,omitempty"`
}

type InvoiceWithOrder struct {
	Invoice
	Order *Order `json:"order,omitempty"`
}

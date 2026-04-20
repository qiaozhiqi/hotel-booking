package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHotelList(c *gin.Context) {
	city := c.Query("city")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	supplierCode := c.Query("supplier_code")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	db := database.GetDB()

	var rows *sql.Rows
	var err error
	var countRows *sql.Rows

	baseQuery := `
		SELECT id, name, address, city, description, rating, image_url, 
		       price_range, min_price, max_price, supplier_id, supplier_code, supplier_name, brand
		FROM hotels 
		WHERE 1=1`

	countQuery := "SELECT COUNT(*) FROM hotels WHERE 1=1"

	var args []interface{}
	var countArgs []interface{}

	if city != "" {
		baseQuery += " AND city = ?"
		countQuery += " AND city = ?"
		args = append(args, city)
		countArgs = append(countArgs, city)
	}

	if supplierCode != "" {
		baseQuery += " AND supplier_code = ?"
		countQuery += " AND supplier_code = ?"
		args = append(args, supplierCode)
		countArgs = append(countArgs, supplierCode)
	}

	baseQuery += " ORDER BY rating DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err = db.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取酒店列表失败",
		})
		return
	}
	defer rows.Close()

	countRows, _ = db.Query(countQuery, countArgs...)

	var hotels []models.Hotel
	for rows.Next() {
		var hotel models.Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
			&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange,
			&hotel.MinPrice, &hotel.MaxPrice, &hotel.SupplierID, &hotel.SupplierCode,
			&hotel.SupplierName, &hotel.Brand)
		if err != nil {
			continue
		}
		hotels = append(hotels, hotel)
	}

	var total int
	if countRows != nil {
		defer countRows.Close()
		if countRows.Next() {
			countRows.Scan(&total)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"hotels":    hotels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetHotelDetail(c *gin.Context) {
	hotelID := c.Param("id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少酒店ID",
		})
		return
	}

	db := database.GetDB()
	var hotel models.Hotel
	err := db.QueryRow(`
		SELECT id, name, address, city, description, rating, image_url, 
		       price_range, min_price, max_price, supplier_id, supplier_code, supplier_name, brand
		FROM hotels WHERE id = ?`, hotelID).Scan(
		&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
		&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange,
		&hotel.MinPrice, &hotel.MaxPrice, &hotel.SupplierID, &hotel.SupplierCode,
		&hotel.SupplierName, &hotel.Brand)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "酒店不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取酒店详情失败",
		})
		return
	}

	rows, err := db.Query(`
		SELECT id, hotel_id, name, description, price, original_price, capacity, area, bed_type, 
		amenities, image_url, total_count, available_count, supplier_code, supplier_name,
		is_price_controlled, price_control_reason, promotion_tag, payment_type, cancel_policy
		FROM rooms WHERE hotel_id = ?`, hotelID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取房型列表失败",
		})
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var isPriceControlled int
		err := rows.Scan(&room.ID, &room.HotelID, &room.Name, &room.Description, &room.Price,
			&room.OriginalPrice, &room.Capacity, &room.Area, &room.BedType, &room.Amenities,
			&room.ImageURL, &room.TotalCount, &room.AvailableCount, &room.SupplierCode,
			&room.SupplierName, &isPriceControlled, &room.PriceControlReason,
			&room.PromotionTag, &room.PaymentType, &room.CancelPolicy)
		if err != nil {
			continue
		}
		room.IsPriceControlled = isPriceControlled == 1
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"hotel": hotel,
			"rooms": rooms,
		},
	})
}

func GetCities(c *gin.Context) {
	db := database.GetDB()
	rows, err := db.Query("SELECT DISTINCT city FROM hotels ORDER BY city")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取城市列表失败",
		})
		return
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err == nil {
			cities = append(cities, city)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    cities,
	})
}

func GetRoomChannelComparison(c *gin.Context) {
	hotelID := c.Param("id")
	roomID := c.Param("room_id")
	
	if hotelID == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少参数",
		})
		return
	}
	
	db := database.GetDB()
	
	var hotel models.Hotel
	var room models.Room
	var isPriceControlled int
	
	err := db.QueryRow(`
		SELECT id, name, address, city, description, rating, image_url, 
		       price_range, min_price, max_price, supplier_id, supplier_code, supplier_name, brand
		FROM hotels WHERE id = ?`, hotelID).Scan(
		&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
		&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange,
		&hotel.MinPrice, &hotel.MaxPrice, &hotel.SupplierID, &hotel.SupplierCode,
		&hotel.SupplierName, &hotel.Brand)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "酒店不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取酒店信息失败",
		})
		return
	}
	
	err = db.QueryRow(`
		SELECT id, hotel_id, name, description, price, original_price, capacity, area, bed_type, 
		amenities, image_url, total_count, available_count, supplier_code, supplier_name,
		is_price_controlled, price_control_reason, promotion_tag, payment_type, cancel_policy
		FROM rooms WHERE id = ? AND hotel_id = ?`, roomID, hotelID).Scan(
		&room.ID, &room.HotelID, &room.Name, &room.Description, &room.Price,
		&room.OriginalPrice, &room.Capacity, &room.Area, &room.BedType, &room.Amenities,
		&room.ImageURL, &room.TotalCount, &room.AvailableCount, &room.SupplierCode,
		&room.SupplierName, &isPriceControlled, &room.PriceControlReason,
		&room.PromotionTag, &room.PaymentType, &room.CancelPolicy)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "房型不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取房型信息失败",
		})
		return
	}
	room.IsPriceControlled = isPriceControlled == 1
	
	channelPrices := generateMockChannelPrices(&hotel, &room)
	
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    channelPrices,
	})
}

func generateMockChannelPrices(hotel *models.Hotel, room *models.Room) models.HotelChannelComparison {
	currentSupplierCode := hotel.SupplierCode
	
	suppliers := []struct {
		Code     string
		Name     string
		Priority int
		Color    string
		Icon     string
	}{
		{"huazhu", "华住酒店集团", 1, "#FF6B35", "🏨"},
		{"jinjiang", "锦江国际酒店集团", 2, "#1E88E5", "🏩"},
		{"rujia", "如家酒店集团", 3, "#27AE60", "🏠"},
	}
	
	standardRoomType := getStandardRoomType(room.Name)
	
	var channelPrices []models.ChannelPriceInfo
	var minPrice float64 = 999999
	var bestPriceChannel string
	var recommendedChannel string
	var lowestPriority int = 999
	
	for _, supplier := range suppliers {
		isCurrentChannel := supplier.Code == currentSupplierCode
		
		var price float64
		var originalPrice float64
		var availableCount int
		var isPriceControlled bool
		var priceControlReason string
		var promotionTag string
		
		if isCurrentChannel {
			price = room.Price
			originalPrice = room.OriginalPrice
			availableCount = room.AvailableCount
			isPriceControlled = room.IsPriceControlled
			priceControlReason = room.PriceControlReason
			promotionTag = room.PromotionTag
		} else {
			priceVariation := 0.9 + float64(20)/100
			price = room.Price * priceVariation
			price = float64(int(price/10) * 10)
			
			originalPrice = price * (1.05 + float64(15)/100)
			originalPrice = float64(int(originalPrice/10) * 10)
			
			baseAvailable := room.TotalCount / 3
			availableCount = baseAvailable + 5
			if availableCount < 1 {
				availableCount = 1
			}
			
			if supplier.Code == "jinjiang" && room.Name == "标准双床房" {
				isPriceControlled = true
				priceControlReason = "集团协议价"
			}
			
			promotions := []string{"", "限时特惠", "会员专享", "提前预订优惠", "连住优惠"}
			promotionTag = promotions[0]
		}
		
		if price < minPrice && availableCount > 0 {
			minPrice = price
			bestPriceChannel = supplier.Code
		}
		
		if supplier.Priority < lowestPriority && availableCount > 0 {
			lowestPriority = supplier.Priority
			recommendedChannel = supplier.Code
		}
		
		isBestPrice := supplier.Code == bestPriceChannel
		isRecommended := supplier.Code == recommendedChannel
		isPriorityChannel := supplier.Priority == 1
		
		paymentTypes := []string{"现付", "预付", "信用住"}
		cancelPolicies := []string{"免费取消(入住前1天)", "部分取消", "不可取消"}
		
		discount := 0.0
		if originalPrice > 0 && originalPrice > price {
			discount = (originalPrice - price) / originalPrice * 100
			discount = float64(int(discount*10) / 10)
		}
		
		channelPrice := models.ChannelPriceInfo{
			SupplierCode:       supplier.Code,
			SupplierName:       supplier.Name,
			SupplierPriority:   supplier.Priority,
			SupplierColor:      supplier.Color,
			SupplierIcon:       supplier.Icon,
			RoomID:             room.ID,
			RoomName:           room.Name,
			Price:              price,
			OriginalPrice:      originalPrice,
			Discount:           discount,
			AvailableCount:     availableCount,
			TotalCount:         room.TotalCount,
			IsPriceControlled:  isPriceControlled,
			PriceControlReason: priceControlReason,
			IsBestPrice:        isBestPrice,
			IsRecommended:      isRecommended,
			IsPriorityChannel:  isPriorityChannel,
			PromotionTag:       promotionTag,
			PaymentType:        paymentTypes[0],
			CancelPolicy:       cancelPolicies[0],
		}
		
		if !isCurrentChannel {
			channelPrice.PaymentType = paymentTypes[0]
			channelPrice.CancelPolicy = cancelPolicies[0]
		} else {
			channelPrice.PaymentType = room.PaymentType
			channelPrice.CancelPolicy = room.CancelPolicy
		}
		
		channelPrices = append(channelPrices, channelPrice)
	}
	
	maxPrice := room.Price * 1.2
	maxPrice = float64(int(maxPrice/10) * 10)
	
	var bestPriceSupplier string
	for _, supplier := range suppliers {
		if supplier.Code == bestPriceChannel {
			bestPriceSupplier = supplier.Name
			break
		}
	}
	
	return models.HotelChannelComparison{
		HotelID:           hotel.ID,
		HotelName:         hotel.Name,
		RoomID:            room.ID,
		RoomName:          room.Name,
		StandardRoomType:  standardRoomType,
		Channels:          channelPrices,
		MinPrice:          minPrice,
		MaxPrice:          maxPrice,
		BestPrice:         minPrice,
		BestPriceSupplier: bestPriceSupplier,
		BestPriceChannel:  bestPriceChannel,
		RecommendedChannel: recommendedChannel,
	}
}

func getStandardRoomType(roomName string) string {
	standardTypes := map[string]string{
		"标准双床房":   "标准双床房",
		"标准大床房":   "标准大床房",
		"豪华大床房":   "豪华大床房",
		"豪华双床房":   "豪华双床房",
		"商务大床房":   "豪华大床房",
		"商务双床房":   "豪华双床房",
		"行政套房":    "行政套房",
		"豪华套房":    "行政套房",
		"精选套房":    "行政套房",
		"家庭房":     "家庭房",
		"家庭亲子房":   "家庭房",
		"经济房":     "标准大床房",
		"特惠双床房":   "标准双床房",
		"特惠大床房":   "标准大床房",
	}
	
	if t, ok := standardTypes[roomName]; ok {
		return t
	}
	
	return "标准大床房"
}

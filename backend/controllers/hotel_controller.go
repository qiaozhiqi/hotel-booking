package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getSupplierByID(db *sql.DB, supplierID int) *models.SupplierInfo {
	if supplierID <= 0 {
		return nil
	}
	var supplier models.SupplierInfo
	err := db.QueryRow(`
		SELECT id, name, code, description, status, priority, price_control 
		FROM suppliers WHERE id = ?`, supplierID).Scan(
		&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Description,
		&supplier.Status, &supplier.Priority, &supplier.PriceControl)
	if err != nil {
		return nil
	}
	return &supplier
}

func GetHotelList(c *gin.Context) {
	city := c.Query("city")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

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

	if city != "" {
		rows, err = db.Query(`
			SELECT h.id, h.name, h.address, h.city, h.description, h.rating, h.image_url, h.price_range, h.supplier_id,
			       s.id, s.name, s.code, s.description, s.status, s.priority, s.price_control
			FROM hotels h
			LEFT JOIN suppliers s ON h.supplier_id = s.id
			WHERE h.city = ? 
			ORDER BY s.priority DESC, h.rating DESC 
			LIMIT ? OFFSET ?`, city, pageSize, offset)
		
		countRows, _ = db.Query("SELECT COUNT(*) FROM hotels WHERE city = ?", city)
	} else {
		rows, err = db.Query(`
			SELECT h.id, h.name, h.address, h.city, h.description, h.rating, h.image_url, h.price_range, h.supplier_id,
			       s.id, s.name, s.code, s.description, s.status, s.priority, s.price_control
			FROM hotels h
			LEFT JOIN suppliers s ON h.supplier_id = s.id
			ORDER BY s.priority DESC, h.rating DESC 
			LIMIT ? OFFSET ?`, pageSize, offset)
		
		countRows, _ = db.Query("SELECT COUNT(*) FROM hotels")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取酒店列表失败",
		})
		return
	}
	defer rows.Close()

	var hotels []models.HotelWithSupplier
	for rows.Next() {
		var hotel models.HotelWithSupplier
		var supplierID sql.NullInt64
		var supplier models.SupplierInfo
		var hasSupplier bool
		
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
			&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange, &supplierID,
			&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Description,
			&supplier.Status, &supplier.Priority, &supplier.PriceControl)
		
		if err != nil {
			continue
		}
		
		if supplierID.Valid && supplierID.Int64 > 0 {
			hasSupplier = true
		}
		
		if hasSupplier && supplier.ID > 0 {
			hotel.Supplier = &supplier
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

func generateChannelPrices(db *sql.DB, basePrice float64, baseAvailable int, baseSupplierID int) []models.ChannelPrice {
	var channelPrices []models.ChannelPrice
	
	baseSupplier := getSupplierByID(db, baseSupplierID)
	
	if baseSupplier != nil {
		controlledPrice := math.Round(basePrice * baseSupplier.PriceControl)
		channelPrices = append(channelPrices, models.ChannelPrice{
			SupplierID:     baseSupplier.ID,
			SupplierName:   baseSupplier.Name,
			SupplierCode:   baseSupplier.Code,
			Price:          controlledPrice,
			OriginalPrice:  basePrice,
			AvailableCount: baseAvailable,
			Priority:       baseSupplier.Priority,
			IsBestPrice:    false,
		})
	}
	
	simulateSuppliers := []struct {
		ID           int
		Name         string
		Code         string
		Priority     int
		PriceControl float64
	}{
		{101, "模拟供应商A", "sim_a", 5, 0.95},
		{102, "模拟供应商B", "sim_b", 3, 1.02},
		{103, "模拟供应商C", "sim_c", 7, 0.98},
	}
	
	for _, sim := range simulateSuppliers {
		if baseSupplier != nil && sim.ID == baseSupplier.ID {
			continue
		}
		
		priceFluctuation := 0.95 + float64(sim.Priority%3)*0.02
		controlledPrice := math.Round(basePrice * sim.PriceControl * priceFluctuation)
		available := baseAvailable - sim.Priority%3
		
		if available < 0 {
			available = 0
		}
		
		channelPrices = append(channelPrices, models.ChannelPrice{
			SupplierID:     sim.ID,
			SupplierName:   sim.Name,
			SupplierCode:   sim.Code,
			Price:          controlledPrice,
			OriginalPrice:  basePrice,
			AvailableCount: available,
			Priority:       sim.Priority,
			IsBestPrice:    false,
		})
	}
	
	if len(channelPrices) > 0 {
		bestPrice := channelPrices[0].Price
		bestIndex := 0
		
		for i, cp := range channelPrices {
			if cp.AvailableCount > 0 && (cp.Price < bestPrice || channelPrices[bestIndex].AvailableCount <= 0) {
				bestPrice = cp.Price
				bestIndex = i
			}
		}
		
		if channelPrices[bestIndex].AvailableCount > 0 {
			channelPrices[bestIndex].IsBestPrice = true
		}
	}
	
	return channelPrices
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
	var hotel models.HotelWithSupplier
	var supplierID sql.NullInt64
	var supplier models.SupplierInfo
	
	err := db.QueryRow(`
		SELECT h.id, h.name, h.address, h.city, h.description, h.rating, h.image_url, h.price_range, h.supplier_id,
		       s.id, s.name, s.code, s.description, s.status, s.priority, s.price_control
		FROM hotels h
		LEFT JOIN suppliers s ON h.supplier_id = s.id
		WHERE h.id = ?`, hotelID).Scan(
		&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
		&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange, &supplierID,
		&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Description,
		&supplier.Status, &supplier.Priority, &supplier.PriceControl)

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

	if supplierID.Valid && supplierID.Int64 > 0 && supplier.ID > 0 {
		hotel.Supplier = &supplier
	}

	rows, err := db.Query(`
		SELECT id, hotel_id, name, description, price, capacity, area, bed_type, 
		amenities, image_url, total_count, available_count, supplier_id
		FROM rooms WHERE hotel_id = ?`, hotelID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取房型列表失败",
		})
		return
	}
	defer rows.Close()

	var rooms []models.RoomWithChannelPrices
	for rows.Next() {
		var room models.RoomWithChannelPrices
		var roomSupplierID sql.NullInt64
		
		err := rows.Scan(&room.ID, &room.HotelID, &room.Name, &room.Description, &room.Price,
			&room.Capacity, &room.Area, &room.BedType, &room.Amenities, &room.ImageURL,
			&room.TotalCount, &room.AvailableCount, &roomSupplierID)
		if err != nil {
			continue
		}
		
		actualSupplierID := 0
		if roomSupplierID.Valid && roomSupplierID.Int64 > 0 {
			actualSupplierID = int(roomSupplierID.Int64)
		} else if hotel.Supplier != nil {
			actualSupplierID = hotel.Supplier.ID
		}
		
		room.ChannelPrices = generateChannelPrices(db, room.Price, room.AvailableCount, actualSupplierID)
		
		if len(room.ChannelPrices) > 0 {
			bestPrice := room.Price
			for _, cp := range room.ChannelPrices {
				if cp.IsBestPrice && cp.AvailableCount > 0 {
					bestPrice = cp.Price
					break
				}
			}
			room.BestPrice = bestPrice
		} else {
			room.BestPrice = room.Price
		}
		
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

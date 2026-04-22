package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hotel-booking/cache"
	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	SupplierCodeQiuguo = "shiji_qiuguo"
)

func getOrCreateQiuguoSupplier() (int, error) {
	db := database.GetDB()

	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", SupplierCodeQiuguo).Scan(&supplierID)
	if err == nil {
		return supplierID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := db.Exec(`
		INSERT INTO suppliers (name, code, description, api_url, status, priority, price_control)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"秋果集团（石基畅联）", SupplierCodeQiuguo,
		"秋果集团是中国知名的中端酒店集团，旗下拥有秋果酒店等多个品牌。通过石基畅联渠道接入，采用推送模式同步酒店数据。",
		"/api/shiji/qiuguo/push", "active", 3, 0.90)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func HandleQiuguoPush(c *gin.Context) {
	var req models.QiuguoPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.QiuguoPushResponse{
			RequestID: req.RequestID,
			Success:   false,
			Code:      400,
			Message:   "请求参数解析失败: " + err.Error(),
		})
		return
	}

	if req.RequestID == "" {
		c.JSON(http.StatusBadRequest, models.QiuguoPushResponse{
			RequestID: "",
			Success:   false,
			Code:      400,
			Message:   "请求ID不能为空",
		})
		return
	}

	log.Printf("收到秋果数据推送: request_id=%s, push_type=%s, hotels=%d, price_inventories=%d",
		req.RequestID, req.PushType, len(req.Hotels), len(req.PriceInventories))

	supplierID, err := getOrCreateQiuguoSupplier()
	if err != nil {
		log.Printf("获取秋果供应商ID失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.QiuguoPushResponse{
			RequestID: req.RequestID,
			Success:   false,
			Code:      500,
			Message:   "供应商初始化失败",
		})
		return
	}

	var processedCount int
	var failedCount int
	var errorDetails []string

	if len(req.Hotels) > 0 {
		for _, hotel := range req.Hotels {
			err := processQiuguoHotel(supplierID, hotel)
			if err != nil {
				failedCount++
				errorDetails = append(errorDetails, fmt.Sprintf("酒店 %s(%s): %v", hotel.Name, hotel.SupplierHotelID, err))
				log.Printf("处理秋果酒店失败: hotel_id=%s, error=%v", hotel.SupplierHotelID, err)
			} else {
				processedCount++
			}
		}
	}

	if len(req.PriceInventories) > 0 {
		for _, pi := range req.PriceInventories {
			err := processQiuguoPriceInventory(supplierID, pi)
			if err != nil {
				failedCount++
				errorDetails = append(errorDetails, fmt.Sprintf("价格库存 %s/%s/%s: %v", pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, err))
				log.Printf("处理秋果价格库存失败: hotel_id=%s, room_id=%s, date=%s, error=%v",
					pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, err)
			} else {
				processedCount++
			}
		}
	}

	response := models.QiuguoPushResponse{
		RequestID:      req.RequestID,
		Success:        true,
		Code:           200,
		Message:        "数据推送处理完成",
		ProcessedCount: processedCount,
		FailedCount:    failedCount,
	}

	if len(errorDetails) > 0 {
		response.ErrorDetails = errorDetails
		if processedCount == 0 {
			response.Success = false
			response.Code = 500
			response.Message = "数据推送处理失败"
		}
	}

	log.Printf("秋果数据推送处理完成: request_id=%s, processed=%d, failed=%d",
		req.RequestID, processedCount, failedCount)

	c.JSON(http.StatusOK, response)
}

func processQiuguoHotel(supplierID int, hotel models.QiuguoPushHotelData) error {
	db := database.GetDB()

	var existingID int
	var localHotelID int
	err := db.QueryRow(`
		SELECT id, local_hotel_id FROM supplier_hotels 
		WHERE supplier_id = ? AND supplier_hotel_id = ?`,
		supplierID, hotel.SupplierHotelID).Scan(&existingID, &localHotelID)

	rawData, _ := json.Marshal(hotel)

	if err == sql.ErrNoRows {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		hotelResult, err := tx.Exec(`
			INSERT INTO hotels (name, address, city, description, rating, image_url, price_range, supplier_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			hotel.Name, hotel.Address, hotel.City, hotel.Description,
			hotel.Rating, hotel.ImageURL, hotel.PriceRange, supplierID)
		if err != nil {
			tx.Rollback()
			return err
		}

		newLocalHotelID, _ := hotelResult.LastInsertId()

		_, err = tx.Exec(`
			INSERT INTO supplier_hotels (supplier_id, supplier_hotel_id, local_hotel_id, hotel_name, 
			city, address, rating, image_url, price_range, raw_data, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			supplierID, hotel.SupplierHotelID, newLocalHotelID, hotel.Name,
			hotel.City, hotel.Address, hotel.Rating, hotel.ImageURL,
			hotel.PriceRange, string(rawData), "synced")
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, room := range hotel.Rooms {
			roomResult, err := tx.Exec(`
				INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
				amenities, image_url, total_count, available_count, supplier_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				newLocalHotelID, room.Name, room.Description, 0.0,
				room.Capacity, room.Area, room.BedType, room.Amenities,
				room.ImageURL, room.TotalCount, room.TotalCount, supplierID)
			if err != nil {
				continue
			}

			localRoomID, _ := roomResult.LastInsertId()
			roomRawData, _ := json.Marshal(room)
			tx.Exec(`
				INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
				local_room_id, room_name, description, price, capacity, area, bed_type, 
				amenities, image_url, total_count, available_count, raw_data, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				supplierID, hotel.SupplierHotelID, room.SupplierRoomID,
				localRoomID, room.Name, room.Description, 0.0,
				room.Capacity, room.Area, room.BedType, room.Amenities,
				room.ImageURL, room.TotalCount, room.TotalCount,
				string(roomRawData), "synced")
		}

		tx.Commit()
		log.Printf("新增秋果酒店: hotel_id=%s, name=%s", hotel.SupplierHotelID, hotel.Name)
	} else if err == nil {
		_, err = db.Exec(`
			UPDATE hotels SET name = ?, address = ?, city = ?, description = ?, 
			rating = ?, image_url = ?, price_range = ? WHERE id = ?`,
			hotel.Name, hotel.Address, hotel.City, hotel.Description,
			hotel.Rating, hotel.ImageURL, hotel.PriceRange, localHotelID)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			UPDATE supplier_hotels SET hotel_name = ?, city = ?, address = ?, rating = ?, 
			image_url = ?, price_range = ?, raw_data = ?, status = ?, updated_at = ? WHERE id = ?`,
			hotel.Name, hotel.City, hotel.Address, hotel.Rating,
			hotel.ImageURL, hotel.PriceRange, string(rawData), "synced", time.Now(), existingID)
		if err != nil {
			return err
		}

		for _, room := range hotel.Rooms {
			var existingRoomID int
			var localRoomID int
			err := db.QueryRow(`
				SELECT id, local_room_id FROM supplier_rooms 
				WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
				supplierID, hotel.SupplierHotelID, room.SupplierRoomID).Scan(&existingRoomID, &localRoomID)

			if err == sql.ErrNoRows {
				roomResult, _ := db.Exec(`
					INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, supplier_id)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					localHotelID, room.Name, room.Description, 0.0,
					room.Capacity, room.Area, room.BedType, room.Amenities,
					room.ImageURL, room.TotalCount, room.TotalCount, supplierID)

				newLocalRoomID, _ := roomResult.LastInsertId()
				roomRawData, _ := json.Marshal(room)
				db.Exec(`
					INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
					local_room_id, room_name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, raw_data, status)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					supplierID, hotel.SupplierHotelID, room.SupplierRoomID,
					newLocalRoomID, room.Name, room.Description, 0.0,
					room.Capacity, room.Area, room.BedType, room.Amenities,
					room.ImageURL, room.TotalCount, room.TotalCount,
					string(roomRawData), "synced")
			} else if err == nil {
				db.Exec(`
					UPDATE rooms SET name = ?, description = ?, capacity = ?, 
					area = ?, bed_type = ?, amenities = ?, image_url = ?, 
					total_count = ? WHERE id = ?`,
					room.Name, room.Description, room.Capacity,
					room.Area, room.BedType, room.Amenities, room.ImageURL,
					room.TotalCount, localRoomID)

				roomRawData, _ := json.Marshal(room)
				db.Exec(`
					UPDATE supplier_rooms SET room_name = ?, description = ?, capacity = ?, 
					area = ?, bed_type = ?, amenities = ?, image_url = ?, 
					total_count = ?, raw_data = ?, status = ?, updated_at = ? WHERE id = ?`,
					room.Name, room.Description, room.Capacity,
					room.Area, room.BedType, room.Amenities, room.ImageURL,
					room.TotalCount, string(roomRawData), "synced", time.Now(), existingRoomID)
			}
		}
		log.Printf("更新秋果酒店: hotel_id=%s, name=%s", hotel.SupplierHotelID, hotel.Name)
	}

	return nil
}

func processQiuguoPriceInventory(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	originalPrice := pi.OriginalPrice
	if originalPrice == 0 {
		originalPrice = pi.Price
	}

	totalCount := pi.TotalCount
	if totalCount == 0 {
		totalCount = 10
	}

	if cfg.EnableRedisCache && cache.GetRedis() != nil {
		err := cache.SetPriceInventory(supplierID, pi)
		if err != nil {
			log.Printf("写入Redis价格库存失败: %v", err)
		} else {
			log.Printf("Redis缓存价格库存: hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
				pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
		}
	}

	var existingID int
	err := db.QueryRow(`
		SELECT id FROM price_inventory 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ? AND date = ?`,
		supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date).Scan(&existingID)

	if err == sql.ErrNoRows {
		_, err = db.Exec(`
			INSERT INTO price_inventory (supplier_id, supplier_hotel_id, supplier_room_id, date, 
			price, original_price, available_count, total_count, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date,
			pi.Price, originalPrice, pi.AvailableCount, totalCount, "active")
		if err != nil {
			return err
		}
		log.Printf("新增秋果价格库存(数据库): hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
			pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
	} else if err == nil {
		_, err = db.Exec(`
			UPDATE price_inventory SET price = ?, original_price = ?, available_count = ?, 
			total_count = ?, status = ?, updated_at = ? WHERE id = ?`,
			pi.Price, originalPrice, pi.AvailableCount, totalCount, "active", time.Now(), existingID)
		if err != nil {
			return err
		}
		log.Printf("更新秋果价格库存(数据库): hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
			pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
	} else {
		return err
	}

	err = updateRoomCurrentPrice(supplierID, pi.SupplierHotelID, pi.SupplierRoomID)
	if err != nil {
		log.Printf("更新房间当前价格失败: %v", err)
	}

	return nil
}

func updateRoomCurrentPrice(supplierID int, supplierHotelID, supplierRoomID string) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	today := time.Now().Format("2006-01-02")

	var currentPrice float64
	var currentAvailable int
	var found bool

	if cfg.EnableRedisCache && cache.GetRedis() != nil {
		pi, err := cache.GetPriceInventory(supplierID, supplierHotelID, supplierRoomID, today)
		if err == nil && pi != nil {
			currentPrice = pi.Price
			currentAvailable = pi.AvailableCount
			found = true
			log.Printf("从Redis获取价格库存: hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
				supplierHotelID, supplierRoomID, today, currentPrice, currentAvailable)
		} else if err != redis.Nil {
			log.Printf("从Redis获取价格库存失败: %v", err)
		}
	}

	if !found {
		err := db.QueryRow(`
			SELECT price, available_count FROM price_inventory 
			WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ? AND date = ?`,
			supplierID, supplierHotelID, supplierRoomID, today).Scan(&currentPrice, &currentAvailable)

		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("从数据库获取价格库存: hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
			supplierHotelID, supplierRoomID, today, currentPrice, currentAvailable)
	}

	var localRoomID int
	err := db.QueryRow(`
		SELECT local_room_id FROM supplier_rooms 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
		supplierID, supplierHotelID, supplierRoomID).Scan(&localRoomID)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE rooms SET price = ?, available_count = ?, updated_at = ? WHERE id = ?`,
		currentPrice, currentAvailable, time.Now(), localRoomID)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE supplier_rooms SET price = ?, available_count = ?, updated_at = ? 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
		currentPrice, currentAvailable, time.Now(), supplierID, supplierHotelID, supplierRoomID)

	return err
}

func GetQiuguoSyncStatus(c *gin.Context) {
	cfg := config.GetConfig()
	db := database.GetDB()

	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", SupplierCodeQiuguo).Scan(&supplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, models.Response{
				Code:    200,
				Message: "获取成功",
				Data: map[string]interface{}{
					"supplier_code":   SupplierCodeQiuguo,
					"supplier_name":   "秋果集团（石基畅联）",
					"status":          "not_initialized",
					"total_hotels":    0,
					"total_rooms":     0,
					"price_inventory_count": 0,
					"last_sync_time":  nil,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询供应商失败",
		})
		return
	}

	var hotelCount int
	db.QueryRow("SELECT COUNT(*) FROM supplier_hotels WHERE supplier_id = ?", supplierID).Scan(&hotelCount)

	var roomCount int
	db.QueryRow("SELECT COUNT(*) FROM supplier_rooms WHERE supplier_id = ?", supplierID).Scan(&roomCount)

	var priceInventoryCount int64
	var minDate, maxDate string
	var lastSyncTime string
	var fromRedis bool

	if cfg.EnableRedisCache && cache.GetRedis() != nil {
		var err error
		priceInventoryCount, err = cache.GetPriceInventoryCount(supplierID)
		if err == nil {
			minDate, maxDate, err = cache.GetPriceInventoryDateRange(supplierID)
			if err == nil {
				lastSyncTime, err = cache.GetPriceInventoryLastUpdate(supplierID)
				if err == nil {
					fromRedis = true
					log.Printf("从Redis获取同步状态: count=%d, min=%s, max=%s, last=%s",
						priceInventoryCount, minDate, maxDate, lastSyncTime)
				}
			}
		}
		if err != nil {
			log.Printf("从Redis获取同步状态失败: %v", err)
		}
	}

	if !fromRedis {
		var dbCount int
		db.QueryRow("SELECT COUNT(*) FROM price_inventory WHERE supplier_id = ?", supplierID).Scan(&dbCount)
		priceInventoryCount = int64(dbCount)

		db.QueryRow(`
			SELECT MIN(date), MAX(date) FROM price_inventory WHERE supplier_id = ?`,
			supplierID).Scan(&minDate, &maxDate)

		var lastUpdatedAt sql.NullTime
		db.QueryRow(`
			SELECT MAX(updated_at) FROM price_inventory WHERE supplier_id = ?`,
			supplierID).Scan(&lastUpdatedAt)

		if lastUpdatedAt.Valid {
			lastSyncTime = lastUpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		log.Printf("从数据库获取同步状态: count=%d, min=%s, max=%s, last=%s",
			priceInventoryCount, minDate, maxDate, lastSyncTime)
	}

	dataSource := "database"
	if fromRedis {
		dataSource = "redis"
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"supplier_code":        SupplierCodeQiuguo,
			"supplier_name":        "秋果集团（石基畅联）",
			"status":               "active",
			"total_hotels":         hotelCount,
			"total_rooms":          roomCount,
			"price_inventory_count": priceInventoryCount,
			"price_inventory_date_range": map[string]string{
				"min": minDate,
				"max": maxDate,
			},
			"last_sync_time": lastSyncTime,
			"data_source":     dataSource,
		},
	})
}

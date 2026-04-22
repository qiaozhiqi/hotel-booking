package services

import (
	"database/sql"
	"fmt"
	"hotel-booking/cache"
	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/models"
	"hotel-booking/suppliers"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type PriceInventoryService struct{}

func NewPriceInventoryService() *PriceInventoryService {
	return &PriceInventoryService{}
}

func (s *PriceInventoryService) SavePriceInventory(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
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

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
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
		log.Printf("新增价格库存(数据库): hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
			pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
	} else if err == nil {
		_, err = db.Exec(`
			UPDATE price_inventory SET price = ?, original_price = ?, available_count = ?, 
			total_count = ?, status = ?, updated_at = ? WHERE id = ?`,
			pi.Price, originalPrice, pi.AvailableCount, totalCount, "active", time.Now(), existingID)
		if err != nil {
			return err
		}
		log.Printf("更新价格库存(数据库): hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
			pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
	} else {
		return err
	}

	err = s.UpdateRoomCurrentPrice(supplierID, pi.SupplierHotelID, pi.SupplierRoomID)
	if err != nil {
		log.Printf("更新房间当前价格失败: %v", err)
	}

	return nil
}

func (s *PriceInventoryService) SaveSupplierPriceInventory(supplierID int, pi suppliers.SupplierPriceInventoryData) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		err := cache.SetSupplierPriceInventory(supplierID, pi)
		if err != nil {
			log.Printf("写入Redis价格库存失败: %v", err)
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
			pi.Price, pi.Price, pi.AvailableCount, 10, "active")
		if err != nil {
			return err
		}
	} else if err == nil {
		_, err = db.Exec(`
			UPDATE price_inventory SET price = ?, available_count = ?, 
			status = ?, updated_at = ? WHERE id = ?`,
			pi.Price, pi.AvailableCount, "active", time.Now(), existingID)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (s *PriceInventoryService) GetPriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	cfg := config.GetConfig()
	db := database.GetDB()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		pi, err := cache.GetPriceInventory(supplierID, supplierHotelID, supplierRoomID, date)
		if err == nil {
			log.Printf("从Redis获取价格库存: hotel_id=%s, room_id=%s, date=%s",
				supplierHotelID, supplierRoomID, date)
			return pi, nil
		} else if err != redis.Nil {
			log.Printf("从Redis获取价格库存失败: %v", err)
		}
	}

	pi := &models.PriceInventory{}
	err := db.QueryRow(`
		SELECT id, supplier_id, supplier_hotel_id, supplier_room_id, date, 
		price, original_price, available_count, total_count, status, created_at, updated_at
		FROM price_inventory 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ? AND date = ?`,
		supplierID, supplierHotelID, supplierRoomID, date).Scan(
		&pi.ID, &pi.SupplierID, &pi.SupplierHotelID, &pi.SupplierRoomID, &pi.Date,
		&pi.Price, &pi.OriginalPrice, &pi.AvailableCount, &pi.TotalCount,
		&pi.Status, &pi.CreatedAt, &pi.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	log.Printf("从数据库获取价格库存: hotel_id=%s, room_id=%s, date=%s",
		supplierHotelID, supplierRoomID, date)

	return pi, nil
}

func (s *PriceInventoryService) UpdateRoomCurrentPrice(supplierID int, supplierHotelID, supplierRoomID string) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	today := time.Now().Format("2006-01-02")

	var currentPrice float64
	var currentAvailable int
	var found bool

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		pi, err := cache.GetPriceInventory(supplierID, supplierHotelID, supplierRoomID, today)
		if err == nil && pi != nil {
			currentPrice = pi.Price
			currentAvailable = pi.AvailableCount
			found = true
			log.Printf("从Redis获取房间当前价格: hotel_id=%s, room_id=%s, price=%.2f, available=%d",
				supplierHotelID, supplierRoomID, currentPrice, currentAvailable)
		} else if err != redis.Nil {
			log.Printf("从Redis获取房间当前价格失败: %v", err)
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
		log.Printf("从数据库获取房间当前价格: hotel_id=%s, room_id=%s, price=%.2f, available=%d",
			supplierHotelID, supplierRoomID, currentPrice, currentAvailable)
	}

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		cache.SetRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID, currentPrice, currentAvailable)
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

func (s *PriceInventoryService) GetSyncStatus(supplierID int) (map[string]interface{}, error) {
	cfg := config.GetConfig()
	db := database.GetDB()

	var priceInventoryCount int64
	var minDate, maxDate string
	var lastSyncTime string
	var dataSource string

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		var err error
		priceInventoryCount, err = cache.GetPriceInventoryCount(supplierID)
		if err == nil {
			minDate, maxDate, err = cache.GetPriceInventoryDateRange(supplierID)
			if err == nil {
				lastSyncTime, err = cache.GetPriceInventoryLastUpdate(supplierID)
				if err == nil {
					dataSource = "redis"
					log.Printf("从Redis获取同步状态: supplier_id=%d, count=%d", supplierID, priceInventoryCount)
				}
			}
		}
		if err != nil {
			log.Printf("从Redis获取同步状态失败: %v", err)
		}
	}

	if dataSource == "" {
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
		dataSource = "database"
		log.Printf("从数据库获取同步状态: supplier_id=%d, count=%d", supplierID, priceInventoryCount)
	}

	return map[string]interface{}{
		"price_inventory_count": priceInventoryCount,
		"price_inventory_date_range": map[string]string{
			"min": minDate,
			"max": maxDate,
		},
		"last_sync_time": lastSyncTime,
		"data_source":     dataSource,
	}, nil
}

func (s *PriceInventoryService) BatchSavePriceInventory(supplierID int, priceInventories []models.QiuguoPushPriceInventoryData) error {
	cfg := config.GetConfig()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		err := cache.BatchSetPriceInventory(supplierID, priceInventories)
		if err != nil {
			log.Printf("批量写入Redis价格库存失败: %v", err)
		} else {
			log.Printf("批量缓存Redis价格库存: %d 条", len(priceInventories))
		}
	}

	for _, pi := range priceInventories {
		err := s.SavePriceInventory(supplierID, pi)
		if err != nil {
			log.Printf("保存价格库存失败: %v", err)
		}
	}

	return nil
}

func (s *PriceInventoryService) ClearPriceInventoryBySupplier(supplierID int) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		err := cache.ClearPriceInventoryBySupplier(supplierID)
		if err != nil {
			log.Printf("清除Redis价格库存失败: %v", err)
		} else {
			log.Printf("已清除Redis价格库存: supplier_id=%d", supplierID)
		}
	}

	_, err := db.Exec("DELETE FROM price_inventory WHERE supplier_id = ?", supplierID)
	if err != nil {
		log.Printf("清除数据库价格库存失败: %v", err)
		return err
	}

	log.Printf("已清除数据库价格库存: supplier_id=%d", supplierID)
	return nil
}

func (s *PriceInventoryService) GetPriceInventoryByDateRange(supplierID int, supplierHotelID, supplierRoomID string, startDate, endDate string) ([]*models.PriceInventory, error) {
	cfg := config.GetConfig()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		dataList, err := cache.GetPriceInventoryByDateRange(supplierID, supplierHotelID, supplierRoomID, startDate, endDate)
		if err == nil && len(dataList) > 0 {
			result := make([]*models.PriceInventory, len(dataList))
			for i, data := range dataList {
				result[i] = data.ToModel()
			}
			log.Printf("从Redis获取价格库存范围: hotel_id=%s, room_id=%s, %s 至 %s, 共%d条",
				supplierHotelID, supplierRoomID, startDate, endDate, len(result))
			return result, nil
		}
	}

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT id, supplier_id, supplier_hotel_id, supplier_room_id, date, 
		price, original_price, available_count, total_count, status, created_at, updated_at
		FROM price_inventory 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ? 
		AND date >= ? AND date <= ?
		ORDER BY date`,
		supplierID, supplierHotelID, supplierRoomID, startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.PriceInventory
	for rows.Next() {
		pi := &models.PriceInventory{}
		err := rows.Scan(
			&pi.ID, &pi.SupplierID, &pi.SupplierHotelID, &pi.SupplierRoomID, &pi.Date,
			&pi.Price, &pi.OriginalPrice, &pi.AvailableCount, &pi.TotalCount,
			&pi.Status, &pi.CreatedAt, &pi.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, pi)
	}

	log.Printf("从数据库获取价格库存范围: hotel_id=%s, room_id=%s, %s 至 %s, 共%d条",
		supplierHotelID, supplierRoomID, startDate, endDate, len(result))

	return result, nil
}

func (s *PriceInventoryService) GetRoomCurrentPrice(supplierID int, supplierHotelID, supplierRoomID string) (*cache.RoomCurrentPrice, error) {
	cfg := config.GetConfig()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		rcp, err := cache.GetRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID)
		if err == nil {
			return rcp, nil
		} else if err != redis.Nil {
			log.Printf("从Redis获取房间当前价格失败: %v", err)
		}
	}

	s.UpdateRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID)

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		return cache.GetRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID)
	}

	return nil, fmt.Errorf("无法获取房间当前价格")
}

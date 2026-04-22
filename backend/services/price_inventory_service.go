package services

import (
	"fmt"
	"hotel-booking/cache"
	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/models"
	"hotel-booking/suppliers"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type PriceInventoryService struct {
	roomSummaryMutex sync.Mutex
	hotelSummaryMutex sync.Mutex
}

func NewPriceInventoryService() *PriceInventoryService {
	return &PriceInventoryService{}
}

func (s *PriceInventoryService) SavePriceInventory(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	cfg := config.GetConfig()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		err := cache.SetPriceInventory(supplierID, pi)
		if err != nil {
			log.Printf("写入Redis价格库存失败: %v", err)
		} else {
			log.Printf("Redis缓存价格库存: hotel_id=%s, room_id=%s, date=%s, price=%.2f, available=%d",
				pi.SupplierHotelID, pi.SupplierRoomID, pi.Date, pi.Price, pi.AvailableCount)
		}
	}

	err := s.updateRoomPriceSummary(supplierID, pi.SupplierHotelID, pi.SupplierRoomID)
	if err != nil {
		log.Printf("更新房间价格汇总失败: %v", err)
	}

	s.UpdateRoomCurrentPrice(supplierID, pi.SupplierHotelID, pi.SupplierRoomID)

	return nil
}

func (s *PriceInventoryService) SaveSupplierPriceInventory(supplierID int, pi suppliers.SupplierPriceInventoryData) error {
	cfg := config.GetConfig()

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		err := cache.SetSupplierPriceInventory(supplierID, pi)
		if err != nil {
			log.Printf("写入Redis价格库存失败: %v", err)
		}
	}

	err := s.updateRoomPriceSummary(supplierID, pi.SupplierHotelID, pi.SupplierRoomID)
	if err != nil {
		log.Printf("更新房间价格汇总失败: %v", err)
	}

	return nil
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

	roomSet := make(map[string]struct{})
	for _, pi := range priceInventories {
		roomKey := fmt.Sprintf("%s:%s", pi.SupplierHotelID, pi.SupplierRoomID)
		roomSet[roomKey] = struct{}{}
	}

	for roomKey := range roomSet {
		var hotelID, roomID string
		fmt.Sscanf(roomKey, "%s:%s", &hotelID, &roomID)
		s.updateRoomPriceSummary(supplierID, hotelID, roomID)
	}

	return nil
}

func (s *PriceInventoryService) GetPriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	cfg := config.GetConfig()

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

	return nil, fmt.Errorf("价格库存数据不存在")
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

	return nil, fmt.Errorf("价格库存数据不存在")
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
		return nil
	}

	if cfg.EnableRedisCache && cache.IsRedisAvailable() {
		cache.SetRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID, currentPrice, currentAvailable)
	}

	var localRoomID int
	err := db.QueryRow(`
		SELECT local_room_id FROM supplier_rooms 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
		supplierID, supplierHotelID, supplierRoomID).Scan(&localRoomID)

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

func (s *PriceInventoryService) updateRoomPriceSummary(supplierID int, supplierHotelID, supplierRoomID string) error {
	cfg := config.GetConfig()
	db := database.GetDB()

	if !cfg.EnableRedisCache || !cache.IsRedisAvailable() {
		return nil
	}

	s.roomSummaryMutex.Lock()
	defer s.roomSummaryMutex.Unlock()

	keys, err := cache.GetAllPriceInventoryKeys(supplierID)
	if err != nil {
		return err
	}

	var roomPrices []float64
	var roomDates []string
	hasInventory := false
	totalCount := 0
	priceSum := 0.0

	for _, key := range keys {
		parts := splitKey(key)
		if len(parts) >= 5 {
			if parts[2] == supplierHotelID && parts[3] == supplierRoomID {
				pi, err := cache.GetPriceInventoryData(supplierID, supplierHotelID, supplierRoomID, parts[4])
				if err != nil {
					continue
				}

				roomPrices = append(roomPrices, pi.Price)
				roomDates = append(roomDates, pi.Date)

				if pi.AvailableCount > 0 {
					hasInventory = true
				}

				totalCount = pi.TotalCount
				priceSum += pi.Price
			}
		}
	}

	if len(roomPrices) == 0 {
		return nil
	}

	sort.Float64s(roomPrices)
	sort.Strings(roomDates)

	minPrice := roomPrices[0]
	maxPrice := roomPrices[len(roomPrices)-1]
	avgPrice := priceSum / float64(len(roomPrices))
	priceRange := fmt.Sprintf("¥%d-¥%d", int(minPrice), int(maxPrice))
	dateRangeStart := roomDates[0]
	dateRangeEnd := roomDates[len(roomDates)-1]

	var existingID int
	err = db.QueryRow(`
		SELECT id FROM room_price_summary 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
		supplierID, supplierHotelID, supplierRoomID).Scan(&existingID)

	if err != nil && err != fmt.Errorf("sql: no rows in result set") {
		if err.Error() != "sql: no rows in result set" && err.Error() != "sql: Rows are closed" {
			if err.Error() != "sql: no rows in result set" {
				log.Printf("查询房间价格汇总失败: %v", err)
			}
		}
	}

	if err == nil && existingID > 0 {
		_, err = db.Exec(`
			UPDATE room_price_summary SET 
			min_price = ?, max_price = ?, avg_price = ?, price_range = ?,
			has_inventory = ?, total_count = ?, 
			date_range_start = ?, date_range_end = ?, updated_at = ?
			WHERE id = ?`,
			minPrice, maxPrice, avgPrice, priceRange,
			hasInventory, totalCount,
			dateRangeStart, dateRangeEnd, time.Now(), existingID)
		if err != nil {
			log.Printf("更新房间价格汇总失败: %v", err)
			return err
		}
		log.Printf("更新房间价格汇总: hotel_id=%s, room_id=%s, price_range=%s",
			supplierHotelID, supplierRoomID, priceRange)
	} else {
		_, err = db.Exec(`
			INSERT INTO room_price_summary 
			(supplier_id, supplier_hotel_id, supplier_room_id, 
			min_price, max_price, avg_price, price_range,
			has_inventory, total_count, 
			date_range_start, date_range_end, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			supplierID, supplierHotelID, supplierRoomID,
			minPrice, maxPrice, avgPrice, priceRange,
			hasInventory, totalCount,
			dateRangeStart, dateRangeEnd, time.Now())
		if err != nil {
			log.Printf("插入房间价格汇总失败: %v", err)
			return err
		}
		log.Printf("新增房间价格汇总: hotel_id=%s, room_id=%s, price_range=%s",
			supplierHotelID, supplierRoomID, priceRange)
	}

	return nil
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
		var roomCount = 0
		db.QueryRow("SELECT COUNT(*) FROM room_price_summary WHERE supplier_id = ?", supplierID).Scan(&roomCount)
		priceInventoryCount = int64(roomCount)

		var dateStart, dateEnd string
		db.QueryRow(`
			SELECT MIN(date_range_start), MAX(date_range_end) 
			FROM room_price_summary WHERE supplier_id = ?`,
			supplierID).Scan(&dateStart, &dateEnd)
		minDate = dateStart
		maxDate = dateEnd

		var lastUpdatedAt time.Time
		db.QueryRow(`
			SELECT MAX(updated_at) FROM room_price_summary WHERE supplier_id = ?`,
			supplierID).Scan(&lastUpdatedAt)

		if !lastUpdatedAt.IsZero() {
			lastSyncTime = lastUpdatedAt.Format("2006-01-02 15:04:05")
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

	_, err := db.Exec("DELETE FROM room_price_summary WHERE supplier_id = ?", supplierID)
	if err != nil {
		log.Printf("清除数据库房间价格汇总失败: %v", err)
		return err
	}

	log.Printf("已清除数据库房间价格汇总: supplier_id=%d", supplierID)
	return nil
}

func (s *PriceInventoryService) GetRoomPriceSummary(supplierID int, supplierHotelID, supplierRoomID string) (*models.RoomPriceSummary, error) {
	db := database.GetDB()

	var summary models.RoomPriceSummary
	err := db.QueryRow(`
		SELECT id, supplier_id, supplier_hotel_id, supplier_room_id,
		min_price, max_price, avg_price, price_range,
		has_inventory, total_count, date_range_start, date_range_end, updated_at
		FROM room_price_summary 
		WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
		supplierID, supplierHotelID, supplierRoomID).Scan(
		&summary.ID, &summary.SupplierID, &summary.SupplierHotelID, &summary.SupplierRoomID,
		&summary.MinPrice, &summary.MaxPrice, &summary.AvgPrice, &summary.PriceRange,
		&summary.HasInventory, &summary.TotalCount, &summary.DateRangeStart, &summary.DateRangeEnd, &summary.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (s *PriceInventoryService) GetHotelPriceSummary(supplierID int, supplierHotelID string) (*models.HotelPriceSummary, error) {
	db := database.GetDB()

	var totalRooms, roomsWithInventory int
	var minPrice, maxPrice, avgPrice float64
	var dateRangeStart, dateRangeEnd string
	var updatedAt time.Time

	err := db.QueryRow(`
		SELECT 
			COUNT(*) as total_rooms,
			SUM(CASE WHEN has_inventory = 1 THEN 1 ELSE 0 END) as rooms_with_inventory,
			MIN(min_price) as min_price,
			MAX(max_price) as max_price,
			AVG(avg_price) as avg_price,
			MIN(date_range_start) as date_range_start,
			MAX(date_range_end) as date_range_end,
			MAX(updated_at) as updated_at
		FROM room_price_summary 
		WHERE supplier_id = ? AND supplier_hotel_id = ?`,
		supplierID, supplierHotelID).Scan(
		&totalRooms, &roomsWithInventory,
		&minPrice, &maxPrice, &avgPrice,
		&dateRangeStart, &dateRangeEnd, &updatedAt)

	if err != nil {
		return nil, err
	}

	priceRange := ""
	if totalRooms > 0 {
		priceRange = fmt.Sprintf("¥%d-¥%d", int(minPrice), int(maxPrice))
	}

	return &models.HotelPriceSummary{
		SupplierID:         supplierID,
		SupplierHotelID: supplierHotelID,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
		AvgPrice:        avgPrice,
		PriceRange:      priceRange,
		TotalRooms:      totalRooms,
		RoomsWithInventory: roomsWithInventory,
		DateRangeStart:  dateRangeStart,
		DateRangeEnd:    dateRangeEnd,
		UpdatedAt:       updatedAt,
	}, nil
}

func splitKey(key string) []string {
	var parts []string
	var current []byte

	for _, c := range key {
		if c == ':' {
			parts = append(parts, string(current))
			current = nil
		} else {
			current = append(current, byte(c))
		}
	}
	if len(current) > 0 {
		parts = append(parts, string(current))
	}
	return parts
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-booking/config"
	"hotel-booking/models"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

const (
	PriceInventoryPrefix = "price_inv"
	PriceInventoryTTL    = 30 * 24 * time.Hour
)

func InitRedis() error {
	cfg := config.GetConfig()
	
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
		return err
	}

	log.Println("Redis连接成功")
	return nil
}

func GetRedis() *redis.Client {
	return RDB
}

func buildPriceInventoryKey(supplierID int, supplierHotelID, supplierRoomID, date string) string {
	return fmt.Sprintf("%s:%d:%s:%s:%s", PriceInventoryPrefix, supplierID, supplierHotelID, supplierRoomID, date)
}

func SetPriceInventory(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

	originalPrice := pi.OriginalPrice
	if originalPrice == 0 {
		originalPrice = pi.Price
	}

	totalCount := pi.TotalCount
	if totalCount == 0 {
		totalCount = 10
	}

	data := map[string]interface{}{
		"supplier_id":       supplierID,
		"supplier_hotel_id": pi.SupplierHotelID,
		"supplier_room_id":  pi.SupplierRoomID,
		"date":              pi.Date,
		"price":             pi.Price,
		"original_price":    originalPrice,
		"available_count":   pi.AvailableCount,
		"total_count":       totalCount,
		"status":            "active",
		"updated_at":        time.Now().Format("2006-01-02 15:04:05"),
	}

	err := RDB.HSet(ctx, key, data).Err()
	if err != nil {
		return err
	}

	err = RDB.Expire(ctx, key, PriceInventoryTTL).Err()
	if err != nil {
		log.Printf("设置Redis过期时间失败: %v", err)
	}

	return nil
}

func GetPriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)

	data, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, redis.Nil
	}

	pi := &models.PriceInventory{}
	if supplierIDStr, ok := data["supplier_id"]; ok {
		var sid int
		fmt.Sscanf(supplierIDStr, "%d", &sid)
		pi.SupplierID = sid
	}
	pi.SupplierHotelID = data["supplier_hotel_id"]
	pi.SupplierRoomID = data["supplier_room_id"]
	pi.Date = data["date"]
	pi.Status = data["status"]

	if priceStr, ok := data["price"]; ok {
		var price float64
		fmt.Sscanf(priceStr, "%f", &price)
		pi.Price = price
	}

	if originalPriceStr, ok := data["original_price"]; ok {
		var originalPrice float64
		fmt.Sscanf(originalPriceStr, "%f", &originalPrice)
		pi.OriginalPrice = originalPrice
	}

	if availableCountStr, ok := data["available_count"]; ok {
		var availableCount int
		fmt.Sscanf(availableCountStr, "%d", &availableCount)
		pi.AvailableCount = availableCount
	}

	if totalCountStr, ok := data["total_count"]; ok {
		var totalCount int
		fmt.Sscanf(totalCountStr, "%d", &totalCount)
		pi.TotalCount = totalCount
	}

	return pi, nil
}

func DeletePriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) error {
	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)
	return RDB.Del(ctx, key).Err()
}

func BatchSetPriceInventory(supplierID int, priceInventories []models.QiuguoPushPriceInventoryData) error {
	pipe := RDB.Pipeline()
	ctx := context.Background()

	for _, pi := range priceInventories {
		key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

		originalPrice := pi.OriginalPrice
		if originalPrice == 0 {
			originalPrice = pi.Price
		}

		totalCount := pi.TotalCount
		if totalCount == 0 {
			totalCount = 10
		}

		data := map[string]interface{}{
			"supplier_id":       supplierID,
			"supplier_hotel_id": pi.SupplierHotelID,
			"supplier_room_id":  pi.SupplierRoomID,
			"date":              pi.Date,
			"price":             pi.Price,
			"original_price":    originalPrice,
			"available_count":   pi.AvailableCount,
			"total_count":       totalCount,
			"status":            "active",
			"updated_at":        time.Now().Format("2006-01-02 15:04:05"),
		}

		pipe.HSet(ctx, key, data)
		pipe.Expire(ctx, key, PriceInventoryTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func GetPriceInventoryCount(supplierID int) (int64, error) {
	pattern := fmt.Sprintf("%s:%d:*", PriceInventoryPrefix, supplierID)
	var count int64
	var cursor uint64

	ctx := context.Background()
	for {
		var keys []string
		var err error
		keys, cursor, err = RDB.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return 0, err
		}
		count += int64(len(keys))
		if cursor == 0 {
			break
		}
	}

	return count, nil
}

func GetPriceInventoryDateRange(supplierID int) (string, string, error) {
	pattern := fmt.Sprintf("%s:%d:*", PriceInventoryPrefix, supplierID)
	var minDate, maxDate string
	var cursor uint64

	ctx := context.Background()
	for {
		var keys []string
		var err error
		keys, cursor, err = RDB.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return "", "", err
		}

		for _, key := range keys {
			parts := splitKey(key)
			if len(parts) >= 5 {
				date := parts[4]
				if minDate == "" || date < minDate {
					minDate = date
				}
				if maxDate == "" || date > maxDate {
					maxDate = date
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return minDate, maxDate, nil
}

func GetPriceInventoryLastUpdate(supplierID int) (string, error) {
	pattern := fmt.Sprintf("%s:%d:*", PriceInventoryPrefix, supplierID)
	var lastUpdate string
	var cursor uint64

	ctx := context.Background()
	for {
		var keys []string
		var err error
		keys, cursor, err = RDB.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return "", err
		}

		for _, key := range keys {
			updatedAt, err := RDB.HGet(ctx, key, "updated_at").Result()
			if err == nil && updatedAt != "" {
				if lastUpdate == "" || updatedAt > lastUpdate {
					lastUpdate = updatedAt
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return lastUpdate, nil
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

func SetPriceInventoryJSON(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

	originalPrice := pi.OriginalPrice
	if originalPrice == 0 {
		originalPrice = pi.Price
	}

	totalCount := pi.TotalCount
	if totalCount == 0 {
		totalCount = 10
	}

	inv := models.PriceInventory{
		SupplierID:       supplierID,
		SupplierHotelID:  pi.SupplierHotelID,
		SupplierRoomID:   pi.SupplierRoomID,
		Date:             pi.Date,
		Price:            pi.Price,
		OriginalPrice:    originalPrice,
		AvailableCount:   pi.AvailableCount,
		TotalCount:       totalCount,
		Status:           "active",
		UpdatedAt:        time.Now(),
	}

	data, err := json.Marshal(inv)
	if err != nil {
		return err
	}

	return RDB.Set(ctx, key, data, PriceInventoryTTL).Err()
}

func GetPriceInventoryJSON(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)

	data, err := RDB.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var pi models.PriceInventory
	err = json.Unmarshal(data, &pi)
	if err != nil {
		return nil, err
	}

	return &pi, nil
}

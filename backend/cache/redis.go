package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/suppliers"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

const (
	PriceInventoryPrefix   = "price_inv"
	PriceInventoryIndex    = "price_inv_index"
	PriceInventoryTTL      = 30 * 24 * time.Hour
	RoomCurrentPricePrefix = "room_price"
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

func IsRedisAvailable() bool {
	return RDB != nil
}

func buildPriceInventoryKey(supplierID int, supplierHotelID, supplierRoomID, date string) string {
	return fmt.Sprintf("%s:%d:%s:%s:%s", PriceInventoryPrefix, supplierID, supplierHotelID, supplierRoomID, date)
}

func buildRoomPriceKey(supplierID int, supplierHotelID, supplierRoomID string) string {
	return fmt.Sprintf("%s:%d:%s:%s", RoomCurrentPricePrefix, supplierID, supplierHotelID, supplierRoomID)
}

func buildSupplierIndexKey(supplierID int) string {
	return fmt.Sprintf("%s:supplier:%d", PriceInventoryIndex, supplierID)
}

func buildHotelIndexKey(supplierID int, supplierHotelID string) string {
	return fmt.Sprintf("%s:hotel:%d:%s", PriceInventoryIndex, supplierID, supplierHotelID)
}

type PriceInventoryData struct {
	SupplierID       int       `json:"supplier_id"`
	SupplierHotelID  string    `json:"supplier_hotel_id"`
	SupplierRoomID   string    `json:"supplier_room_id"`
	Date             string    `json:"date"`
	Price            float64   `json:"price"`
	OriginalPrice    float64   `json:"original_price"`
	AvailableCount   int       `json:"available_count"`
	TotalCount       int       `json:"total_count"`
	Status           string    `json:"status"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewPriceInventoryData(supplierID int, pi models.QiuguoPushPriceInventoryData) *PriceInventoryData {
	originalPrice := pi.OriginalPrice
	if originalPrice == 0 {
		originalPrice = pi.Price
	}

	totalCount := pi.TotalCount
	if totalCount == 0 {
		totalCount = 10
	}

	return &PriceInventoryData{
		SupplierID:      supplierID,
		SupplierHotelID: pi.SupplierHotelID,
		SupplierRoomID:  pi.SupplierRoomID,
		Date:            pi.Date,
		Price:           pi.Price,
		OriginalPrice:   originalPrice,
		AvailableCount:  pi.AvailableCount,
		TotalCount:      totalCount,
		Status:          "active",
		UpdatedAt:       time.Now(),
	}
}

func NewPriceInventoryDataFromSupplier(supplierID int, pi suppliers.SupplierPriceInventoryData) *PriceInventoryData {
	return &PriceInventoryData{
		SupplierID:      supplierID,
		SupplierHotelID: pi.SupplierHotelID,
		SupplierRoomID:  pi.SupplierRoomID,
		Date:            pi.Date,
		Price:           pi.Price,
		OriginalPrice:   pi.Price,
		AvailableCount:  pi.AvailableCount,
		TotalCount:      10,
		Status:          "active",
		UpdatedAt:       time.Now(),
	}
}

func (p *PriceInventoryData) ToHash() map[string]interface{} {
	return map[string]interface{}{
		"supplier_id":       p.SupplierID,
		"supplier_hotel_id": p.SupplierHotelID,
		"supplier_room_id":  p.SupplierRoomID,
		"date":              p.Date,
		"price":             p.Price,
		"original_price":    p.OriginalPrice,
		"available_count":   p.AvailableCount,
		"total_count":       p.TotalCount,
		"status":            p.Status,
		"updated_at":        p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (p *PriceInventoryData) ToModel() *models.PriceInventory {
	return &models.PriceInventory{
		SupplierID:       p.SupplierID,
		SupplierHotelID:  p.SupplierHotelID,
		SupplierRoomID:   p.SupplierRoomID,
		Date:             p.Date,
		Price:            p.Price,
		OriginalPrice:    p.OriginalPrice,
		AvailableCount:   p.AvailableCount,
		TotalCount:       p.TotalCount,
		Status:           p.Status,
		UpdatedAt:        p.UpdatedAt,
	}
}

func parsePriceInventoryFromHash(data map[string]string) *PriceInventoryData {
	if len(data) == 0 {
		return nil
	}

	pi := &PriceInventoryData{}
	pi.SupplierHotelID = data["supplier_hotel_id"]
	pi.SupplierRoomID = data["supplier_room_id"]
	pi.Date = data["date"]
	pi.Status = data["status"]

	if supplierIDStr, ok := data["supplier_id"]; ok {
		if sid, err := strconv.Atoi(supplierIDStr); err == nil {
			pi.SupplierID = sid
		}
	}

	if priceStr, ok := data["price"]; ok {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			pi.Price = price
		}
	}

	if originalPriceStr, ok := data["original_price"]; ok {
		if originalPrice, err := strconv.ParseFloat(originalPriceStr, 64); err == nil {
			pi.OriginalPrice = originalPrice
		}
	}

	if availableCountStr, ok := data["available_count"]; ok {
		if availableCount, err := strconv.Atoi(availableCountStr); err == nil {
			pi.AvailableCount = availableCount
		}
	}

	if totalCountStr, ok := data["total_count"]; ok {
		if totalCount, err := strconv.Atoi(totalCountStr); err == nil {
			pi.TotalCount = totalCount
		}
	}

	if updatedAtStr, ok := data["updated_at"]; ok {
		if updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			pi.UpdatedAt = updatedAt
		}
	}

	return pi
}

func SetPriceInventoryData(pi *PriceInventoryData) error {
	if !IsRedisAvailable() {
		return nil
	}

	key := buildPriceInventoryKey(pi.SupplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

	err := RDB.HSet(ctx, key, pi.ToHash()).Err()
	if err != nil {
		return err
	}

	err = RDB.Expire(ctx, key, PriceInventoryTTL).Err()
	if err != nil {
		log.Printf("设置Redis过期时间失败: %v", err)
	}

	supplierIndexKey := buildSupplierIndexKey(pi.SupplierID)
	hotelIndexKey := buildHotelIndexKey(pi.SupplierID, pi.SupplierHotelID)

	score := float64(time.Now().Unix())
	RDB.ZAdd(ctx, supplierIndexKey, &redis.Z{Score: score, Member: key})
	RDB.ZAdd(ctx, hotelIndexKey, &redis.Z{Score: score, Member: key})

	RDB.Expire(ctx, supplierIndexKey, PriceInventoryTTL)
	RDB.Expire(ctx, hotelIndexKey, PriceInventoryTTL)

	return nil
}

func GetPriceInventoryData(supplierID int, supplierHotelID, supplierRoomID, date string) (*PriceInventoryData, error) {
	if !IsRedisAvailable() {
		return nil, redis.Nil
	}

	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)

	data, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	pi := parsePriceInventoryFromHash(data)
	if pi == nil {
		return nil, redis.Nil
	}

	return pi, nil
}

func GetPriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	pi, err := GetPriceInventoryData(supplierID, supplierHotelID, supplierRoomID, date)
	if err != nil {
		return nil, err
	}
	return pi.ToModel(), nil
}

func SetPriceInventory(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	data := NewPriceInventoryData(supplierID, pi)
	return SetPriceInventoryData(data)
}

func SetSupplierPriceInventory(supplierID int, pi suppliers.SupplierPriceInventoryData) error {
	data := NewPriceInventoryDataFromSupplier(supplierID, pi)
	return SetPriceInventoryData(data)
}

func DeletePriceInventory(supplierID int, supplierHotelID, supplierRoomID, date string) error {
	if !IsRedisAvailable() {
		return nil
	}

	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)

	supplierIndexKey := buildSupplierIndexKey(supplierID)
	hotelIndexKey := buildHotelIndexKey(supplierID, supplierHotelID)

	RDB.ZRem(ctx, supplierIndexKey, key)
	RDB.ZRem(ctx, hotelIndexKey, key)

	return RDB.Del(ctx, key).Err()
}

func BatchSetPriceInventory(supplierID int, priceInventories []models.QiuguoPushPriceInventoryData) error {
	if !IsRedisAvailable() {
		return nil
	}

	pipe := RDB.Pipeline()

	for _, pi := range priceInventories {
		data := NewPriceInventoryData(supplierID, pi)
		key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

		pipe.HSet(ctx, key, data.ToHash())
		pipe.Expire(ctx, key, PriceInventoryTTL)

		supplierIndexKey := buildSupplierIndexKey(supplierID)
		hotelIndexKey := buildHotelIndexKey(supplierID, pi.SupplierHotelID)
		score := float64(time.Now().Unix())

		pipe.ZAdd(ctx, supplierIndexKey, &redis.Z{Score: score, Member: key})
		pipe.ZAdd(ctx, hotelIndexKey, &redis.Z{Score: score, Member: key})
	}

	_, err := pipe.Exec(ctx)
	return err
}

func BatchSetSupplierPriceInventory(supplierID int, priceInventories []suppliers.SupplierPriceInventoryData) error {
	if !IsRedisAvailable() {
		return nil
	}

	pipe := RDB.Pipeline()

	for _, pi := range priceInventories {
		data := NewPriceInventoryDataFromSupplier(supplierID, pi)
		key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

		pipe.HSet(ctx, key, data.ToHash())
		pipe.Expire(ctx, key, PriceInventoryTTL)

		supplierIndexKey := buildSupplierIndexKey(supplierID)
		hotelIndexKey := buildHotelIndexKey(supplierID, pi.SupplierHotelID)
		score := float64(time.Now().Unix())

		pipe.ZAdd(ctx, supplierIndexKey, &redis.Z{Score: score, Member: key})
		pipe.ZAdd(ctx, hotelIndexKey, &redis.Z{Score: score, Member: key})
	}

	_, err := pipe.Exec(ctx)
	return err
}

func GetPriceInventoryCount(supplierID int) (int64, error) {
	if !IsRedisAvailable() {
		return 0, nil
	}

	indexKey := buildSupplierIndexKey(supplierID)
	return RDB.ZCard(ctx, indexKey).Result()
}

func GetPriceInventoryDateRange(supplierID int) (string, string, error) {
	if !IsRedisAvailable() {
		return "", "", nil
	}

	indexKey := buildSupplierIndexKey(supplierID)
	keys, err := RDB.ZRange(ctx, indexKey, 0, -1).Result()
	if err != nil {
		return "", "", err
	}

	if len(keys) == 0 {
		return "", "", nil
	}

	var minDate, maxDate string
	for _, key := range keys {
		parts := strings.Split(key, ":")
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

	return minDate, maxDate, nil
}

func GetPriceInventoryLastUpdate(supplierID int) (string, error) {
	if !IsRedisAvailable() {
		return "", nil
	}

	indexKey := buildSupplierIndexKey(supplierID)
	keys, err := RDB.ZRevRange(ctx, indexKey, 0, 100).Result()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	var lastUpdate string
	for _, key := range keys {
		updatedAt, err := RDB.HGet(ctx, key, "updated_at").Result()
		if err == nil && updatedAt != "" {
			if lastUpdate == "" || updatedAt > lastUpdate {
				lastUpdate = updatedAt
			}
		}
	}

	return lastUpdate, nil
}

func GetPriceInventoryByDateRange(supplierID int, supplierHotelID, supplierRoomID string, startDate, endDate string) ([]*PriceInventoryData, error) {
	if !IsRedisAvailable() {
		return nil, nil
	}

	var result []*PriceInventoryData

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		pi, err := GetPriceInventoryData(supplierID, supplierHotelID, supplierRoomID, date)
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Printf("获取价格库存失败: %v", err)
			continue
		}
		result = append(result, pi)
	}

	return result, nil
}

func GetHotelPriceInventory(supplierID int, supplierHotelID string, date string) ([]*PriceInventoryData, error) {
	if !IsRedisAvailable() {
		return nil, nil
	}

	hotelIndexKey := buildHotelIndexKey(supplierID, supplierHotelID)
	keys, err := RDB.ZRange(ctx, hotelIndexKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []*PriceInventoryData
	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) >= 5 && parts[4] == date {
			data, err := RDB.HGetAll(ctx, key).Result()
			if err == nil {
				pi := parsePriceInventoryFromHash(data)
				if pi != nil {
					result = append(result, pi)
				}
			}
		}
	}

	return result, nil
}

type RoomCurrentPrice struct {
	SupplierID       int       `json:"supplier_id"`
	SupplierHotelID  string    `json:"supplier_hotel_id"`
	SupplierRoomID   string    `json:"supplier_room_id"`
	Price            float64   `json:"price"`
	AvailableCount   int       `json:"available_count"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func SetRoomCurrentPrice(supplierID int, supplierHotelID, supplierRoomID string, price float64, availableCount int) error {
	if !IsRedisAvailable() {
		return nil
	}

	key := buildRoomPriceKey(supplierID, supplierHotelID, supplierRoomID)

	data := map[string]interface{}{
		"supplier_id":       supplierID,
		"supplier_hotel_id": supplierHotelID,
		"supplier_room_id":  supplierRoomID,
		"price":             price,
		"available_count":   availableCount,
		"updated_at":        time.Now().Format("2006-01-02 15:04:05"),
	}

	err := RDB.HSet(ctx, key, data).Err()
	if err != nil {
		return err
	}

	return RDB.Expire(ctx, key, PriceInventoryTTL).Err()
}

func GetRoomCurrentPrice(supplierID int, supplierHotelID, supplierRoomID string) (*RoomCurrentPrice, error) {
	if !IsRedisAvailable() {
		return nil, redis.Nil
	}

	key := buildRoomPriceKey(supplierID, supplierHotelID, supplierRoomID)

	data, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, redis.Nil
	}

	rcp := &RoomCurrentPrice{}
	rcp.SupplierHotelID = data["supplier_hotel_id"]
	rcp.SupplierRoomID = data["supplier_room_id"]

	if supplierIDStr, ok := data["supplier_id"]; ok {
		if sid, err := strconv.Atoi(supplierIDStr); err == nil {
			rcp.SupplierID = sid
		}
	}

	if priceStr, ok := data["price"]; ok {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			rcp.Price = price
		}
	}

	if availableCountStr, ok := data["available_count"]; ok {
		if availableCount, err := strconv.Atoi(availableCountStr); err == nil {
			rcp.AvailableCount = availableCount
		}
	}

	if updatedAtStr, ok := data["updated_at"]; ok {
		if updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			rcp.UpdatedAt = updatedAt
		}
	}

	return rcp, nil
}

func SetPriceInventoryJSON(supplierID int, pi models.QiuguoPushPriceInventoryData) error {
	if !IsRedisAvailable() {
		return nil
	}

	data := NewPriceInventoryData(supplierID, pi)
	key := buildPriceInventoryKey(supplierID, pi.SupplierHotelID, pi.SupplierRoomID, pi.Date)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return RDB.Set(ctx, key, jsonData, PriceInventoryTTL).Err()
}

func GetPriceInventoryJSON(supplierID int, supplierHotelID, supplierRoomID, date string) (*models.PriceInventory, error) {
	if !IsRedisAvailable() {
		return nil, redis.Nil
	}

	key := buildPriceInventoryKey(supplierID, supplierHotelID, supplierRoomID, date)

	data, err := RDB.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var pi PriceInventoryData
	err = json.Unmarshal(data, &pi)
	if err != nil {
		return nil, err
	}

	return pi.ToModel(), nil
}

func ClearPriceInventoryBySupplier(supplierID int) error {
	if !IsRedisAvailable() {
		return nil
	}

	indexKey := buildSupplierIndexKey(supplierID)
	keys, err := RDB.ZRange(ctx, indexKey, 0, -1).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		RDB.Del(ctx, keys...)
	}

	RDB.Del(ctx, indexKey)

	pattern := fmt.Sprintf("%s:hotel:%d:*", PriceInventoryIndex, supplierID)
	var cursor uint64
	for {
		var hotelKeys []string
		var err error
		hotelKeys, cursor, err = RDB.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			break
		}
		if len(hotelKeys) > 0 {
			RDB.Del(ctx, hotelKeys...)
		}
		if cursor == 0 {
			break
		}
	}

	return nil
}

func GetAllPriceInventoryKeys(supplierID int) ([]string, error) {
	if !IsRedisAvailable() {
		return nil, nil
	}

	indexKey := buildSupplierIndexKey(supplierID)
	return RDB.ZRange(ctx, indexKey, 0, -1).Result()
}

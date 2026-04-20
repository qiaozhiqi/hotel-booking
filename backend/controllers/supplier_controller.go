package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hotel-booking/database"
	"hotel-booking/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SupplierCodeHuazhu = "huazhu"
)

type MockHotelData struct {
	SupplierHotelID string           `json:"hotel_id"`
	Name            string           `json:"name"`
	City            string           `json:"city"`
	Address         string           `json:"address"`
	Description     string           `json:"description"`
	Rating          float64          `json:"rating"`
	ImageURL        string           `json:"image_url"`
	PriceRange      string           `json:"price_range"`
	Rooms           []MockRoomData   `json:"rooms"`
}

type MockRoomData struct {
	SupplierRoomID  string  `json:"room_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Capacity        int     `json:"capacity"`
	Area            int     `json:"area"`
	BedType         string  `json:"bed_type"`
	Amenities       string  `json:"amenities"`
	ImageURL        string  `json:"image_url"`
	TotalCount      int     `json:"total_count"`
	AvailableCount  int     `json:"available_count"`
}

func GetHuazhuMockHotels() []MockHotelData {
	rand.Seed(time.Now().UnixNano())
	
	cities := []string{"上海", "北京", "广州", "深圳", "杭州", "南京", "成都", "武汉", "西安", "重庆"}
	hotelTypes := []string{"全季酒店", "汉庭酒店", "桔子酒店", "美居酒店", "宜必思酒店", "禧玥酒店", "花间堂", "漫心酒店", "星程酒店", "怡莱酒店"}
	addressSuffixes := []string{"南京路店", "外滩店", "陆家嘴店", "王府井店", "三里屯店", "天河店", "科技园店", "西湖店", "春熙路店", "解放碑店"}
	
	hotels := make([]MockHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 200 + rand.Intn(300)
		maxPrice := minPrice + 200 + rand.Intn(300)
		
		hotels[i] = MockHotelData{
			SupplierHotelID: fmt.Sprintf("HZ-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 100+rand.Intn(900)),
			Description:     fmt.Sprintf("%s位于%s核心商圈，交通便利，设施完善，提供优质的住宿体验。酒店拥有各类客房，配备现代化设施，是商务出行和休闲旅游的理想选择。", hotelType, city),
			Rating:          4.0 + float64(rand.Intn(10))/10,
			ImageURL:        getHotelImageURL(i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateMockRooms(i, minPrice),
		}
	}
	
	return hotels
}

func getRandomDistrict(city string) string {
	districts := map[string][]string{
		"上海": {"黄浦", "徐汇", "长宁", "静安", "普陀", "虹口", "杨浦", "浦东"},
		"北京": {"东城", "西城", "朝阳", "海淀", "丰台", "石景山"},
		"广州": {"天河", "越秀", "海珠", "荔湾", "白云", "番禺"},
		"深圳": {"南山", "福田", "罗湖", "宝安", "龙岗", "龙华"},
		"杭州": {"西湖", "上城", "下城", "江干", "拱墅", "滨江"},
		"南京": {"玄武", "秦淮", "建邺", "鼓楼", "浦口", "栖霞"},
		"成都": {"锦江", "青羊", "金牛", "武侯", "成华", "高新"},
		"武汉": {"江岸", "江汉", "硚口", "汉阳", "武昌", "青山"},
		"西安": {"新城", "碑林", "莲湖", "雁塔", "未央", "灞桥"},
		"重庆": {"渝中", "大渡口", "江北", "沙坪坝", "九龙坡", "南岸"},
	}
	if d, ok := districts[city]; ok {
		return d[rand.Intn(len(d))]
	}
	return "中心"
}

func getRandomRoad() string {
	roads := []string{"中山", "人民", "解放", "建国", "和平", "新华", "建设", "发展", "科技", "商务"}
	return roads[rand.Intn(len(roads))]
}

func getHotelImageURL(index int) string {
	themes := []string{
		"modern luxury hotel lobby interior design",
		"elegant hotel entrance with glass facade",
		"contemporary hotel building night view",
		"luxury hotel suite with city view",
		"boutique hotel lobby with art decor",
		"business hotel exterior modern architecture",
		"resort hotel with swimming pool",
		"city hotel with skyline view",
		"designer hotel interior minimalist style",
		"grand hotel entrance with porte cochere",
		"urban hotel with green rooftop",
		"classic hotel with modern renovation",
	}
	theme := themes[index%len(themes)]
	return fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=%s&image_size=landscape_4_3", theme)
}

func getRoomImageURL(hotelIndex, roomIndex int) string {
	themes := []string{
		"luxury hotel bedroom king size bed",
		"modern hotel room twin beds",
		"hotel suite with living area",
		"deluxe hotel room with city view",
		"executive hotel room work desk",
	}
	theme := themes[roomIndex%len(themes)]
	return fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=%s&image_size=landscape_4_3", theme)
}

func generateMockRooms(hotelIndex int, basePrice int) []MockRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"标准双床房", "经济实惠的标准客房，两张单人床，适合商务出行。配备免费WiFi、空调、电视和24小时热水。", 1.0, 2, 28, "双床"},
		{"标准大床房", "舒适的标准客房，一张大床，适合情侣或单人入住。配备免费WiFi、空调、电视和24小时热水。", 1.1, 2, 30, "大床"},
		{"豪华大床房", "宽敞明亮的豪华客房，配备高品质床品和现代化设施。享有城市美景，是商务和休闲的理想选择。", 1.5, 2, 40, "大床"},
		{"行政套房", "高端行政套房，独立客厅和卧室，配备宽大办公桌和高速网络。享受行政楼层礼遇，适合高端商务人士。", 2.5, 2, 65, "大床"},
		{"家庭房", "温馨舒适的家庭房，配备大床和单人床，适合家庭出行。空间宽敞，设施齐全，让您享受家的温暖。", 1.8, 4, 50, "大床+双床"},
	}
	
	rooms := make([]MockRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 10 + rand.Intn(30)
		availableCount := totalCount - rand.Intn(5)
		
		rooms[i] = MockRoomData{
			SupplierRoomID: fmt.Sprintf("HZ-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           float64(int(price/10) * 10),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 空调, 电视, 24小时热水, 吹风机, 电热水壶",
			ImageURL:        getRoomImageURL(hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

func InitSupplier() error {
	db := database.GetDB()
	
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM suppliers WHERE code = ?", SupplierCodeHuazhu).Scan(&count)
	if err != nil {
		return err
	}
	
	if count == 0 {
		_, err = db.Exec(`
			INSERT INTO suppliers (name, code, description, api_url, status)
			VALUES (?, ?, ?, ?, ?)`,
			"华住酒店集团", SupplierCodeHuazhu, 
			"华住酒店集团是中国领先的多品牌酒店集团，旗下包括全季、汉庭、桔子、美居等多个知名酒店品牌。",
			"/api/mock/huazhu", "active")
		if err != nil {
			return err
		}
	}
	
	return nil
}

func GetSupplierList(c *gin.Context) {
	db := database.GetDB()
	
	rows, err := db.Query(`
		SELECT id, name, code, description, api_url, status, created_at, updated_at
		FROM suppliers ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取供应商列表失败",
		})
		return
	}
	defer rows.Close()
	
	var suppliers []models.Supplier
	for rows.Next() {
		var supplier models.Supplier
		err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Description,
			&supplier.APIURL, &supplier.Status, &supplier.CreatedAt, &supplier.UpdatedAt)
		if err != nil {
			continue
		}
		suppliers = append(suppliers, supplier)
	}
	
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    suppliers,
	})
}

func MockHuazhuGetHotels(c *gin.Context) {
	hotels := GetHuazhuMockHotels()
	
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"hotels": hotels,
			"total":  len(hotels),
		},
	})
}

func MockHuazhuGetHotelDetail(c *gin.Context) {
	hotelID := c.Param("id")
	hotels := GetHuazhuMockHotels()
	
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			c.JSON(http.StatusOK, models.Response{
				Code:    200,
				Message: "获取成功",
				Data:    hotel,
			})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, models.Response{
		Code:    404,
		Message: "酒店不存在",
	})
}

func PullSupplierData(c *gin.Context) {
	supplierCode := c.Param("code")
	
	if supplierCode != SupplierCodeHuazhu {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "不支持的供应商",
		})
		return
	}
	
	db := database.GetDB()
	
	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&supplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = InitSupplier()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{
					Code:    500,
					Message: "初始化供应商失败",
				})
				return
			}
			db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&supplierID)
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "查询供应商失败",
			})
			return
		}
	}
	
	mockHotels := GetHuazhuMockHotels()
	
	pulledCount := 0
	updatedCount := 0
	
	for _, mockHotel := range mockHotels {
		var existingID int
		err := db.QueryRow(`
			SELECT id FROM supplier_hotels WHERE supplier_id = ? AND supplier_hotel_id = ?`,
			supplierID, mockHotel.SupplierHotelID).Scan(&existingID)
		
		rawData, _ := json.Marshal(mockHotel)
		
		if err == sql.ErrNoRows {
			tx, err := db.Begin()
			if err != nil {
				continue
			}
			
			hotelResult, err := tx.Exec(`
				INSERT INTO hotels (name, address, city, description, rating, image_url, price_range, supplier_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				mockHotel.Name, mockHotel.Address, mockHotel.City, mockHotel.Description,
				mockHotel.Rating, mockHotel.ImageURL, mockHotel.PriceRange, supplierID)
			if err != nil {
				tx.Rollback()
				continue
			}
			
			localHotelID, _ := hotelResult.LastInsertId()
			
			_, err = tx.Exec(`
				INSERT INTO supplier_hotels (supplier_id, supplier_hotel_id, local_hotel_id, hotel_name, 
				city, address, rating, image_url, price_range, raw_data, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				supplierID, mockHotel.SupplierHotelID, localHotelID, mockHotel.Name,
				mockHotel.City, mockHotel.Address, mockHotel.Rating, mockHotel.ImageURL,
				mockHotel.PriceRange, string(rawData), "synced")
			if err != nil {
				tx.Rollback()
				continue
			}
			
			for _, mockRoom := range mockHotel.Rooms {
				roomResult, err := tx.Exec(`
					INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, supplier_id)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					localHotelID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
					mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
					mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount, supplierID)
				if err != nil {
					continue
				}
				
				localRoomID, _ := roomResult.LastInsertId()
				
				roomRawData, _ := json.Marshal(mockRoom)
				tx.Exec(`
					INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
					local_room_id, room_name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, raw_data, status)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					supplierID, mockHotel.SupplierHotelID, mockRoom.SupplierRoomID,
					localRoomID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
					mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
					mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount,
					string(roomRawData), "synced")
			}
			
			tx.Commit()
			pulledCount++
		} else if err == nil {
			var localHotelID int
			db.QueryRow("SELECT local_hotel_id FROM supplier_hotels WHERE id = ?", existingID).Scan(&localHotelID)
			
			_, err = db.Exec(`
				UPDATE hotels SET name = ?, address = ?, city = ?, description = ?, 
				rating = ?, image_url = ?, price_range = ? WHERE id = ?`,
				mockHotel.Name, mockHotel.Address, mockHotel.City, mockHotel.Description,
				mockHotel.Rating, mockHotel.ImageURL, mockHotel.PriceRange, localHotelID)
			
			_, err = db.Exec(`
				UPDATE supplier_hotels SET hotel_name = ?, city = ?, address = ?, rating = ?, 
				image_url = ?, price_range = ?, raw_data = ?, status = ? WHERE id = ?`,
				mockHotel.Name, mockHotel.City, mockHotel.Address, mockHotel.Rating,
				mockHotel.ImageURL, mockHotel.PriceRange, string(rawData), "synced", existingID)
			
			for _, mockRoom := range mockHotel.Rooms {
				var existingRoomID int
				var localRoomID int
				err := db.QueryRow(`
					SELECT id, local_room_id FROM supplier_rooms 
					WHERE supplier_id = ? AND supplier_hotel_id = ? AND supplier_room_id = ?`,
					supplierID, mockHotel.SupplierHotelID, mockRoom.SupplierRoomID).Scan(&existingRoomID, &localRoomID)
				
				if err == sql.ErrNoRows {
					roomResult, _ := db.Exec(`
						INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
						amenities, image_url, total_count, available_count, supplier_id)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
						localHotelID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
						mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
						mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount, supplierID)
					
					newLocalRoomID, _ := roomResult.LastInsertId()
					roomRawData, _ := json.Marshal(mockRoom)
					db.Exec(`
						INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
						local_room_id, room_name, description, price, capacity, area, bed_type, 
						amenities, image_url, total_count, available_count, raw_data, status)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
						supplierID, mockHotel.SupplierHotelID, mockRoom.SupplierRoomID,
						newLocalRoomID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
						mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
						mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount,
						string(roomRawData), "synced")
				} else if err == nil {
					db.Exec(`
						UPDATE rooms SET name = ?, description = ?, price = ?, capacity = ?, 
						area = ?, bed_type = ?, amenities = ?, image_url = ?, 
						total_count = ?, available_count = ? WHERE id = ?`,
						mockRoom.Name, mockRoom.Description, mockRoom.Price, mockRoom.Capacity,
						mockRoom.Area, mockRoom.BedType, mockRoom.Amenities, mockRoom.ImageURL,
						mockRoom.TotalCount, mockRoom.AvailableCount, localRoomID)
					
					roomRawData, _ := json.Marshal(mockRoom)
					db.Exec(`
						UPDATE supplier_rooms SET room_name = ?, description = ?, price = ?, 
						capacity = ?, area = ?, bed_type = ?, amenities = ?, image_url = ?, 
						total_count = ?, available_count = ?, raw_data = ?, status = ? WHERE id = ?`,
						mockRoom.Name, mockRoom.Description, mockRoom.Price, mockRoom.Capacity,
						mockRoom.Area, mockRoom.BedType, mockRoom.Amenities, mockRoom.ImageURL,
						mockRoom.TotalCount, mockRoom.AvailableCount, string(roomRawData), "synced", existingRoomID)
				}
			}
			updatedCount++
		}
	}
	
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "数据拉取成功",
		Data: map[string]interface{}{
			"supplier":      "华住酒店集团",
			"pulled_count":  pulledCount,
			"updated_count": updatedCount,
			"total_count":   len(mockHotels),
		},
	})
}

func GetSyncStatus(c *gin.Context) {
	supplierCode := c.Param("code")
	db := database.GetDB()
	
	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&supplierID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "供应商不存在",
		})
		return
	}
	
	var hotelCount int
	db.QueryRow("SELECT COUNT(*) FROM supplier_hotels WHERE supplier_id = ?", supplierID).Scan(&hotelCount)
	
	var roomCount int
	db.QueryRow("SELECT COUNT(*) FROM supplier_rooms WHERE supplier_id = ?", supplierID).Scan(&roomCount)
	
	var syncedHotelCount int
	db.QueryRow("SELECT COUNT(*) FROM supplier_hotels WHERE supplier_id = ? AND status = ?", supplierID, "synced").Scan(&syncedHotelCount)
	
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"supplier_code":      supplierCode,
			"total_hotels":       hotelCount,
			"total_rooms":        roomCount,
			"synced_hotels":      syncedHotelCount,
			"last_sync_time":     time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

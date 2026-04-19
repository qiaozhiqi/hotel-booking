package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type RuJiaAdapter struct{}

func NewRuJiaAdapter() *RuJiaAdapter {
	return &RuJiaAdapter{}
}

func (r *RuJiaAdapter) GetCode() string {
	return "rujia"
}

func (r *RuJiaAdapter) GetName() string {
	return "如家酒店集团"
}

func (r *RuJiaAdapter) GetDescription() string {
	return "如家酒店集团是中国知名的经济型连锁酒店品牌，旗下包括如家、莫泰、和颐等多个品牌，致力于为旅客提供干净、温馨的住宿体验。"
}

func (r *RuJiaAdapter) GetAPIURL() string {
	return "/api/mock/rujia"
}

func (r *RuJiaAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateRuJiaMockData(), nil
}

func (r *RuJiaAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateRuJiaMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateRuJiaMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 1000)
	
	cities := []string{"北京", "上海", "广州", "深圳", "苏州", "无锡", "宁波", "福州", "厦门", "青岛"}
	hotelTypes := []string{"如家酒店", "莫泰酒店", "和颐酒店", "如家精选", "如家商旅", "派酒店", "云上四季"}
	addressSuffixes := []string{"西单店", "国贸店", "中关村店", "陆家嘴店", "人民广场店", "天河店", "科技园店", "观前街店", "西湖店", "中山路店"}
	
	hotels := make([]SupplierHotelData, 10)
	
	for i := 0; i < 10; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 150 + rand.Intn(200)
		maxPrice := minPrice + 150 + rand.Intn(200)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("RJ-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s街%d号", city, getRandomDistrict(city), getRandomRoad(), 50+rand.Intn(500)),
			Description:     fmt.Sprintf("%s秉承如家酒店连锁的特色——三大“统一”性：统一建筑设施；统一服务；统一硬件设施。酒店内外由名家设计，风格简约、别致，设施齐全舒适，展现的是一个“干净、温馨”的住宿环境。", hotelType),
			Rating:          3.8 + float64(rand.Intn(12))/10,
			ImageURL:        getHotelImageURL("rujia", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateRuJiaRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateRuJiaRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"特惠双床房", "经济实惠的特惠客房，适合预算有限的旅客。房间干净整洁，配备基本设施。", 0.9, 2, 22, "双床"},
		{"标准大床房", "舒适的标准客房，一张大床，适合情侣或单人入住。配备免费WiFi、空调、电视。", 1.0, 2, 25, "大床"},
		{"商务大床房", "专为商务旅客设计，配备宽大办公桌和高速网络。房间宽敞明亮，工作休息两不误。", 1.3, 2, 32, "大床"},
		{"精选套房", "宽敞的精选套房，独立起居空间，享受更舒适的住宿体验。适合需要更多空间的旅客。", 2.0, 2, 50, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 15 + rand.Intn(25)
		availableCount := totalCount - rand.Intn(8)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("RJ-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           float64(int(price/10) * 10),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 空调, 电视, 24小时热水, 吹风机",
			ImageURL:        getRoomImageURL("rujia", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

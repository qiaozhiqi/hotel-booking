package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiLvdiAdapter struct{}

func NewShijiLvdiAdapter() *ShijiLvdiAdapter {
	return &ShijiLvdiAdapter{}
}

func (s *ShijiLvdiAdapter) GetCode() string {
	return "shiji_lvdi"
}

func (s *ShijiLvdiAdapter) GetName() string {
	return "绿地酒店集团（石基畅联）"
}

func (s *ShijiLvdiAdapter) GetDescription() string {
	return "绿地酒店集团是中国知名的酒店管理集团，旗下拥有绿地铂瑞、绿地铂骊、绿地假日等多个酒店品牌，致力于为宾客提供高品质的住宿体验。通过石基畅联渠道接入。"
}

func (s *ShijiLvdiAdapter) GetAPIURL() string {
	return "/api/mock/shiji_lvdi"
}

func (s *ShijiLvdiAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiLvdiMockData(), nil
}

func (s *ShijiLvdiAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiLvdiMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiLvdiMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3500)
	
	cities := []string{"上海", "北京", "广州", "深圳", "南京", "杭州", "成都", "武汉", "西安", "重庆", "苏州", "无锡"}
	hotelTypes := []string{"绿地铂瑞酒店", "绿地铂骊酒店", "绿地假日酒店", "绿地康养居酒店", "绿地魔奇酒店"}
	addressSuffixes := []string{"黄浦店", "朝阳店", "天河店", "龙华店", "鼓楼店", "西湖店", "武侯店", "武昌店", "雁塔店", "江北店", "工业园店", "滨湖店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 350 + rand.Intn(300)
		maxPrice := minPrice + 300 + rand.Intn(300)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-LD-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 90+rand.Intn(550)),
			Description:     fmt.Sprintf("%s是绿地酒店集团旗下的知名酒店品牌，秉承\"让生活更美好\"的企业理念，为宾客提供舒适、便捷、高品质的住宿体验。酒店设施完善，服务周到，是商务出行和休闲旅游的理想选择。", hotelType),
			Rating:          roundToOneDecimal(4.1 + float64(rand.Intn(9))/10),
			ImageURL:        getHotelImageURL("shiji_lvdi", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiLvdiRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiLvdiRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"标准双床房", "舒适的标准双床房，配备优质床品和现代化设施。房间干净整洁，是商务出行的经济实惠选择。", 1.0, 2, 24, "双床"},
		{"豪华大床房", "宽敞的豪华大床房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.15, 2, 30, "大床"},
		{"行政客房", "位于行政楼层的高级客房，可享受行政服务。配备办公区域，专为商务旅客设计。", 1.4, 2, 36, "大床/双床"},
		{"豪华套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具，享受尊贵住宿体验。", 1.8, 2, 48, "大床"},
		{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 3.5, 4, 100, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 18 + rand.Intn(15)
		availableCount := totalCount - rand.Intn(8)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-LD-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋",
			ImageURL:        getRoomImageURL("shiji_lvdi", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiMarriottAdapter struct{}

func NewShijiMarriottAdapter() *ShijiMarriottAdapter {
	return &ShijiMarriottAdapter{}
}

func (s *ShijiMarriottAdapter) GetCode() string {
	return "shiji_marriott"
}

func (s *ShijiMarriottAdapter) GetName() string {
	return "万豪国际酒店集团（石基畅联）"
}

func (s *ShijiMarriottAdapter) GetDescription() string {
	return "万豪国际集团是全球首屈一指的酒店管理公司，旗下拥有万豪、喜来登、丽思卡尔顿等多个知名酒店品牌。通过石基畅联渠道接入。"
}

func (s *ShijiMarriottAdapter) GetAPIURL() string {
	return "/api/mock/shiji_marriott"
}

func (s *ShijiMarriottAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiMarriottMockData(), nil
}

func (s *ShijiMarriottAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiMarriottMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiMarriottMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3000)
	
	cities := []string{"上海", "北京", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "苏州", "三亚"}
	hotelTypes := []string{"万豪酒店", "喜来登酒店", "威斯汀酒店", "丽思卡尔顿酒店", "JW万豪酒店", "万怡酒店", "万丽酒店", "福朋喜来登"}
	addressSuffixes := []string{"陆家嘴店", "国贸店", "天河店", "福田店", "西湖店", "春熙路店", "解放碑店", "光谷店", "大雁塔店", "玄武湖店", "金鸡湖店", "亚龙湾店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 500 + rand.Intn(500)
		maxPrice := minPrice + 500 + rand.Intn(500)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-MA-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 100+rand.Intn(900)),
			Description:     fmt.Sprintf("%s是万豪国际集团旗下的高端酒店品牌，秉承卓越服务理念，为宾客提供无与伦比的住宿体验。酒店地理位置优越，设施豪华完善，是商务旅行和休闲度假的理想之选。", hotelType),
			Rating:          roundToOneDecimal(4.5 + float64(rand.Intn(5))/10),
			ImageURL:        getHotelImageURL("shiji_marriott", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiMarriottRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiMarriottRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"豪华客房", "宽敞明亮的豪华客房，配备高品质床品和现代化设施。享受城市美景，是商务和休闲的理想选择。", 1.0, 2, 35, "大床/双床"},
		{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊礼遇。配备宽大办公桌和高速网络，专为商务精英打造。", 1.4, 2, 42, "大床/双床"},
		{"小型套房", "独立客厅与卧室的小型套房，空间宽敞舒适。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.0, 2, 55, "大床"},
		{"行政套房", "高端行政套房，独立客厅和卧室，配备宽大办公桌和高速网络。享受行政楼层礼遇，适合高端商务人士。", 2.8, 2, 75, "大床"},
		{"豪华套房", "顶级豪华套房，宽敞的独立客厅和卧室，配备奢华家具和设施。享受私人管家服务，是尊贵宾客的首选。", 4.0, 3, 120, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 8 + rand.Intn(25)
		availableCount := totalCount - rand.Intn(4)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-MA-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费矿泉水",
			ImageURL:        getRoomImageURL("shiji_marriott", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

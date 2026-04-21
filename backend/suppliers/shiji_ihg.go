package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiIHGAdapter struct{}

func NewShijiIHGAdapter() *ShijiIHGAdapter {
	return &ShijiIHGAdapter{}
}

func (s *ShijiIHGAdapter) GetCode() string {
	return "shiji_ihg"
}

func (s *ShijiIHGAdapter) GetName() string {
	return "洲际酒店集团（石基畅联）"
}

func (s *ShijiIHGAdapter) GetDescription() string {
	return "洲际酒店集团（IHG）是全球最大的酒店管理公司之一，旗下拥有洲际、皇冠假日、假日酒店等多个知名品牌。通过石基畅联渠道接入。"
}

func (s *ShijiIHGAdapter) GetAPIURL() string {
	return "/api/mock/shiji_ihg"
}

func (s *ShijiIHGAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiIHGMockData(), nil
}

func (s *ShijiIHGAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiIHGMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiIHGMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3200)
	
	cities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "天津", "大连"}
	hotelTypes := []string{"洲际酒店", "皇冠假日酒店", "假日酒店", "智选假日酒店", "英迪格酒店", "VOCO酒店", "逸衡酒店"}
	addressSuffixes := []string{"三里屯店", "静安店", "白云店", "南山店", "滨江店", "锦江区店", "渝中区店", "武昌店", "高新区店", "鼓楼店", "和平店", "中山区店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 450 + rand.Intn(400)
		maxPrice := minPrice + 400 + rand.Intn(400)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-IHG-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 60+rand.Intn(700)),
			Description:     fmt.Sprintf("%s是洲际酒店集团旗下的知名酒店品牌，秉承宾至如归的服务理念，为宾客提供优质的住宿体验。酒店设施完善，服务周到，是商务出行和休闲旅游的理想选择。", hotelType),
			Rating:          roundToOneDecimal(4.3 + float64(rand.Intn(7))/10),
			ImageURL:        getHotelImageURL("shiji_ihg", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiIHGRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiIHGRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"标准房", "舒适的标准客房，配备优质床品和基本设施。干净整洁，是商务出行的经济实惠选择。", 1.0, 2, 28, "大床/双床"},
		{"豪华房", "宽敞的豪华客房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.25, 2, 36, "大床/双床"},
		{"行政房", "位于行政楼层的高级客房，可享受行政俱乐部服务。配备办公区域，专为商务旅客设计。", 1.5, 2, 42, "大床/双床"},
		{"套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具，享受尊贵住宿体验。", 2.0, 2, 58, "大床"},
		{"洲际套房", "高端洲际套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 3.5, 3, 100, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 12 + rand.Intn(20)
		availableCount := totalCount - rand.Intn(6)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-IHG-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋",
			ImageURL:        getRoomImageURL("shiji_ihg", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

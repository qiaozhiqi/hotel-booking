package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiHiltonAdapter struct{}

func NewShijiHiltonAdapter() *ShijiHiltonAdapter {
	return &ShijiHiltonAdapter{}
}

func (s *ShijiHiltonAdapter) GetCode() string {
	return "shiji_hilton"
}

func (s *ShijiHiltonAdapter) GetName() string {
	return "希尔顿酒店集团（石基畅联）"
}

func (s *ShijiHiltonAdapter) GetDescription() string {
	return "希尔顿全球控股有限公司是全球最知名的酒店管理公司之一，旗下拥有希尔顿、华尔道夫、康莱德等多个奢华及高端酒店品牌。通过石基畅联渠道接入。"
}

func (s *ShijiHiltonAdapter) GetAPIURL() string {
	return "/api/mock/shiji_hilton"
}

func (s *ShijiHiltonAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiHiltonMockData(), nil
}

func (s *ShijiHiltonAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiHiltonMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiHiltonMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3100)
	
	cities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "厦门", "青岛"}
	hotelTypes := []string{"希尔顿酒店", "华尔道夫酒店", "康莱德酒店", "希尔顿逸林酒店", "希尔顿花园酒店", "欢朋酒店", "希尔顿嘉悦里"}
	addressSuffixes := []string{"王府井店", "外滩店", "珠江新城店", "华侨城店", "钱江新城店", "高新店", "江北嘴店", "江汉路店", "曲江店", "河西店", "环岛路店", "五四广场店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 480 + rand.Intn(450)
		maxPrice := minPrice + 450 + rand.Intn(450)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-HI-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 80+rand.Intn(800)),
			Description:     fmt.Sprintf("%s是希尔顿全球旗下的知名酒店品牌，以卓越的服务品质和舒适的住宿环境著称。酒店设计时尚典雅，设施先进完善，为宾客提供难忘的入住体验。", hotelType),
			Rating:          roundToOneDecimal(4.4 + float64(rand.Intn(6))/10),
			ImageURL:        getHotelImageURL("shiji_hilton", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiHiltonRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiHiltonRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"客房", "舒适温馨的标准客房，配备希尔顿特色睡床和现代化设施。享受优质睡眠体验，是商务出行的理想选择。", 1.0, 2, 32, "大床/双床"},
		{"豪华客房", "宽敞的豪华客房，配备高品质家具和现代化设施。享有城市美景，提供更舒适的住宿体验。", 1.3, 2, 40, "大床/双床"},
		{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊服务。配备办公区域和高速网络，专为商务旅客设计。", 1.6, 2, 45, "大床/双床"},
		{"套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.2, 2, 60, "大床"},
		{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和私人管家服务，是尊贵宾客的首选。", 5.0, 4, 150, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 10 + rand.Intn(22)
		availableCount := totalCount - rand.Intn(5)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-HI-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 希尔顿睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费咖啡茶包",
			ImageURL:        getRoomImageURL("shiji_hilton", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

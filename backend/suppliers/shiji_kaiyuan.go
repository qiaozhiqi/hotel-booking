package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiKaiyuanAdapter struct{}

func NewShijiKaiyuanAdapter() *ShijiKaiyuanAdapter {
	return &ShijiKaiyuanAdapter{}
}

func (s *ShijiKaiyuanAdapter) GetCode() string {
	return "shiji_kaiyuan"
}

func (s *ShijiKaiyuanAdapter) GetName() string {
	return "开元酒店集团（石基畅联）"
}

func (s *ShijiKaiyuanAdapter) GetDescription() string {
	return "开元酒店集团是中国最大的民营高星级连锁酒店集团，旗下拥有开元名都、开元大酒店、开元度假村等多个知名品牌。通过石基畅联渠道接入。"
}

func (s *ShijiKaiyuanAdapter) GetAPIURL() string {
	return "/api/mock/shiji_kaiyuan"
}

func (s *ShijiKaiyuanAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiKaiyuanMockData(), nil
}

func (s *ShijiKaiyuanAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiKaiyuanMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiKaiyuanMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3300)
	
	cities := []string{"杭州", "上海", "北京", "宁波", "温州", "绍兴", "嘉兴", "湖州", "金华", "台州", "南京", "苏州"}
	hotelTypes := []string{"开元名都大酒店", "开元大酒店", "开元度假村", "开元名庭酒店", "开元曼居酒店", "开元颐居酒店", "开元森泊度假酒店"}
	addressSuffixes := []string{"西湖店", "松江店", "朝阳店", "鄞州店", "鹿城店", "越城店", "南湖店", "吴兴店", "婺城店", "椒江店", "江宁店", "姑苏店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 380 + rand.Intn(350)
		maxPrice := minPrice + 350 + rand.Intn(350)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-KY-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 50+rand.Intn(600)),
			Description:     fmt.Sprintf("%s是开元酒店集团旗下的知名酒店品牌，秉承\"东方文化与国际标准完美融合\"的服务理念，为宾客提供高品质的住宿体验。酒店环境优雅，设施完善，服务周到。", hotelType),
			Rating:          roundToOneDecimal(4.2 + float64(rand.Intn(8))/10),
			ImageURL:        getHotelImageURL("shiji_kaiyuan", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiKaiyuanRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiKaiyuanRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"标准双床房", "舒适的标准双床房，配备开元特色睡床和现代化设施。房间干净整洁，是商务出行的理想选择。", 1.0, 2, 26, "双床"},
		{"豪华大床房", "宽敞的豪华大床房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.2, 2, 32, "大床"},
		{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊服务。配备办公区域，专为商务旅客设计。", 1.5, 2, 38, "大床/双床"},
		{"豪华套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.0, 2, 52, "大床"},
		{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 4.0, 4, 120, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 15 + rand.Intn(18)
		availableCount := totalCount - rand.Intn(7)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-KY-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 开元睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费茶包",
			ImageURL:        getRoomImageURL("shiji_kaiyuan", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

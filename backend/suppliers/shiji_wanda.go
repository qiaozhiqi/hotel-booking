package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiWandaAdapter struct{}

func NewShijiWandaAdapter() *ShijiWandaAdapter {
	return &ShijiWandaAdapter{}
}

func (s *ShijiWandaAdapter) GetCode() string {
	return "shiji_wanda"
}

func (s *ShijiWandaAdapter) GetName() string {
	return "万达酒店集团（石基畅联）"
}

func (s *ShijiWandaAdapter) GetDescription() string {
	return "万达酒店集团是中国知名的高端酒店管理集团，旗下拥有万达文华、万达嘉华、万达瑞华等多个奢华及高端酒店品牌。通过石基畅联渠道接入。"
}

func (s *ShijiWandaAdapter) GetAPIURL() string {
	return "/api/mock/shiji_wanda"
}

func (s *ShijiWandaAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateShijiWandaMockData(), nil
}

func (s *ShijiWandaAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateShijiWandaMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateShijiWandaMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 3400)
	
	cities := []string{"北京", "上海", "广州", "深圳", "成都", "重庆", "武汉", "西安", "南京", "杭州", "青岛", "大连"}
	hotelTypes := []string{"万达文华酒店", "万达嘉华酒店", "万达瑞华酒店", "万达锦华酒店", "万达美华酒店", "万达颐华酒店", "万达悦华酒店"}
	addressSuffixes := []string{"CBD店", "五角场店", "海珠店", "宝安店", "金牛店", "南岸店", "洪山店", "未央店", "建邺店", "拱墅店", "市南店", "西岗店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 420 + rand.Intn(380)
		maxPrice := minPrice + 380 + rand.Intn(380)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("SJ-WD-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 70+rand.Intn(650)),
			Description:     fmt.Sprintf("%s是万达酒店集团旗下的高端酒店品牌，秉承\"华仕生活、从容尊享\"的服务理念，为宾客提供卓越的住宿体验。酒店设计奢华典雅，设施完善，服务贴心周到。", hotelType),
			Rating:          roundToOneDecimal(4.3 + float64(rand.Intn(7))/10),
			ImageURL:        getHotelImageURL("shiji_wanda", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateShijiWandaRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateShijiWandaRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"豪华大床房", "宽敞舒适的豪华大床房，配备万达特色睡床和现代化设施。房间设计典雅，是商务出行的理想选择。", 1.0, 2, 30, "大床"},
		{"行政套房", "位于行政楼层的高级套房，独立客厅与卧室。可享受行政酒廊服务，专为商务精英打造。", 1.8, 2, 48, "大床"},
		{"豪华套房", "宽敞的豪华套房，独立客厅与卧室，配备高品质家具。享受城市美景，提供尊贵的住宿体验。", 2.5, 2, 65, "大床"},
		{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室、餐厅和书房。配备奢华家具和私人管家服务，是尊贵宾客的首选。", 5.0, 4, 180, "大床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 10 + rand.Intn(20)
		availableCount := totalCount - rand.Intn(5)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("SJ-WD-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 中央空调, 万达睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费矿泉水",
			ImageURL:        getRoomImageURL("shiji_wanda", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

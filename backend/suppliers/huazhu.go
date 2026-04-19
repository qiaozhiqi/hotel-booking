package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type HuazhuAdapter struct{}

func NewHuazhuAdapter() *HuazhuAdapter {
	return &HuazhuAdapter{}
}

func (h *HuazhuAdapter) GetCode() string {
	return "huazhu"
}

func (h *HuazhuAdapter) GetName() string {
	return "华住酒店集团"
}

func (h *HuazhuAdapter) GetDescription() string {
	return "华住酒店集团是中国领先的多品牌酒店集团，旗下包括全季、汉庭、桔子、美居等多个知名酒店品牌。"
}

func (h *HuazhuAdapter) GetAPIURL() string {
	return "/api/mock/huazhu"
}

func (h *HuazhuAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateHuazhuMockData(), nil
}

func (h *HuazhuAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateHuazhuMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateHuazhuMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano())
	
	cities := []string{"上海", "北京", "广州", "深圳", "杭州", "南京", "成都", "武汉", "西安", "重庆"}
	hotelTypes := []string{"全季酒店", "汉庭酒店", "桔子酒店", "美居酒店", "宜必思酒店", "禧玥酒店", "花间堂", "漫心酒店", "星程酒店", "怡莱酒店"}
	addressSuffixes := []string{"南京路店", "外滩店", "陆家嘴店", "王府井店", "三里屯店", "天河店", "科技园店", "西湖店", "春熙路店", "解放碑店"}
	
	hotels := make([]SupplierHotelData, 12)
	
	for i := 0; i < 12; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 200 + rand.Intn(300)
		maxPrice := minPrice + 200 + rand.Intn(300)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("HZ-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 100+rand.Intn(900)),
			Description:     fmt.Sprintf("%s位于%s核心商圈，交通便利，设施完善，提供优质的住宿体验。酒店拥有各类客房，配备现代化设施，是商务出行和休闲旅游的理想选择。", hotelType, city),
			Rating:          roundToOneDecimal(4.0 + float64(rand.Intn(10))/10),
			ImageURL:        getHotelImageURL("huazhu", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateHuazhuRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateHuazhuRooms(hotelIndex int, basePrice int) []SupplierRoomData {
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
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 10 + rand.Intn(30)
		availableCount := totalCount - rand.Intn(5)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("HZ-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           roundPrice(price),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 空调, 电视, 24小时热水, 吹风机, 电热水壶",
			ImageURL:        getRoomImageURL("huazhu", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

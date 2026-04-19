package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type JinJiangAdapter struct{}

func NewJinJiangAdapter() *JinJiangAdapter {
	return &JinJiangAdapter{}
}

func (j *JinJiangAdapter) GetCode() string {
	return "jinjiang"
}

func (j *JinJiangAdapter) GetName() string {
	return "锦江国际酒店集团"
}

func (j *JinJiangAdapter) GetDescription() string {
	return "锦江国际集团是中国规模最大的综合性旅游企业集团之一，旗下拥有锦江饭店、昆仑饭店、锦江之星等知名酒店品牌。"
}

func (j *JinJiangAdapter) GetAPIURL() string {
	return "/api/mock/jinjiang"
}

func (j *JinJiangAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return generateJinJiangMockData(), nil
}

func (j *JinJiangAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := generateJinJiangMockData()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func generateJinJiangMockData() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + 2000)
	
	cities := []string{"上海", "北京", "成都", "重庆", "武汉", "长沙", "郑州", "济南", "沈阳", "哈尔滨"}
	hotelTypes := []string{"锦江酒店", "锦江之星", "白玉兰酒店", "丽枫酒店", "喆啡酒店", "潮漫酒店", "7天连锁", "IU酒店"}
	addressSuffixes := []string{"人民广场店", "外滩中心店", "中关村店", "春熙路店", "解放碑店", "光谷店", "五一广场店", "二七广场店", "泉城广场店", "中央大街店"}
	
	hotels := make([]SupplierHotelData, 11)
	
	for i := 0; i < 11; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 180 + rand.Intn(250)
		maxPrice := minPrice + 180 + rand.Intn(250)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("JJ-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s大道%d号", city, getRandomDistrict(city), getRandomRoad(), 80+rand.Intn(600)),
			Description:     fmt.Sprintf("%s是锦江国际集团旗下品牌，秉承自然、健康、人性化的生活方式，为您提供物超所值的住宿体验。酒店设计时尚现代，设施齐全，是商务出行和休闲旅游的理想选择。", hotelType),
			Rating:          4.1 + float64(rand.Intn(9))/10,
			ImageURL:        getHotelImageURL("jinjiang", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           generateJinJiangRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateJinJiangRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name        string
		description string
		priceMulti  float64
		capacity    int
		area        int
		bedType     string
	}{
		{"经济房", "简洁舒适的经济客房，适合预算敏感的旅客。房间干净整洁，配备基本生活设施。", 0.85, 2, 20, "大床"},
		{"标准双床房", "经济实惠的标准客房，两张单人床，适合商务出行。配备免费WiFi、空调、电视和24小时热水。", 1.0, 2, 24, "双床"},
		{"商务大床房", "宽敞舒适的商务客房，配备优质床品和现代化设施。适合商务旅客和休闲游客。", 1.4, 2, 35, "大床"},
		{"豪华套房", "高端豪华套房，独立客厅和卧室，配备高品质家具和设备。享受尊贵的住宿体验。", 2.2, 2, 55, "大床"},
		{"家庭亲子房", "温馨舒适的家庭房，配备大床和儿童床，适合家庭出行。空间宽敞，设施齐全。", 1.6, 3, 40, "大床+单人床"},
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		totalCount := 12 + rand.Intn(28)
		availableCount := totalCount - rand.Intn(6)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("JJ-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:            rt.name,
			Description:     rt.description,
			Price:           float64(int(price/10) * 10),
			Capacity:        rt.capacity,
			Area:            rt.area,
			BedType:         rt.bedType,
			Amenities:       "免费WiFi, 空调, 电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水",
			ImageURL:        getRoomImageURL("jinjiang", hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

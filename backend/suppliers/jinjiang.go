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

func (j *JinJiangAdapter) GetPriority() int {
	return 2
}

func (j *JinJiangAdapter) GetColor() string {
	return "#1E88E5"
}

func (j *JinJiangAdapter) GetIcon() string {
	return "🏩"
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
	rand.Seed(time.Now().UnixNano() + 20000)
	
	cities := []string{"上海", "北京", "成都", "重庆", "武汉", "长沙", "郑州", "济南", "沈阳", "哈尔滨"}
	hotelTypes := []string{"锦江酒店", "锦江之星", "白玉兰酒店", "丽枫酒店", "喆啡酒店", "潮漫酒店", "7天连锁", "IU酒店"}
	addressSuffixes := []string{"人民广场店", "外滩中心店", "中关村店", "春熙路店", "解放碑店", "光谷店", "五一广场店", "二七广场店", "泉城广场店", "中央大街店"}
	
	brandMap := map[string]string{
		"锦江酒店":  "锦江",
		"锦江之星":  "锦江之星",
		"白玉兰酒店": "白玉兰",
		"丽枫酒店":  "丽枫",
		"喆啡酒店":  "喆啡",
		"潮漫酒店":  "潮漫",
		"7天连锁":   "7天",
		"IU酒店":   "IU",
	}
	
	hotels := make([]SupplierHotelData, 11)
	
	for i := 0; i < 11; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := 180 + rand.Intn(250)
		maxPrice := minPrice + 180 + rand.Intn(250)
		brand := brandMap[hotelType]
		if brand == "" {
			brand = "锦江"
		}
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("JJ-HOTEL-%04d", i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s大道%d号", city, getRandomDistrict(city), getRandomRoad(), 80+rand.Intn(600)),
			Description:     fmt.Sprintf("%s是锦江国际集团旗下品牌，秉承自然、健康、人性化的生活方式，为您提供物超所值的住宿体验。酒店设计时尚现代，设施齐全，是商务出行和休闲旅游的理想选择。", hotelType),
			Rating:          roundToOneDecimal(4.1 + float64(rand.Intn(9))/10),
			ImageURL:        getHotelImageURL("jinjiang", i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			MinPrice:        float64(minPrice),
			MaxPrice:        float64(maxPrice),
			Brand:           brand,
			Rooms:           generateJinJiangRooms(i, minPrice),
		}
	}
	
	return hotels
}

func generateJinJiangRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := []struct {
		name             string
		description      string
		priceMulti       float64
		capacity         int
		area             int
		bedType          string
		standardRoomType string
	}{
		{"经济房", "简洁舒适的经济客房，适合预算敏感的旅客。房间干净整洁，配备基本生活设施。", 0.85, 2, 20, "大床", "标准大床房"},
		{"标准双床房", "经济实惠的标准客房，两张单人床，适合商务出行。配备免费WiFi、空调、电视和24小时热水。", 1.0, 2, 24, "双床", "标准双床房"},
		{"商务大床房", "宽敞舒适的商务客房，配备优质床品和现代化设施。适合商务旅客和休闲游客。", 1.4, 2, 35, "大床", "豪华大床房"},
		{"豪华套房", "高端豪华套房，独立客厅和卧室，配备高品质家具和设备。享受尊贵的住宿体验。", 2.2, 2, 55, "大床", "行政套房"},
		{"家庭亲子房", "温馨舒适的家庭房，配备大床和儿童床，适合家庭出行。空间宽敞，设施齐全。", 1.6, 3, 40, "大床+单人床", "家庭房"},
	}
	
	promotionTags := []string{"会员专享价", "新客特惠", "限时闪购", "连住2晚9折", "", "", ""}
	paymentTypes := []string{"现付", "预付", "信用住"}
	cancelPolicies := []string{"免费取消(入住前1天)", "部分取消", "不可取消"}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.priceMulti
		originalPrice := price * (1.08 + float64(rand.Intn(25))/100)
		totalCount := 12 + rand.Intn(28)
		availableCount := totalCount - rand.Intn(6)
		
		isPriceControlled := false
		priceControlReason := ""
		if i == 1 && rand.Intn(2) == 0 {
			isPriceControlled = true
			priceControlReason = "集团控价"
		}
		
		promotionTag := ""
		if rand.Intn(3) == 0 {
			promotionTag = promotionTags[rand.Intn(len(promotionTags))]
		}
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID:     fmt.Sprintf("JJ-ROOM-%04d-%02d", hotelIndex+1, i+1),
			Name:               rt.name,
			Description:        rt.description,
			Price:              roundPrice(price),
			OriginalPrice:      roundPrice(originalPrice),
			Capacity:           rt.capacity,
			Area:               rt.area,
			BedType:            rt.bedType,
			StandardRoomType:   rt.standardRoomType,
			Amenities:          "免费WiFi, 空调, 电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水",
			ImageURL:           getRoomImageURL("jinjiang", hotelIndex, i),
			TotalCount:         totalCount,
			AvailableCount:     availableCount,
			IsPriceControlled:  isPriceControlled,
			PriceControlReason: priceControlReason,
			PromotionTag:       promotionTag,
			PaymentType:        paymentTypes[rand.Intn(len(paymentTypes))],
			CancelPolicy:       cancelPolicies[rand.Intn(len(cancelPolicies))],
		}
	}
	
	return rooms
}

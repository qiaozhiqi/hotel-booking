package suppliers

import (
	"fmt"
	"math/rand"
	"time"
)

type ShijiSubSupplierConfig struct {
	Code           string
	Name           string
	Description    string
	APIURL         string
	HotelTypes     []string
	Cities         []string
	AddressSuffixes []string
	MinPriceBase   int
	MaxPriceBase   int
	RatingBase     float64
	RatingRange    float64
	HotelIDPrefix  string
	RoomIDPrefix   string
	ImageThemes    []string
	RoomTypes      []ShijiRoomTypeConfig
	HotelCount     int
}

type ShijiRoomTypeConfig struct {
	Name        string
	Description string
	PriceMulti  float64
	Capacity    int
	Area        int
	BedType     string
	Amenities   string
}

type ShijiGenericAdapter struct {
	config ShijiSubSupplierConfig
}

func NewShijiGenericAdapter(config ShijiSubSupplierConfig) *ShijiGenericAdapter {
	return &ShijiGenericAdapter{config: config}
}

func (s *ShijiGenericAdapter) GetCode() string {
	return s.config.Code
}

func (s *ShijiGenericAdapter) GetName() string {
	return s.config.Name
}

func (s *ShijiGenericAdapter) GetDescription() string {
	return s.config.Description
}

func (s *ShijiGenericAdapter) GetAPIURL() string {
	return s.config.APIURL
}

func (s *ShijiGenericAdapter) FetchHotels() ([]SupplierHotelData, error) {
	return s.generateMockHotels(), nil
}

func (s *ShijiGenericAdapter) FetchHotelDetail(hotelID string) (*SupplierHotelData, error) {
	hotels := s.generateMockHotels()
	for _, hotel := range hotels {
		if hotel.SupplierHotelID == hotelID {
			return &hotel, nil
		}
	}
	return nil, fmt.Errorf("hotel not found: %s", hotelID)
}

func (s *ShijiGenericAdapter) generateMockHotels() []SupplierHotelData {
	rand.Seed(time.Now().UnixNano() + int64(s.hashCode(s.config.Code)))
	
	cities := s.config.Cities
	if len(cities) == 0 {
		cities = []string{"上海", "北京", "广州", "深圳", "杭州", "南京", "成都", "武汉", "西安", "重庆"}
	}
	
	hotelTypes := s.config.HotelTypes
	if len(hotelTypes) == 0 {
		hotelTypes = []string{"精选酒店", "商务酒店", "度假酒店"}
	}
	
	addressSuffixes := s.config.AddressSuffixes
	if len(addressSuffixes) == 0 {
		addressSuffixes = []string{"中心店", "广场店", "路店"}
	}
	
	hotelCount := s.config.HotelCount
	if hotelCount <= 0 {
		hotelCount = 10
	}
	
	hotels := make([]SupplierHotelData, hotelCount)
	
	for i := 0; i < hotelCount; i++ {
		city := cities[i%len(cities)]
		hotelType := hotelTypes[i%len(hotelTypes)]
		addressSuffix := addressSuffixes[i%len(addressSuffixes)]
		
		minPrice := s.config.MinPriceBase + rand.Intn(s.config.MaxPriceBase)
		maxPrice := minPrice + s.config.MinPriceBase/2 + rand.Intn(s.config.MaxPriceBase/2)
		
		hotels[i] = SupplierHotelData{
			SupplierHotelID: fmt.Sprintf("%s-HOTEL-%04d", s.config.HotelIDPrefix, i+1),
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, getRandomDistrict(city), getRandomRoad(), 50+rand.Intn(500)),
			Description:     fmt.Sprintf("%s是%s旗下的酒店品牌，为宾客提供优质的住宿体验。酒店地理位置优越，设施完善，服务周到，是商务出行和休闲旅游的理想选择。", hotelType, s.config.Name),
			Rating:          roundToOneDecimal(s.config.RatingBase + float64(rand.Intn(int(s.config.RatingRange*10)))/10),
			ImageURL:        s.getHotelImageURL(i),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           s.generateRooms(i, minPrice),
		}
	}
	
	return hotels
}

func (s *ShijiGenericAdapter) generateRooms(hotelIndex int, basePrice int) []SupplierRoomData {
	roomTypes := s.config.RoomTypes
	if len(roomTypes) == 0 {
		roomTypes = []ShijiRoomTypeConfig{
			{"标准双床房", "经济实惠的标准客房，适合商务出行。", 1.0, 2, 25, "双床", "免费WiFi, 空调, 电视, 24小时热水"},
			{"豪华大床房", "宽敞舒适的豪华客房，配备高品质设施。", 1.3, 2, 32, "大床", "免费WiFi, 空调, 电视, 24小时热水, 迷你吧"},
			{"行政套房", "高端行政套房，独立客厅和卧室。", 2.0, 2, 50, "大床", "免费WiFi, 空调, 电视, 24小时热水, 迷你吧, 保险箱"},
		}
	}
	
	rooms := make([]SupplierRoomData, len(roomTypes))
	for i, rt := range roomTypes {
		price := float64(basePrice) * rt.PriceMulti
		totalCount := 10 + rand.Intn(25)
		availableCount := totalCount - rand.Intn(5)
		
		rooms[i] = SupplierRoomData{
			SupplierRoomID: fmt.Sprintf("%s-ROOM-%04d-%02d", s.config.RoomIDPrefix, hotelIndex+1, i+1),
			Name:            rt.Name,
			Description:     rt.Description,
			Price:           roundPrice(price),
			Capacity:        rt.Capacity,
			Area:            rt.Area,
			BedType:         rt.BedType,
			Amenities:       rt.Amenities,
			ImageURL:        getRoomImageURL(s.config.Code, hotelIndex, i),
			TotalCount:      totalCount,
			AvailableCount:  availableCount,
		}
	}
	
	return rooms
}

func (s *ShijiGenericAdapter) getHotelImageURL(index int) string {
	themes := s.config.ImageThemes
	if len(themes) == 0 {
		themes = []string{
			"modern luxury hotel lobby interior design",
			"elegant hotel entrance with glass facade",
			"contemporary hotel building night view",
		}
	}
	theme := themes[index%len(themes)]
	return fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=%s&image_size=landscape_4_3", theme)
}

func (s *ShijiGenericAdapter) hashCode(str string) int {
	h := 0
	for i := 0; i < len(str); i++ {
		h = 31*h + int(str[i])
	}
	return h
}

var ShijiSubSuppliers = []ShijiSubSupplierConfig{
	{
		Code:         "shiji_marriott",
		Name:         "万豪国际酒店集团（石基畅联）",
		Description:  "万豪国际集团是全球首屈一指的酒店管理公司，旗下拥有万豪、喜来登、丽思卡尔顿等多个知名酒店品牌。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_marriott",
		HotelTypes:   []string{"万豪酒店", "喜来登酒店", "威斯汀酒店", "丽思卡尔顿酒店", "JW万豪酒店", "万怡酒店", "万丽酒店", "福朋喜来登"},
		Cities:       []string{"上海", "北京", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "苏州", "三亚"},
		AddressSuffixes: []string{"陆家嘴店", "国贸店", "天河店", "福田店", "西湖店", "春熙路店", "解放碑店", "光谷店", "大雁塔店", "玄武湖店", "金鸡湖店", "亚龙湾店"},
		MinPriceBase: 500,
		MaxPriceBase: 500,
		RatingBase:   4.5,
		RatingRange:  0.5,
		HotelIDPrefix: "SJ-MA",
		RoomIDPrefix:  "SJ-MA",
		HotelCount:   12,
		ImageThemes: []string{
			"marriott luxury hotel grand lobby chandelier",
			"elegant hotel entrance with valet parking",
			"modern skyscraper hotel night view city lights",
			"luxury hotel suite marble bathroom",
			"ritz carlton style lobby with classical decor",
			"sheraton hotel exterior glass facade",
			"westin hotel lobby with modern art",
			"jw marriott executive lounge view",
			"courtyard by marriott exterior clean design",
			"renaissance hotel boutique style interior",
			"marriott resort infinity pool ocean view",
			"marriott bonvoy elite lounge",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"豪华客房", "宽敞明亮的豪华客房，配备高品质床品和现代化设施。享受城市美景，是商务和休闲的理想选择。", 1.0, 2, 35, "大床/双床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费矿泉水"},
			{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊礼遇。配备宽大办公桌和高速网络，专为商务精英打造。", 1.4, 2, 42, "大床/双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"小型套房", "独立客厅与卧室的小型套房，空间宽敞舒适。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.0, 2, 55, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费茶包"},
			{"行政套房", "高端行政套房，独立客厅和卧室，配备宽大办公桌和高速网络。享受行政楼层礼遇，适合高端商务人士。", 2.8, 2, 75, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务"},
			{"豪华套房", "顶级豪华套房，宽敞的独立客厅和卧室，配备奢华家具和设施。享受私人管家服务，是尊贵宾客的首选。", 4.0, 3, 120, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务, 专属餐厅"},
		},
	},
	{
		Code:         "shiji_hilton",
		Name:         "希尔顿酒店集团（石基畅联）",
		Description:  "希尔顿全球控股有限公司是全球最知名的酒店管理公司之一，旗下拥有希尔顿、华尔道夫、康莱德等多个奢华及高端酒店品牌。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_hilton",
		HotelTypes:   []string{"希尔顿酒店", "华尔道夫酒店", "康莱德酒店", "希尔顿逸林酒店", "希尔顿花园酒店", "欢朋酒店", "希尔顿嘉悦里"},
		Cities:       []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "厦门", "青岛"},
		AddressSuffixes: []string{"王府井店", "外滩店", "珠江新城店", "华侨城店", "钱江新城店", "高新店", "江北嘴店", "江汉路店", "曲江店", "河西店", "环岛路店", "五四广场店"},
		MinPriceBase: 480,
		MaxPriceBase: 450,
		RatingBase:   4.4,
		RatingRange:  0.6,
		HotelIDPrefix: "SJ-HI",
		RoomIDPrefix:  "SJ-HI",
		HotelCount:   12,
		ImageThemes: []string{
			"hilton hotel grand lobby marble columns",
			"waldorf astoria luxury entrance",
			"conrad hotel modern minimalist design",
			"hilton suite with panoramic city view",
			"doubletree by hilton warm welcoming lobby",
			"hilton garden inn modern exterior",
			"hampton by hilton bright clean interior",
			"canopy by hilton boutique style lobby",
			"hilton executive club lounge",
			"waldorf astoria presidential suite",
			"conrad hotel art collection display",
			"hilton resort beachfront view",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"客房", "舒适温馨的标准客房，配备希尔顿特色睡床和现代化设施。享受优质睡眠体验，是商务出行的理想选择。", 1.0, 2, 32, "大床/双床", "免费WiFi, 中央空调, 希尔顿睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费咖啡茶包"},
			{"豪华客房", "宽敞的豪华客房，配备高品质家具和现代化设施。享有城市美景，提供更舒适的住宿体验。", 1.3, 2, 40, "大床/双床", "免费WiFi, 中央空调, 希尔顿睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊服务。配备办公区域和高速网络，专为商务旅客设计。", 1.6, 2, 45, "大床/双床", "免费WiFi, 中央空调, 希尔顿睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 行政酒廊礼遇"},
			{"套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.2, 2, 60, "大床", "免费WiFi, 中央空调, 希尔顿睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 独立客厅"},
			{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和私人管家服务，是尊贵宾客的首选。", 5.0, 4, 150, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务, 专属餐厅"},
		},
	},
	{
		Code:         "shiji_ihg",
		Name:         "洲际酒店集团（石基畅联）",
		Description:  "洲际酒店集团（IHG）是全球最大的酒店管理公司之一，旗下拥有洲际、皇冠假日、假日酒店等多个知名品牌。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_ihg",
		HotelTypes:   []string{"洲际酒店", "皇冠假日酒店", "假日酒店", "智选假日酒店", "英迪格酒店", "VOCO酒店", "逸衡酒店"},
		Cities:       []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "重庆", "武汉", "西安", "南京", "天津", "大连"},
		AddressSuffixes: []string{"三里屯店", "静安店", "白云店", "南山店", "滨江店", "锦江区店", "渝中区店", "武昌店", "高新区店", "鼓楼店", "和平店", "中山区店"},
		MinPriceBase: 450,
		MaxPriceBase: 400,
		RatingBase:   4.3,
		RatingRange:  0.7,
		HotelIDPrefix: "SJ-IHG",
		RoomIDPrefix:  "SJ-IHG",
		HotelCount:   12,
		ImageThemes: []string{
			"intercontinental hotel grand lobby",
			"crowne plaza business hotel exterior",
			"holiday inn modern family friendly lobby",
			"holiday inn express clean simple design",
			"hotel indigo boutique style interior",
			"voco hotel colorful vibrant lobby",
			"even wellness hotel fitness area",
			"intercontinental suite luxury decor",
			"crowne plaza club lounge",
			"holiday inn family suite",
			"hotel indigo local art display",
			"ihg rewards club lounge",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"标准房", "舒适的标准客房，配备优质床品和基本设施。干净整洁，是商务出行的经济实惠选择。", 1.0, 2, 28, "大床/双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"豪华房", "宽敞的豪华客房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.25, 2, 36, "大床/双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"行政房", "位于行政楼层的高级客房，可享受行政俱乐部服务。配备办公区域，专为商务旅客设计。", 1.5, 2, 42, "大床/双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 行政俱乐部礼遇"},
			{"套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具，享受尊贵住宿体验。", 2.0, 2, 58, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 独立客厅"},
			{"洲际套房", "高端洲际套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 3.5, 3, 100, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务"},
		},
	},
	{
		Code:         "shiji_kaiyuan",
		Name:         "开元酒店集团（石基畅联）",
		Description:  "开元酒店集团是中国最大的民营高星级连锁酒店集团，旗下拥有开元名都、开元大酒店、开元度假村等多个知名品牌。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_kaiyuan",
		HotelTypes:   []string{"开元名都大酒店", "开元大酒店", "开元度假村", "开元名庭酒店", "开元曼居酒店", "开元颐居酒店", "开元森泊度假酒店"},
		Cities:       []string{"杭州", "上海", "北京", "宁波", "温州", "绍兴", "嘉兴", "湖州", "金华", "台州", "南京", "苏州"},
		AddressSuffixes: []string{"西湖店", "松江店", "朝阳店", "鄞州店", "鹿城店", "越城店", "南湖店", "吴兴店", "婺城店", "椒江店", "江宁店", "姑苏店"},
		MinPriceBase: 380,
		MaxPriceBase: 350,
		RatingBase:   4.2,
		RatingRange:  0.8,
		HotelIDPrefix: "SJ-KY",
		RoomIDPrefix:  "SJ-KY",
		HotelCount:   12,
		ImageThemes: []string{
			"kaiyuan chinese luxury hotel lobby",
			"new century grand hotel entrance",
			"kaiyuan resort with chinese garden",
			"kaiyuan mingting business hotel",
			"kaiyuan manju boutique hotel",
			"kaiyuan yiju eco hotel",
			"kaiyuan senbo resort water park",
			"kaiyuan luxury suite oriental style",
			"kaiyuan executive lounge",
			"kaiyuan resort mountain view",
			"kaiyuan hotel chinese restaurant",
			"kaiyuan wellness spa area",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"标准双床房", "舒适的标准双床房，配备开元特色睡床和现代化设施。房间干净整洁，是商务出行的理想选择。", 1.0, 2, 26, "双床", "免费WiFi, 中央空调, 开元睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费茶包"},
			{"豪华大床房", "宽敞的豪华大床房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.2, 2, 32, "大床", "免费WiFi, 中央空调, 开元睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"行政客房", "位于行政楼层的高级客房，可享受行政酒廊服务。配备办公区域，专为商务旅客设计。", 1.5, 2, 38, "大床/双床", "免费WiFi, 中央空调, 开元睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 行政酒廊礼遇"},
			{"豪华套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具和现代化设施，享受尊贵住宿体验。", 2.0, 2, 52, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 独立客厅"},
			{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 4.0, 4, 120, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务"},
		},
	},
	{
		Code:         "shiji_wanda",
		Name:         "万达酒店集团（石基畅联）",
		Description:  "万达酒店集团是中国知名的高端酒店管理集团，旗下拥有万达文华、万达嘉华、万达瑞华等多个奢华及高端酒店品牌。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_wanda",
		HotelTypes:   []string{"万达文华酒店", "万达嘉华酒店", "万达瑞华酒店", "万达锦华酒店", "万达美华酒店", "万达颐华酒店", "万达悦华酒店"},
		Cities:       []string{"北京", "上海", "广州", "深圳", "成都", "重庆", "武汉", "西安", "南京", "杭州", "青岛", "大连"},
		AddressSuffixes: []string{"CBD店", "五角场店", "海珠店", "宝安店", "金牛店", "南岸店", "洪山店", "未央店", "建邺店", "拱墅店", "市南店", "西岗店"},
		MinPriceBase: 420,
		MaxPriceBase: 380,
		RatingBase:   4.3,
		RatingRange:  0.7,
		HotelIDPrefix: "SJ-WD",
		RoomIDPrefix:  "SJ-WD",
		HotelCount:   12,
		ImageThemes: []string{
			"wanda vista luxury hotel grand lobby",
			"wanda realm hotel elegant entrance",
			"wanda jinjiang business hotel",
			"wanda meihua modern design hotel",
			"wanda yihua boutique hotel",
			"wanda yuehua budget hotel",
			"wanda suite oriental luxury style",
			"wanda executive club lounge",
			"wanda hotel ballroom entrance",
			"wanda resort infinity pool",
			"wanda presidential suite",
			"wanda hotel spa wellness",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"豪华大床房", "宽敞舒适的豪华大床房，配备万达特色睡床和现代化设施。房间设计典雅，是商务出行的理想选择。", 1.0, 2, 30, "大床", "免费WiFi, 中央空调, 万达睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 免费矿泉水"},
			{"行政套房", "位于行政楼层的高级套房，独立客厅与卧室。可享受行政酒廊服务，专为商务精英打造。", 1.8, 2, 48, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 行政酒廊礼遇"},
			{"豪华套房", "宽敞的豪华套房，独立客厅与卧室，配备高品质家具。享受城市美景，提供尊贵的住宿体验。", 2.5, 2, 65, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 独立客厅"},
			{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室、餐厅和书房。配备奢华家具和私人管家服务，是尊贵宾客的首选。", 5.0, 4, 180, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务, 专属餐厅"},
		},
	},
	{
		Code:         "shiji_lvdi",
		Name:         "绿地酒店集团（石基畅联）",
		Description:  "绿地酒店集团是中国知名的酒店管理集团，旗下拥有绿地铂瑞、绿地铂骊、绿地假日等多个酒店品牌，致力于为宾客提供高品质的住宿体验。通过石基畅联渠道接入。",
		APIURL:       "/api/mock/shiji_lvdi",
		HotelTypes:   []string{"绿地铂瑞酒店", "绿地铂骊酒店", "绿地假日酒店", "绿地康养居酒店", "绿地魔奇酒店"},
		Cities:       []string{"上海", "北京", "广州", "深圳", "南京", "杭州", "成都", "武汉", "西安", "重庆", "苏州", "无锡"},
		AddressSuffixes: []string{"黄浦店", "朝阳店", "天河店", "龙华店", "鼓楼店", "西湖店", "武侯店", "武昌店", "雁塔店", "江北店", "工业园店", "滨湖店"},
		MinPriceBase: 350,
		MaxPriceBase: 300,
		RatingBase:   4.1,
		RatingRange:  0.9,
		HotelIDPrefix: "SJ-LD",
		RoomIDPrefix:  "SJ-LD",
		HotelCount:   12,
		ImageThemes: []string{
			"greenland primus luxury hotel lobby",
			"greenland qhotel modern entrance",
			"greenland holiday hotel exterior",
			"greenland kangyangju wellness hotel",
			"greenland moqi boutique hotel",
			"greenland suite contemporary design",
			"greenland executive lounge",
			"greenland business center",
			"greenland resort garden view",
			"greenland hotel restaurant interior",
			"greenland family friendly suite",
			"greenland eco hotel design",
		},
		RoomTypes: []ShijiRoomTypeConfig{
			{"标准双床房", "舒适的标准双床房，配备优质床品和现代化设施。房间干净整洁，是商务出行的经济实惠选择。", 1.0, 2, 24, "双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"豪华大床房", "宽敞的豪华大床房，配备高品质家具和现代化设施。享受城市美景，提供更舒适的住宿体验。", 1.15, 2, 30, "大床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋"},
			{"行政客房", "位于行政楼层的高级客房，可享受行政服务。配备办公区域，专为商务旅客设计。", 1.4, 2, 36, "大床/双床", "免费WiFi, 中央空调, 优质睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 行政服务"},
			{"豪华套房", "独立客厅与卧室的豪华套房，空间宽敞明亮。配备高品质家具，享受尊贵住宿体验。", 1.8, 2, 48, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 独立客厅"},
			{"总统套房", "顶级总统套房，拥有宽敞的客厅、卧室和餐厅。配备奢华家具和专属服务，是尊贵宾客的首选。", 3.5, 4, 100, "大床", "免费WiFi, 中央空调, 奢华睡床, 智能电视, 24小时热水, 迷你吧, 保险箱, 熨斗, 浴袍拖鞋, 私人管家服务"},
		},
	},
}

func InitShijiSubSuppliers() {
	for _, config := range ShijiSubSuppliers {
		RegisterAdapter(NewShijiGenericAdapter(config))
	}
}

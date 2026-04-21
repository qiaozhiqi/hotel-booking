package suppliers

import (
	"fmt"
	"math"
)

func roundToOneDecimal(val float64) float64 {
	return math.Round(val*10) / 10
}

func roundPrice(price float64) float64 {
	return math.Round(price/10) * 10
}

func getRandomDistrict(city string) string {
	districts := map[string][]string{
		"上海": {"黄浦", "徐汇", "长宁", "静安", "普陀", "虹口", "杨浦", "浦东"},
		"北京": {"东城", "西城", "朝阳", "海淀", "丰台", "石景山"},
		"广州": {"天河", "越秀", "海珠", "荔湾", "白云", "番禺"},
		"深圳": {"南山", "福田", "罗湖", "宝安", "龙岗", "龙华"},
		"杭州": {"西湖", "上城", "下城", "江干", "拱墅", "滨江"},
		"南京": {"玄武", "秦淮", "建邺", "鼓楼", "浦口", "栖霞"},
		"成都": {"锦江", "青羊", "金牛", "武侯", "成华", "高新"},
		"武汉": {"江岸", "江汉", "硚口", "汉阳", "武昌", "青山"},
		"西安": {"新城", "碑林", "莲湖", "雁塔", "未央", "灞桥"},
		"重庆": {"渝中", "大渡口", "江北", "沙坪坝", "九龙坡", "南岸"},
		"苏州": {"姑苏", "虎丘", "吴中", "相城", "苏州工业园", "吴江"},
		"无锡": {"锡山", "惠山", "滨湖", "梁溪", "新吴", "江阴"},
		"宁波": {"海曙", "江北", "北仑", "镇海", "鄞州", "奉化"},
		"福州": {"鼓楼", "台江", "仓山", "马尾", "晋安", "长乐"},
		"厦门": {"思明", "海沧", "湖里", "集美", "同安", "翔安"},
		"青岛": {"市南", "市北", "黄岛", "崂山", "李沧", "城阳"},
		"长沙": {"芙蓉", "天心", "岳麓", "开福", "雨花", "望城"},
		"郑州": {"中原", "二七", "管城", "金水", "上街", "惠济"},
		"济南": {"历下", "市中", "槐荫", "天桥", "历城", "长清"},
		"沈阳": {"和平", "沈河", "大东", "皇姑", "铁西", "苏家屯"},
		"哈尔滨": {"道里", "南岗", "道外", "平房", "松北", "香坊"},
	}
	if d, ok := districts[city]; ok {
		return d[0]
	}
	return "中心"
}

func getRandomRoad() string {
	roads := []string{"中山", "人民", "解放", "建国", "和平", "新华", "建设", "发展", "科技", "商务", "世纪", "广场", "友谊", "工农", "跃进"}
	return roads[0]
}

func getHotelImageURL(supplierCode string, index int) string {
	themes := map[string][]string{
		"huazhu": {
			"modern luxury hotel lobby interior design",
			"elegant hotel entrance with glass facade",
			"contemporary hotel building night view",
			"luxury hotel suite with city view",
			"boutique hotel lobby with art decor",
			"business hotel exterior modern architecture",
			"resort hotel with swimming pool",
			"city hotel with skyline view",
			"designer hotel interior minimalist style",
			"grand hotel entrance with porte cochere",
			"urban hotel with green rooftop",
			"classic hotel with modern renovation",
		},
		"rujia": {
			"budget hotel exterior clean simple design",
			"economy hotel reception friendly staff",
			"simple hotel room tidy organized",
			"chain hotel building modern exterior",
			"affordable hotel lobby welcoming",
			"standard hotel room basic amenities",
			"city budget hotel convenient location",
			"comfortable hotel room clean bedding",
			"value hotel exterior professional",
			"practical hotel room functional design",
		},
		"jinjiang": {
			"international hotel grand entrance",
			"luxury hotel lobby marble floor",
			"premium hotel suite elegant design",
			"five star hotel exterior impressive",
			"high end hotel room luxurious furnishing",
			"hotel chain building modern skyline",
			"upscale hotel reception grand",
			"deluxe hotel room city view",
			"premium hotel building contemporary",
			"luxury hotel amenities premium",
			"world class hotel exterior iconic",
		},
		"shiji_marriott": {
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
		"shiji_hilton": {
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
		"shiji_ihg": {
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
		"shiji_kaiyuan": {
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
		"shiji_wanda": {
			"wanda vista luxury hotel grand lobby",
			"wanda realm hotel elegant entrance",
			"wanda jinjiang business hotel",
			"wanda mehua modern design hotel",
			"wanda yihua boutique hotel",
			"wanda yuehua budget hotel",
			"wanda suite oriental luxury style",
			"wanda executive club lounge",
			"wanda hotel ballroom entrance",
			"wanda resort infinity pool",
			"wanda presidential suite",
			"wanda hotel spa wellness",
		},
		"shiji_lvdi": {
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
	}
	
	supplierThemes, ok := themes[supplierCode]
	if !ok {
		supplierThemes = themes["huazhu"]
	}
	
	theme := supplierThemes[index%len(supplierThemes)]
	return fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=%s&image_size=landscape_4_3", theme)
}

func getRoomImageURL(supplierCode string, hotelIndex, roomIndex int) string {
	themes := []string{
		"luxury hotel bedroom king size bed",
		"modern hotel room twin beds",
		"hotel suite with living area",
		"deluxe hotel room with city view",
		"executive hotel room work desk",
		"comfortable hotel guest room",
		"premium hotel bedroom design",
		"standard hotel room clean simple",
	}
	
	theme := themes[roomIndex%len(themes)]
	return fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=%s&image_size=landscape_4_3", theme)
}

func InitSuppliers() {
	RegisterAdapter(NewHuazhuAdapter())
	RegisterAdapter(NewRuJiaAdapter())
	RegisterAdapter(NewJinJiangAdapter())
	
	RegisterAdapter(NewShijiMarriottAdapter())
	RegisterAdapter(NewShijiHiltonAdapter())
	RegisterAdapter(NewShijiIHGAdapter())
	RegisterAdapter(NewShijiKaiyuanAdapter())
	RegisterAdapter(NewShijiWandaAdapter())
	RegisterAdapter(NewShijiLvdiAdapter())
}

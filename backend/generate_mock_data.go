//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"hotel-booking/models"
	"math/rand"
	"os"
	"time"
)

const (
	NumHotels    = 3
	NumRooms     = 5
	NumDays      = 90
	BasePrice    = 280.0
	MaxPriceVar  = 100.0
	BaseTotal    = 15
)

func generateMockData() models.QiuguoPushRequest {
	rand.Seed(time.Now().UnixNano())

	hotels := make([]models.QiuguoPushHotelData, NumHotels)
	priceInventories := make([]models.QiuguoPushPriceInventoryData, 0, NumHotels*NumRooms*NumDays)

	cities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "武汉", "西安", "南京", "苏州"}
	districts := map[string][]string{
		"北京": {"海淀", "朝阳", "西城", "东城", "丰台", "石景山"},
		"上海": {"黄浦", "徐汇", "长宁", "静安", "普陀", "浦东"},
		"广州": {"天河", "越秀", "海珠", "荔湾", "白云", "番禺"},
		"深圳": {"南山", "福田", "罗湖", "宝安", "龙岗", "龙华"},
		"杭州": {"西湖", "上城", "下城", "江干", "拱墅", "滨江"},
		"成都": {"锦江", "青羊", "金牛", "武侯", "成华", "高新"},
		"武汉": {"江岸", "江汉", "硚口", "汉阳", "武昌", "青山"},
		"西安": {"新城", "碑林", "莲湖", "雁塔", "未央", "灞桥"},
		"南京": {"玄武", "秦淮", "建邺", "鼓楼", "浦口", "栖霞"},
		"苏州": {"姑苏", "虎丘", "吴中", "相城", "吴江", "工业园区"},
	}

	hotelTypes := []string{"秋果酒店", "秋果S酒店", "秋果Q酒店"}
	addressSuffixes := []string{"中关村店", "国贸店", "天河店", "南山店", "西湖店", "春熙路店", "光谷店", "大雁塔店", "新街口店", "观前街店"}

	roomTypes := []struct {
		name        string
		description string
		capacity    int
		area        int
		bedType     string
		amenities   string
		priceMulti  float64
	}{
		{"精选大床房", "舒适温馨的精选大床房，配备高品质床品和现代化设施。简约设计风格，为您提供优质的睡眠体验。", 2, 26, "大床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水", 1.0},
		{"精选双床房", "干净整洁的精选双床房，两张舒适单人床，适合双人入住或商务出行。", 2, 28, "双床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水", 1.05},
		{"高级大床房", "空间更宽敞的高级大床房，配备高品质家具和床品。享有城市美景，提供更舒适的住宿体验。", 2, 32, "大床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水, 迷你吧", 1.3},
		{"商务套房", "高端商务套房，独立客厅与卧室，配备宽大办公桌和高速网络。专为商务精英打造，享受尊贵住宿体验。", 2, 48, "大床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水, 迷你吧, 保险箱", 2.0},
		{"家庭亲子房", "温馨舒适的家庭亲子房，配备大床和儿童床，适合家庭出行。空间宽敞，设施齐全，提供温馨的家庭住宿体验。", 4, 45, "大床+儿童床", "免费WiFi, 中央空调, 智能电视, 24小时热水, 吹风机, 电热水壶, 免费矿泉水, 儿童用品", 1.8},
	}

	// 生成酒店数据
	for h := 0; h < NumHotels; h++ {
		city := cities[rand.Intn(len(cities))]
		cityDistricts := districts[city]
		district := cityDistricts[rand.Intn(len(cityDistricts))]
		hotelType := hotelTypes[h%len(hotelTypes)]
		addressSuffix := addressSuffixes[h%len(addressSuffixes)]

		minPrice := int(BasePrice) + rand.Intn(int(MaxPriceVar))
		maxPrice := minPrice + int(BasePrice/2) + rand.Intn(int(MaxPriceVar/2))

		hotelID := fmt.Sprintf("QG-HOTEL-%04d", h+1)

		rooms := make([]models.QiuguoPushRoomData, NumRooms)
		for r := 0; r < NumRooms; r++ {
			rt := roomTypes[r]
			rooms[r] = models.QiuguoPushRoomData{
				SupplierRoomID: fmt.Sprintf("QG-ROOM-%04d-%02d", h+1, r+1),
				Name:           rt.name,
				Description:    rt.description,
				Capacity:       rt.capacity,
				Area:           rt.area,
				BedType:        rt.bedType,
				Amenities:      rt.amenities,
				ImageURL:       fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=modern%%20hotel%%20room%%20type%%20%d&image_size=landscape_4_3", r+1),
				TotalCount:     BaseTotal + rand.Intn(10),
			}
		}

		hotels[h] = models.QiuguoPushHotelData{
			SupplierHotelID: hotelID,
			Name:            fmt.Sprintf("%s(%s%s)", hotelType, city, addressSuffix),
			City:            city,
			Address:         fmt.Sprintf("%s市%s区%s路%d号", city, district, getRandomRoad(), 100+rand.Intn(900)),
			Description:     fmt.Sprintf("%s位于%s核心商圈，交通便利，设施完善。酒店拥有各类客房，配备现代化设施，是商务出行和休闲旅游的理想选择。", hotelType, city),
			Rating:          4.3 + float64(rand.Intn(5))/10.0,
			ImageURL:        fmt.Sprintf("https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=qiuguo%%20hotel%%20lobby%%20%d&image_size=landscape_4_3", h+1),
			PriceRange:      fmt.Sprintf("¥%d-¥%d", minPrice, maxPrice),
			Rooms:           rooms,
		}

		// 生成90天的价格库存数据
		now := time.Now()
		for d := 0; d < NumDays; d++ {
			date := now.AddDate(0, 0, d).Format("2006-01-02")
			weekday := now.AddDate(0, 0, d).Weekday()

			// 周末价格上浮
			weekendMulti := 1.0
			if weekday == time.Saturday || weekday == time.Sunday {
				weekendMulti = 1.15 + rand.Float64()*0.1
			}

			for r := 0; r < NumRooms; r++ {
				rt := roomTypes[r]
				baseRoomPrice := BasePrice * rt.priceMulti * weekendMulti

				// 随机价格波动
				priceFluctuation := 0.9 + rand.Float64()*0.2
				price := baseRoomPrice * priceFluctuation

				// 随机可用数量
				totalCount := rooms[r].TotalCount
				availableCount := rand.Intn(totalCount + 1) // 0 到 totalCount

				priceInventories = append(priceInventories, models.QiuguoPushPriceInventoryData{
					SupplierHotelID: hotelID,
					SupplierRoomID:  fmt.Sprintf("QG-ROOM-%04d-%02d", h+1, r+1),
					Date:            date,
					Price:           roundPrice(price),
					OriginalPrice:   roundPrice(price * 1.15), // 原价稍高
					AvailableCount:  availableCount,
					TotalCount:      totalCount,
				})
			}
		}
	}

	return models.QiuguoPushRequest{
		RequestID:       fmt.Sprintf("MOCK-PUSH-%s", time.Now().Format("20060102150405")),
		PushType:        "full",
		Hotels:          hotels,
		PriceInventories: priceInventories,
		StartDate:       time.Now().Format("2006-01-02"),
		EndDate:         time.Now().AddDate(0, 0, NumDays-1).Format("2006-01-02"),
	}
}

func getRandomRoad() string {
	roads := []string{"中山", "人民", "解放", "建国", "和平", "新华", "建设", "发展", "科技", "商务", "金融", "文化"}
	return roads[rand.Intn(len(roads))]
}

func roundPrice(price float64) float64 {
	return float64(int(price/10)*10)
}

func main() {
	fmt.Println("正在生成秋果集团测试数据...")
	fmt.Printf("配置: %d 家酒店, %d 种房型, %d 天价格库存\n", NumHotels, NumRooms, NumDays)
	fmt.Printf("预计生成: %d 条价格库存记录\n", NumHotels*NumRooms*NumDays)

	data := generateMockData()

	fmt.Printf("\n实际生成: \n")
	fmt.Printf("  - 酒店: %d 家\n", len(data.Hotels))
	fmt.Printf("  - 价格库存: %d 条\n", len(data.PriceInventories))
	fmt.Printf("  - 日期范围: %s 至 %s\n", data.StartDate, data.EndDate)

	// 输出一些统计信息
	hotelStats := make(map[string]int)
	for _, pi := range data.PriceInventories {
		hotelStats[pi.SupplierHotelID]++
	}

	fmt.Println("\n各酒店价格库存统计:")
	for hotelID, count := range hotelStats {
		fmt.Printf("  %s: %d 条\n", hotelID, count)
	}

	// 写入JSON文件
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("qiuguo_mock_data_%s.json", time.Now().Format("20060102"))
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ 测试数据已成功生成并保存到: %s\n", filename)
	fmt.Println("\n使用以下命令推送数据:")
	fmt.Printf("  curl -s -X POST http://localhost:8081/api/shiji/qiuguo/push \\\n")
	fmt.Printf("    -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("    -d @%s\n", filename)
}

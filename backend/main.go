package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hotel-booking/cache"
	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/routes"
	"hotel-booking/suppliers"
	"log"
)

func columnExists(db *sql.DB, tableName, columnName string) bool {
	cfg := config.GetConfig()
	
	if cfg.DBType == "sqlite" {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?
		`, tableName, columnName).Scan(&count)
		if err != nil {
			return false
		}
		return count > 0
	} else {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM information_schema.columns 
			WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?
		`, tableName, columnName).Scan(&count)
		if err != nil {
			return false
		}
		return count > 0
	}
}

func initDatabaseTables() error {
	db := database.GetDB()
	cfg := config.GetConfig()
	
	if cfg.DBType == "sqlite" {
		createUsersSQL := `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				email TEXT,
				phone TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createHotelsSQL := `
			CREATE TABLE IF NOT EXISTS hotels (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				address TEXT NOT NULL,
				city TEXT NOT NULL,
				description TEXT,
				rating REAL DEFAULT 0.0,
				image_url TEXT,
				price_range TEXT,
				supplier_id INTEGER DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createRoomsSQL := `
			CREATE TABLE IF NOT EXISTS rooms (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				hotel_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				description TEXT,
				price REAL NOT NULL,
				capacity INTEGER NOT NULL,
				area INTEGER,
				bed_type TEXT,
				amenities TEXT,
				image_url TEXT,
				total_count INTEGER NOT NULL DEFAULT 1,
				available_count INTEGER NOT NULL DEFAULT 1,
				supplier_id INTEGER DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createOrdersSQL := `
			CREATE TABLE IF NOT EXISTS orders (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				order_no TEXT NOT NULL UNIQUE,
				user_id INTEGER NOT NULL,
				hotel_id INTEGER NOT NULL,
				room_id INTEGER NOT NULL,
				check_in DATE NOT NULL,
				check_out DATE NOT NULL,
				guest_name TEXT NOT NULL,
				guest_phone TEXT NOT NULL,
				total_amount REAL NOT NULL,
				status TEXT DEFAULT 'pending',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createSuppliersSQL := `
			CREATE TABLE IF NOT EXISTS suppliers (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				code TEXT NOT NULL UNIQUE,
				description TEXT,
				api_url TEXT,
				api_key TEXT,
				secret_key TEXT,
				status TEXT DEFAULT 'active',
				priority INTEGER DEFAULT 0,
				price_control REAL DEFAULT 1.0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createSupplierHotelsSQL := `
			CREATE TABLE IF NOT EXISTS supplier_hotels (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				supplier_id INTEGER NOT NULL,
				supplier_hotel_id TEXT NOT NULL,
				local_hotel_id INTEGER NOT NULL,
				hotel_name TEXT NOT NULL,
				city TEXT,
				address TEXT,
				rating REAL DEFAULT 0.0,
				image_url TEXT,
				price_range TEXT,
				raw_data TEXT,
				status TEXT DEFAULT 'pending',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(supplier_id, supplier_hotel_id)
			)
		`
		
		createSupplierRoomsSQL := `
			CREATE TABLE IF NOT EXISTS supplier_rooms (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				supplier_id INTEGER NOT NULL,
				supplier_hotel_id TEXT NOT NULL,
				supplier_room_id TEXT NOT NULL,
				local_room_id INTEGER NOT NULL,
				room_name TEXT NOT NULL,
				description TEXT,
				price REAL NOT NULL,
				capacity INTEGER NOT NULL,
				area INTEGER,
				bed_type TEXT,
				amenities TEXT,
				image_url TEXT,
				total_count INTEGER NOT NULL DEFAULT 1,
				available_count INTEGER NOT NULL DEFAULT 1,
				raw_data TEXT,
				status TEXT DEFAULT 'pending',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(supplier_id, supplier_hotel_id, supplier_room_id)
			)
		`
		
		createPriceInventorySQL := `
			CREATE TABLE IF NOT EXISTS price_inventory (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				supplier_id INTEGER NOT NULL,
				supplier_hotel_id TEXT NOT NULL,
				supplier_room_id TEXT NOT NULL,
				date TEXT NOT NULL,
				price REAL NOT NULL DEFAULT 0.0,
				original_price REAL NOT NULL DEFAULT 0.0,
				available_count INTEGER NOT NULL DEFAULT 0,
				total_count INTEGER NOT NULL DEFAULT 0,
				status TEXT DEFAULT 'active',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(supplier_id, supplier_hotel_id, supplier_room_id, date)
			)
		`
		
		createRoomPriceSummarySQL := `
			CREATE TABLE IF NOT EXISTS room_price_summary (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				supplier_id INTEGER NOT NULL,
				supplier_hotel_id TEXT NOT NULL,
				supplier_room_id TEXT NOT NULL,
				min_price REAL NOT NULL DEFAULT 0.0,
				max_price REAL NOT NULL DEFAULT 0.0,
				avg_price REAL NOT NULL DEFAULT 0.0,
				price_range TEXT,
				has_inventory INTEGER NOT NULL DEFAULT 0,
				total_count INTEGER NOT NULL DEFAULT 0,
				date_range_start TEXT,
				date_range_end TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(supplier_id, supplier_hotel_id, supplier_room_id)
			)
		`
		
		createHotelPriceSummarySQL := `
			CREATE TABLE IF NOT EXISTS hotel_price_summary (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				supplier_id INTEGER NOT NULL,
				supplier_hotel_id TEXT NOT NULL,
				min_price REAL NOT NULL DEFAULT 0.0,
				max_price REAL NOT NULL DEFAULT 0.0,
				avg_price REAL NOT NULL DEFAULT 0.0,
				price_range TEXT,
				total_rooms INTEGER NOT NULL DEFAULT 0,
				rooms_with_inventory INTEGER NOT NULL DEFAULT 0,
				date_range_start TEXT,
				date_range_end TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(supplier_id, supplier_hotel_id)
			)
		`

		createGuestsSQL := `
			CREATE TABLE IF NOT EXISTS guests (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				phone TEXT NOT NULL,
				id_type TEXT,
				id_number TEXT,
				is_default INTEGER DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`

		db.Exec(createUsersSQL)
		db.Exec(createHotelsSQL)
		db.Exec(createRoomsSQL)
		db.Exec(createOrdersSQL)
		db.Exec(createSuppliersSQL)
		db.Exec(createSupplierHotelsSQL)
		db.Exec(createSupplierRoomsSQL)
		db.Exec(createPriceInventorySQL)
		db.Exec(createRoomPriceSummarySQL)
		db.Exec(createHotelPriceSummarySQL)
		db.Exec(createGuestsSQL)
	} else {
		createSuppliersSQL := `
			CREATE TABLE IF NOT EXISTS suppliers (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				code VARCHAR(50) NOT NULL UNIQUE,
				description TEXT,
				api_url VARCHAR(255),
				api_key VARCHAR(255),
				secret_key VARCHAR(255),
				status ENUM('active', 'inactive') DEFAULT 'active',
				priority INT DEFAULT 0,
				price_control DECIMAL(5,2) DEFAULT 1.00,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		createSupplierHotelsSQL := `
			CREATE TABLE IF NOT EXISTS supplier_hotels (
				id INT AUTO_INCREMENT PRIMARY KEY,
				supplier_id INT NOT NULL,
				supplier_hotel_id VARCHAR(100) NOT NULL,
				local_hotel_id INT NOT NULL,
				hotel_name VARCHAR(100) NOT NULL,
				city VARCHAR(50),
				address VARCHAR(255),
				rating DECIMAL(2,1) DEFAULT 0.0,
				image_url VARCHAR(255),
				price_range VARCHAR(50),
				raw_data TEXT,
				status ENUM('pending', 'synced', 'failed') DEFAULT 'pending',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				UNIQUE KEY uk_supplier_hotel (supplier_id, supplier_hotel_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		createSupplierRoomsSQL := `
			CREATE TABLE IF NOT EXISTS supplier_rooms (
				id INT AUTO_INCREMENT PRIMARY KEY,
				supplier_id INT NOT NULL,
				supplier_hotel_id VARCHAR(100) NOT NULL,
				supplier_room_id VARCHAR(100) NOT NULL,
				local_room_id INT NOT NULL,
				room_name VARCHAR(100) NOT NULL,
				description TEXT,
				price DECIMAL(10,2) NOT NULL,
				capacity INT NOT NULL,
				area INT,
				bed_type VARCHAR(50),
				amenities TEXT,
				image_url VARCHAR(255),
				total_count INT NOT NULL DEFAULT 1,
				available_count INT NOT NULL DEFAULT 1,
				raw_data TEXT,
				status ENUM('pending', 'synced', 'failed') DEFAULT 'pending',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				UNIQUE KEY uk_supplier_room (supplier_id, supplier_hotel_id, supplier_room_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		createPriceInventorySQL := `
			CREATE TABLE IF NOT EXISTS price_inventory (
				id INT AUTO_INCREMENT PRIMARY KEY,
				supplier_id INT NOT NULL,
				supplier_hotel_id VARCHAR(100) NOT NULL,
				supplier_room_id VARCHAR(100) NOT NULL,
				date DATE NOT NULL,
				price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				original_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				available_count INT NOT NULL DEFAULT 0,
				total_count INT NOT NULL DEFAULT 0,
				status ENUM('active', 'inactive') DEFAULT 'active',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				UNIQUE KEY uk_price_inventory_date (supplier_id, supplier_hotel_id, supplier_room_id, date),
				KEY idx_supplier_date (supplier_id, date),
				KEY idx_hotel_room_date (supplier_hotel_id, supplier_room_id, date)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		createRoomPriceSummarySQL := `
			CREATE TABLE IF NOT EXISTS room_price_summary (
				id INT AUTO_INCREMENT PRIMARY KEY,
				supplier_id INT NOT NULL,
				supplier_hotel_id VARCHAR(100) NOT NULL,
				supplier_room_id VARCHAR(100) NOT NULL,
				min_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				max_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				avg_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				price_range VARCHAR(50),
				has_inventory TINYINT(1) NOT NULL DEFAULT 0,
				total_count INT NOT NULL DEFAULT 0,
				date_range_start DATE,
				date_range_end DATE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				UNIQUE KEY uk_room_summary (supplier_id, supplier_hotel_id, supplier_room_id),
				KEY idx_supplier_id (supplier_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		createHotelPriceSummarySQL := `
			CREATE TABLE IF NOT EXISTS hotel_price_summary (
				id INT AUTO_INCREMENT PRIMARY KEY,
				supplier_id INT NOT NULL,
				supplier_hotel_id VARCHAR(100) NOT NULL,
				min_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				max_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				avg_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				price_range VARCHAR(50),
				total_rooms INT NOT NULL DEFAULT 0,
				rooms_with_inventory INT NOT NULL DEFAULT 0,
				date_range_start DATE,
				date_range_end DATE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				UNIQUE KEY uk_hotel_summary (supplier_id, supplier_hotel_id),
				KEY idx_supplier_id (supplier_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`
		
		db.Exec(createSuppliersSQL)
		db.Exec(createSupplierHotelsSQL)
		db.Exec(createSupplierRoomsSQL)
		db.Exec(createPriceInventorySQL)
		db.Exec(createRoomPriceSummarySQL)
		db.Exec(createHotelPriceSummarySQL)
	}
	
	migrateDatabase()
	
	return nil
}

func migrateDatabase() {
	db := database.GetDB()
	cfg := config.GetConfig()
	
	if !columnExists(db, "suppliers", "secret_key") {
		log.Println("正在添加 secret_key 字段到 suppliers 表...")
		if cfg.DBType == "sqlite" {
			db.Exec(`ALTER TABLE suppliers ADD COLUMN secret_key TEXT`)
		} else {
			db.Exec(`ALTER TABLE suppliers ADD COLUMN secret_key VARCHAR(255)`)
		}
		log.Println("secret_key 字段添加完成")
	}
}

func getSupplierPriority(code string) int {
	priorityMap := map[string]int{
		"shiji_marriott": 10,
		"shiji_hilton":   9,
		"shiji_ihg":      8,
		"shiji_wanda":    7,
		"shiji_kaiyuan":  6,
		"shiji_lvdi":     5,
		"huazhu":         4,
		"shiji_qiuguo":   3,
		"jinjiang":       3,
		"rujia":          2,
	}
	if p, ok := priorityMap[code]; ok {
		return p
	}
	return 1
}

func getSupplierPriceControl(code string) float64 {
	priceControlMap := map[string]float64{
		"shiji_marriott": 1.15,
		"shiji_hilton":   1.10,
		"shiji_ihg":      1.08,
		"shiji_wanda":    1.05,
		"shiji_kaiyuan":  1.02,
		"shiji_lvdi":     1.00,
		"huazhu":         0.95,
		"shiji_qiuguo":   0.90,
		"jinjiang":       0.92,
		"rujia":          0.88,
	}
	if p, ok := priceControlMap[code]; ok {
		return p
	}
	return 1.0
}

func initSupplierRecords() error {
	db := database.GetDB()
	adapters := suppliers.GetAllAdapters()
	
	for _, adapter := range adapters {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM suppliers WHERE code = ?", adapter.GetCode()).Scan(&count)
		if err != nil {
			log.Printf("检查供应商 %s 失败: %v", adapter.GetCode(), err)
			continue
		}
		
		if count == 0 {
			priority := getSupplierPriority(adapter.GetCode())
			priceControl := getSupplierPriceControl(adapter.GetCode())
			
			_, err = db.Exec(`
				INSERT INTO suppliers (name, code, description, api_url, status, priority, price_control)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				adapter.GetName(), adapter.GetCode(), adapter.GetDescription(),
				adapter.GetAPIURL(), "active", priority, priceControl)
			if err != nil {
				log.Printf("初始化供应商 %s 失败: %v", adapter.GetName(), err)
			} else {
				log.Printf("已注册供应商: %s (优先级: %d, 控价系数: %.2f)", adapter.GetName(), priority, priceControl)
			}
		}
	}
	
	return nil
}

func initTestUsers() error {
	db := database.GetDB()
	
	testUsers := []struct {
		username string
		password string
		email    string
		phone    string
	}{
		{"admin", "admin123", "admin@example.com", "13800138001"},
		{"testuser", "test123", "test@example.com", "13800138002"},
	}
	
	for _, u := range testUsers {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", u.username).Scan(&count)
		if err != nil {
			log.Printf("检查用户 %s 失败: %v", u.username, err)
			continue
		}
		
		if count == 0 {
			_, err = db.Exec(`
				INSERT INTO users (username, password, email, phone)
				VALUES (?, ?, ?, ?)`,
				u.username, u.password, u.email, u.phone)
			if err != nil {
				log.Printf("初始化用户 %s 失败: %v", u.username, err)
			} else {
				log.Printf("已注册测试用户: %s", u.username)
			}
		}
	}
	
	return nil
}

func initTestGuests() error {
	db := database.GetDB()

	testGuests := []struct {
		userID    int
		name      string
		phone     string
		idType    string
		idNumber  string
		isDefault bool
	}{
		{1, "张三", "13800138001", "身份证", "110101199001011234", true},
		{1, "李四", "13900139002", "身份证", "110101199002022345", false},
		{2, "王五", "13700137003", "护照", "E12345678", true},
	}

	for _, g := range testGuests {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM guests WHERE user_id = ? AND name = ?", g.userID, g.name).Scan(&count)
		if err != nil {
			log.Printf("检查入住人 %s 失败: %v", g.name, err)
			continue
		}

		if count == 0 {
			if g.isDefault {
				db.Exec("UPDATE guests SET is_default = 0 WHERE user_id = ?", g.userID)
			}

			isDefaultInt := 0
			if g.isDefault {
				isDefaultInt = 1
			}

			_, err = db.Exec(`
				INSERT INTO guests (user_id, name, phone, id_type, id_number, is_default)
				VALUES (?, ?, ?, ?, ?, ?)`,
				g.userID, g.name, g.phone, g.idType, g.idNumber, isDefaultInt)
			if err != nil {
				log.Printf("初始化入住人 %s 失败: %v", g.name, err)
			} else {
				log.Printf("已注册测试入住人: %s", g.name)
			}
		}
	}

	return nil
}

func pullAndSyncSupplier(adapter suppliers.SupplierAdapter) (int, error) {
	db := database.GetDB()
	
	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", adapter.GetCode()).Scan(&supplierID)
	if err != nil {
		return 0, fmt.Errorf("获取供应商ID失败: %v", err)
	}
	
	hotels, err := adapter.FetchHotels()
	if err != nil {
		return 0, err
	}
	
	pulledCount := 0
	
	for _, hotelData := range hotels {
		var existingID int
		err := db.QueryRow(`
			SELECT id FROM supplier_hotels WHERE supplier_id = ? AND supplier_hotel_id = ?`,
			supplierID, hotelData.SupplierHotelID).Scan(&existingID)
		
		rawData, _ := json.Marshal(hotelData)
		
		if err == sql.ErrNoRows {
			tx, err := db.Begin()
			if err != nil {
				continue
			}
			
			hotelResult, err := tx.Exec(`
				INSERT INTO hotels (name, address, city, description, rating, image_url, price_range, supplier_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				hotelData.Name, hotelData.Address, hotelData.City, hotelData.Description,
				hotelData.Rating, hotelData.ImageURL, hotelData.PriceRange, supplierID)
			if err != nil {
				tx.Rollback()
				continue
			}
			
			localHotelID, _ := hotelResult.LastInsertId()
			
			_, err = tx.Exec(`
				INSERT INTO supplier_hotels (supplier_id, supplier_hotel_id, local_hotel_id, hotel_name, 
				city, address, rating, image_url, price_range, raw_data, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				supplierID, hotelData.SupplierHotelID, localHotelID, hotelData.Name,
				hotelData.City, hotelData.Address, hotelData.Rating, hotelData.ImageURL,
				hotelData.PriceRange, string(rawData), "synced")
			if err != nil {
				tx.Rollback()
				continue
			}
			
			for _, roomData := range hotelData.Rooms {
				roomResult, err := tx.Exec(`
					INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, supplier_id)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					localHotelID, roomData.Name, roomData.Description, roomData.Price,
					roomData.Capacity, roomData.Area, roomData.BedType, roomData.Amenities,
					roomData.ImageURL, roomData.TotalCount, roomData.AvailableCount, supplierID)
				if err != nil {
					continue
				}
				
				localRoomID, _ := roomResult.LastInsertId()
				
				roomRawData, _ := json.Marshal(roomData)
				tx.Exec(`
					INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
					local_room_id, room_name, description, price, capacity, area, bed_type, 
					amenities, image_url, total_count, available_count, raw_data, status)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					supplierID, hotelData.SupplierHotelID, roomData.SupplierRoomID,
					localRoomID, roomData.Name, roomData.Description, roomData.Price,
					roomData.Capacity, roomData.Area, roomData.BedType, roomData.Amenities,
					roomData.ImageURL, roomData.TotalCount, roomData.AvailableCount,
					string(roomRawData), "synced")
			}
			
			tx.Commit()
			pulledCount++
		}
	}
	
	return pulledCount, nil
}

func autoPullAllSuppliers() {
	db := database.GetDB()
	
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM hotels WHERE supplier_id > 0").Scan(&count)
	if err != nil {
		log.Printf("检查供应商数据失败: %v", err)
		return
	}
	
	if count == 0 {
		log.Println("检测到没有供应商数据，开始自动拉取所有供应商数据...")
		
		adapters := suppliers.GetAllAdapters()
		totalPulled := 0
		
		for _, adapter := range adapters {
			log.Printf("正在拉取 %s 的数据...", adapter.GetName())
			pulled, err := pullAndSyncSupplier(adapter)
			if err != nil {
				log.Printf("拉取 %s 数据失败: %v", adapter.GetName(), err)
			} else {
				log.Printf("拉取 %s 完成，新增 %d 家酒店", adapter.GetName(), pulled)
				totalPulled += pulled
			}
		}
		
		log.Printf("所有供应商数据拉取完成，共新增 %d 家酒店", totalPulled)
	} else {
		log.Printf("已存在 %d 家供应商酒店数据", count)
	}
}

func main() {
	suppliers.InitSuppliers()
	
	err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	
	err = initDatabaseTables()
	if err != nil {
		log.Printf("数据库表初始化警告: %v", err)
	}
	
	err = initSupplierRecords()
	if err != nil {
		log.Printf("供应商记录初始化警告: %v", err)
	}
	
	err = initTestUsers()
	if err != nil {
		log.Printf("测试用户初始化警告: %v", err)
	}

	err = initTestGuests()
	if err != nil {
		log.Printf("测试入住人初始化警告: %v", err)
	}
	
	cfg := config.GetConfig()
	if cfg.EnableRedisCache {
		err = cache.InitRedis()
		if err != nil {
			log.Printf("Redis连接失败，将禁用缓存功能: %v", err)
			cfg.EnableRedisCache = false
		}
	}
	
	autoPullAllSuppliers()

	r := routes.SetupRouter()

	log.Printf("服务器启动在端口 %d...", cfg.ServerPort)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

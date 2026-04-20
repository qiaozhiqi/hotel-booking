package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
				min_price REAL DEFAULT 0.0,
				max_price REAL DEFAULT 0.0,
				supplier_id INTEGER DEFAULT 0,
				supplier_code TEXT,
				supplier_name TEXT,
				brand TEXT,
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
				original_price REAL DEFAULT 0.0,
				capacity INTEGER NOT NULL,
				area INTEGER,
				bed_type TEXT,
				amenities TEXT,
				image_url TEXT,
				total_count INTEGER NOT NULL DEFAULT 1,
				available_count INTEGER NOT NULL DEFAULT 1,
				supplier_id INTEGER DEFAULT 0,
				supplier_code TEXT,
				supplier_name TEXT,
				is_price_controlled INTEGER DEFAULT 0,
				price_control_reason TEXT,
				promotion_tag TEXT,
				payment_type TEXT,
				cancel_policy TEXT,
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
				status TEXT DEFAULT 'active',
				priority INTEGER DEFAULT 0,
				color TEXT,
				icon TEXT,
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
		
		db.Exec(createUsersSQL)
		db.Exec(createHotelsSQL)
		db.Exec(createRoomsSQL)
		db.Exec(createOrdersSQL)
		db.Exec(createSuppliersSQL)
		db.Exec(createSupplierHotelsSQL)
		db.Exec(createSupplierRoomsSQL)
	} else {
		createSuppliersSQL := `
			CREATE TABLE IF NOT EXISTS suppliers (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				code VARCHAR(50) NOT NULL UNIQUE,
				description TEXT,
				api_url VARCHAR(255),
				api_key VARCHAR(255),
				status ENUM('active', 'inactive') DEFAULT 'active',
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
		
		db.Exec(createSuppliersSQL)
		db.Exec(createSupplierHotelsSQL)
		db.Exec(createSupplierRoomsSQL)
	}
	
	return nil
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
			_, err = db.Exec(`
				INSERT INTO suppliers (name, code, description, api_url, status, priority, color, icon)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				adapter.GetName(), adapter.GetCode(), adapter.GetDescription(),
				adapter.GetAPIURL(), "active", adapter.GetPriority(),
				adapter.GetColor(), adapter.GetIcon())
			if err != nil {
				log.Printf("初始化供应商 %s 失败: %v", adapter.GetName(), err)
			} else {
				log.Printf("已注册供应商: %s (优先级: %d)", adapter.GetName(), adapter.GetPriority())
			}
		} else {
			_, err = db.Exec(`
				UPDATE suppliers SET priority = ?, color = ?, icon = ? WHERE code = ?`,
				adapter.GetPriority(), adapter.GetColor(), adapter.GetIcon(), adapter.GetCode())
			if err != nil {
				log.Printf("更新供应商 %s 信息失败: %v", adapter.GetName(), err)
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
	
	supplierCode := adapter.GetCode()
	supplierName := adapter.GetName()
	
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
				INSERT INTO hotels (name, address, city, description, rating, image_url, 
				price_range, min_price, max_price, supplier_id, supplier_code, supplier_name, brand)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				hotelData.Name, hotelData.Address, hotelData.City, hotelData.Description,
				hotelData.Rating, hotelData.ImageURL, hotelData.PriceRange,
				hotelData.MinPrice, hotelData.MaxPrice, supplierID,
				supplierCode, supplierName, hotelData.Brand)
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
				isPriceControlled := 0
				if roomData.IsPriceControlled {
					isPriceControlled = 1
				}
				
				roomResult, err := tx.Exec(`
					INSERT INTO rooms (hotel_id, name, description, price, original_price, 
					capacity, area, bed_type, amenities, image_url, total_count, available_count, 
					supplier_id, supplier_code, supplier_name, is_price_controlled, 
					price_control_reason, promotion_tag, payment_type, cancel_policy)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					localHotelID, roomData.Name, roomData.Description, roomData.Price,
					roomData.OriginalPrice, roomData.Capacity, roomData.Area, roomData.BedType,
					roomData.Amenities, roomData.ImageURL, roomData.TotalCount, roomData.AvailableCount,
					supplierID, supplierCode, supplierName, isPriceControlled,
					roomData.PriceControlReason, roomData.PromotionTag, roomData.PaymentType,
					roomData.CancelPolicy)
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
	
	autoPullAllSuppliers()

	cfg := config.GetConfig()
	r := routes.SetupRouter()

	log.Printf("服务器启动在端口 %d...", cfg.ServerPort)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

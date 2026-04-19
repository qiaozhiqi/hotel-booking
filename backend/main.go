package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hotel-booking/config"
	"hotel-booking/controllers"
	"hotel-booking/database"
	"hotel-booking/routes"
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
	
	var createSuppliersSQL string
	var createSupplierHotelsSQL string
	var createSupplierRoomsSQL string
	
	if cfg.DBType == "sqlite" {
		createSuppliersSQL = `
			CREATE TABLE IF NOT EXISTS suppliers (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				code TEXT NOT NULL UNIQUE,
				description TEXT,
				api_url TEXT,
				api_key TEXT,
				status TEXT DEFAULT 'active',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
		
		createSupplierHotelsSQL = `
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
		
		createSupplierRoomsSQL = `
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
	} else {
		createSuppliersSQL = `
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
		
		createSupplierHotelsSQL = `
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
		
		createSupplierRoomsSQL = `
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
	}
	
	_, err := db.Exec(createSuppliersSQL)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(createSupplierHotelsSQL)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(createSupplierRoomsSQL)
	if err != nil {
		return err
	}
	
	if !columnExists(db, "hotels", "supplier_id") {
		_, err = db.Exec("ALTER TABLE hotels ADD COLUMN supplier_id INTEGER DEFAULT 0")
		if err != nil {
			log.Printf("添加列 supplier_id 到 hotels 表失败: %v", err)
		}
	}
	
	if !columnExists(db, "rooms", "supplier_id") {
		_, err = db.Exec("ALTER TABLE rooms ADD COLUMN supplier_id INTEGER DEFAULT 0")
		if err != nil {
			log.Printf("添加列 supplier_id 到 rooms 表失败: %v", err)
		}
	}
	
	return nil
}

func autoPullSupplierData() {
	db := database.GetDB()
	
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM hotels WHERE supplier_id > 0").Scan(&count)
	if err != nil {
		log.Printf("检查供应商数据失败: %v", err)
		return
	}
	
	if count == 0 {
		log.Println("检测到没有供应商数据，开始自动拉取华住酒店数据...")
		
		mockHotels := controllers.GetHuazhuMockHotels()
		
		var supplierID int
		err = db.QueryRow("SELECT id FROM suppliers WHERE code = ?", "huazhu").Scan(&supplierID)
		if err != nil {
			log.Printf("获取供应商ID失败: %v", err)
			return
		}
		
		pulledCount := 0
		for _, mockHotel := range mockHotels {
			var existingID int
			err := db.QueryRow(`
				SELECT id FROM supplier_hotels WHERE supplier_id = ? AND supplier_hotel_id = ?`,
				supplierID, mockHotel.SupplierHotelID).Scan(&existingID)
			
			rawData, _ := json.Marshal(mockHotel)
			
			if err == sql.ErrNoRows {
				tx, err := db.Begin()
				if err != nil {
					continue
				}
				
				hotelResult, err := tx.Exec(`
					INSERT INTO hotels (name, address, city, description, rating, image_url, price_range, supplier_id)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
					mockHotel.Name, mockHotel.Address, mockHotel.City, mockHotel.Description,
					mockHotel.Rating, mockHotel.ImageURL, mockHotel.PriceRange, supplierID)
				if err != nil {
					tx.Rollback()
					continue
				}
				
				localHotelID, _ := hotelResult.LastInsertId()
				
				_, err = tx.Exec(`
					INSERT INTO supplier_hotels (supplier_id, supplier_hotel_id, local_hotel_id, hotel_name, 
					city, address, rating, image_url, price_range, raw_data, status)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					supplierID, mockHotel.SupplierHotelID, localHotelID, mockHotel.Name,
					mockHotel.City, mockHotel.Address, mockHotel.Rating, mockHotel.ImageURL,
					mockHotel.PriceRange, string(rawData), "synced")
				if err != nil {
					tx.Rollback()
					continue
				}
				
				for _, mockRoom := range mockHotel.Rooms {
					roomResult, err := tx.Exec(`
						INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, 
						amenities, image_url, total_count, available_count, supplier_id)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
						localHotelID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
						mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
						mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount, supplierID)
					if err != nil {
						continue
					}
					
					localRoomID, _ := roomResult.LastInsertId()
					
					roomRawData, _ := json.Marshal(mockRoom)
					tx.Exec(`
						INSERT INTO supplier_rooms (supplier_id, supplier_hotel_id, supplier_room_id, 
						local_room_id, room_name, description, price, capacity, area, bed_type, 
						amenities, image_url, total_count, available_count, raw_data, status)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
						supplierID, mockHotel.SupplierHotelID, mockRoom.SupplierRoomID,
						localRoomID, mockRoom.Name, mockRoom.Description, mockRoom.Price,
						mockRoom.Capacity, mockRoom.Area, mockRoom.BedType, mockRoom.Amenities,
						mockRoom.ImageURL, mockRoom.TotalCount, mockRoom.AvailableCount,
						string(roomRawData), "synced")
				}
				
				tx.Commit()
				pulledCount++
			}
		}
		
		log.Printf("自动拉取完成，新增 %d 家酒店数据", pulledCount)
	} else {
		log.Printf("已存在 %d 家供应商酒店数据", count)
	}
}

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	
	err = initDatabaseTables()
	if err != nil {
		log.Printf("数据库表初始化警告: %v", err)
	}
	
	err = controllers.InitSupplier()
	if err != nil {
		log.Printf("供应商初始化警告: %v", err)
	}
	
	autoPullSupplierData()

	cfg := config.GetConfig()
	r := routes.SetupRouter()

	log.Printf("服务器启动在端口 %d...", cfg.ServerPort)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

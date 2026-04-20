package database

import (
	"database/sql"
	"hotel-booking/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	cfg := config.GetConfig()
	var err error
	var driverName string
	
	if cfg.DBType == "sqlite" {
		driverName = "sqlite3"
	} else {
		driverName = "mysql"
	}
	
	DB, err = sql.Open(driverName, cfg.GetDSN())
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	if cfg.DBType == "sqlite" {
		log.Println("SQLite数据库连接成功")
	} else {
		log.Println("MySQL数据库连接成功")
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}

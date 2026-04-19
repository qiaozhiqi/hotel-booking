package database

import (
	"database/sql"
	"hotel-booking/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	cfg := config.GetConfig()
	var err error
	DB, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("MySQL数据库连接成功")
	return nil
}

func GetDB() *sql.DB {
	return DB
}

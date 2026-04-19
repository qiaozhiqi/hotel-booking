package config

import (
	"fmt"
)

type Config struct {
	DBType     string
	DBPath     string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
}

func GetConfig() *Config {
	return &Config{
		DBType:     "sqlite",
		DBPath:     "./database/hotel_booking.db",
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "root",
		DBPassword: "root",
		DBName:     "hotel_booking",
		ServerPort: 8081,
	}
}

func (c *Config) GetDSN() string {
	if c.DBType == "sqlite" {
		return c.DBPath
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

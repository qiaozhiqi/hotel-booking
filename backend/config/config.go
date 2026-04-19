package config

import (
	"fmt"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
}

func GetConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "root",
		DBPassword: "root",
		DBName:     "hotel_booking",
		ServerPort: 8080,
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

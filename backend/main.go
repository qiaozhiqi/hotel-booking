package main

import (
	"fmt"
	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/routes"
	"log"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	cfg := config.GetConfig()
	r := routes.SetupRouter()

	log.Printf("服务器启动在端口 %d...", cfg.ServerPort)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

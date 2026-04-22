//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"hotel-booking/database"
	"hotel-booking/security"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run init_credentials.go <supplier_code>")
		fmt.Println("示例: go run init_credentials.go shiji_qiuguo")
		os.Exit(1)
	}

	supplierCode := os.Args[1]

	err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	credentials, err := security.GenerateTestCredentials(supplierCode)
	if err != nil {
		log.Fatalf("生成凭证失败: %v", err)
	}

	fmt.Println("\n========================================")
	fmt.Printf("供应商: %s\n", supplierCode)
	fmt.Println("========================================")
	fmt.Printf("API Key:    %s\n", credentials.APIKey)
	fmt.Printf("Secret Key: %s\n", credentials.SecretKey)
	fmt.Println("========================================")
	fmt.Println("\n请妥善保管以上凭证，不要泄露给第三方！")
	fmt.Println("\n推送接口: POST http://localhost:8081/api/shiji/qiuguo/push")
	fmt.Println("\n请求头:")
	fmt.Printf("  X-API-Key:    %s\n", credentials.APIKey)
	fmt.Printf("  X-Signature:  <HMAC-SHA256签名>\n")
	fmt.Printf("  X-Timestamp:  <Unix时间戳秒数>\n")
	fmt.Printf("  X-Request-ID: <唯一请求ID>\n")
}

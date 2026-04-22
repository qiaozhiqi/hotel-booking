//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-booking/security"
	"io"
	"net/http"
	"time"
)

func main() {
	apiKey := "17768509324904480004904480008000"
	secretKey := "1776850932490450000490450000017768509324904500004904500000177685"
	baseURL := "http://localhost:8081"

	testData := map[string]interface{}{
		"request_id": fmt.Sprintf("TEST-SEC-%d", time.Now().Unix()),
		"push_type":  "test",
		"price_inventories": []map[string]interface{}{
			{
				"hotel_id":        "QG-HOTEL-0001",
				"room_id":         "QG-ROOM-0001-01",
				"date":            time.Now().Format("2006-01-02"),
				"price":           280.0,
				"available_count": 5,
			},
		},
	}

	fmt.Println("========================================")
	fmt.Println("安全验证测试")
	fmt.Println("========================================")

	fmt.Println("\n--- 测试1: 不携带安全凭证 ---")
	testWithoutCredentials(baseURL, testData)

	fmt.Println("\n--- 测试2: 使用正确的安全凭证 ---")
	testWithCorrectCredentials(baseURL, apiKey, secretKey, testData)

	fmt.Println("\n--- 测试3: 使用错误的签名 ---")
	testWithWrongSignature(baseURL, apiKey, testData)

	fmt.Println("\n--- 测试4: 使用过期的时间戳 ---")
	testWithExpiredTimestamp(baseURL, apiKey, secretKey, testData)

	fmt.Println("\n========================================")
	fmt.Println("测试完成")
	fmt.Println("========================================")
}

func testWithoutCredentials(baseURL string, testData map[string]interface{}) {
	bodyBytes, _ := json.Marshal(testData)

	req, _ := http.NewRequest("POST", baseURL+"/api/shiji/qiuguo/push", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(respBody))
	fmt.Printf("期望: 401 未授权 ✓\n")
}

func testWithCorrectCredentials(baseURL string, apiKey, secretKey string, testData map[string]interface{}) {
	bodyBytes, _ := json.Marshal(testData)
	timestamp := time.Now().Unix()
	requestID := fmt.Sprintf("TEST-CORRECT-%d", timestamp)

	testData["request_id"] = requestID
	bodyBytes, _ = json.Marshal(testData)

	signature := security.GenerateSignature(
		secretKey,
		"POST",
		"/api/shiji/qiuguo/push",
		timestamp,
		requestID,
		bodyBytes,
	)

	req, _ := http.NewRequest("POST", baseURL+"/api/shiji/qiuguo/push", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(respBody))
	fmt.Printf("期望: 200 成功 ✓\n")
}

func testWithWrongSignature(baseURL string, apiKey string, testData map[string]interface{}) {
	bodyBytes, _ := json.Marshal(testData)
	timestamp := time.Now().Unix()
	requestID := fmt.Sprintf("TEST-WRONG-%d", timestamp)

	testData["request_id"] = requestID
	bodyBytes, _ = json.Marshal(testData)

	wrongSignature := "this-is-a-wrong-signature-1234567890"

	req, _ := http.NewRequest("POST", baseURL+"/api/shiji/qiuguo/push", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Signature", wrongSignature)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(respBody))
	fmt.Printf("期望: 401 签名验证失败 ✓\n")
}

func testWithExpiredTimestamp(baseURL string, apiKey, secretKey string, testData map[string]interface{}) {
	bodyBytes, _ := json.Marshal(testData)
	expiredTimestamp := time.Now().Unix() - 600 // 10分钟前
	requestID := fmt.Sprintf("TEST-EXPIRED-%d", expiredTimestamp)

	testData["request_id"] = requestID
	bodyBytes, _ = json.Marshal(testData)

	signature := security.GenerateSignature(
		secretKey,
		"POST",
		"/api/shiji/qiuguo/push",
		expiredTimestamp,
		requestID,
		bodyBytes,
	)

	req, _ := http.NewRequest("POST", baseURL+"/api/shiji/qiuguo/push", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", expiredTimestamp))
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应: %s\n", string(respBody))
	fmt.Printf("期望: 403 请求已过期 ✓\n")
}

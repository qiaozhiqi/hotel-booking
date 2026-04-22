package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hotel-booking/database"
	"hotel-booking/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	HeaderAPIKey       = "X-API-Key"
	HeaderSignature    = "X-Signature"
	HeaderTimestamp    = "X-Timestamp"
	HeaderRequestID    = "X-Request-ID"

	MaxTimeDriftSeconds = 300
	RequestIDExpireHours = 24
)

var (
	requestIDCache = sync.Map{}
)

type CachedRequestID struct {
	Timestamp time.Time
}

type SupplierCredentials struct {
	APIKey    string
	SecretKey string
}

func GetSupplierCredentials(supplierCode string) (*SupplierCredentials, error) {
	db := database.GetDB()

	var apiKey, secretKey string
	err := db.QueryRow(`
		SELECT api_key, secret_key FROM suppliers WHERE code = ?`,
		supplierCode).Scan(&apiKey, &secretKey)

	if err != nil {
		return nil, err
	}

	return &SupplierCredentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}, nil
}

func GenerateSignature(secretKey string, method string, path string, timestamp int64, requestID string, body []byte) string {
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s",
		strings.ToUpper(method),
		path,
		timestamp,
		requestID,
		string(body))

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifySignature(secretKey string, method string, path string, timestamp int64, requestID string, body []byte, signature string) bool {
	expectedSignature := GenerateSignature(secretKey, method, path, timestamp, requestID, body)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func IsRequestIDDuplicate(requestID string) bool {
	if requestID == "" {
		return false
	}

	_, exists := requestIDCache.Load(requestID)
	if exists {
		return true
	}

	requestIDCache.Store(requestID, CachedRequestID{Timestamp: time.Now()})
	return false
}

func CleanupExpiredRequestIDs() {
	expireTime := time.Now().Add(-RequestIDExpireHours * time.Hour)

	requestIDCache.Range(func(key, value interface{}) bool {
		if cached, ok := value.(CachedRequestID); ok {
			if cached.Timestamp.Before(expireTime) {
				requestIDCache.Delete(key)
			}
		}
		return true
	})
}

func SupplierPushAuthMiddleware(supplierCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(HeaderAPIKey)
		signature := c.GetHeader(HeaderSignature)
		timestampStr := c.GetHeader(HeaderTimestamp)
		requestID := c.GetHeader(HeaderRequestID)

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      401,
				Message:   "缺少API Key",
			})
			c.Abort()
			return
		}

		if signature == "" {
			c.JSON(http.StatusUnauthorized, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      401,
				Message:   "缺少请求签名",
			})
			c.Abort()
			return
		}

		if timestampStr == "" {
			c.JSON(http.StatusUnauthorized, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      401,
				Message:   "缺少时间戳",
			})
			c.Abort()
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      400,
				Message:   "时间戳格式错误",
			})
			c.Abort()
			return
		}

		now := time.Now().Unix()
		timeDiff := now - timestamp
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > MaxTimeDriftSeconds {
			c.JSON(http.StatusForbidden, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      403,
				Message:   "请求已过期",
			})
			c.Abort()
			return
		}

		if IsRequestIDDuplicate(requestID) {
			c.JSON(http.StatusConflict, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      409,
				Message:   "重复请求",
			})
			c.Abort()
			return
		}

		credentials, err := GetSupplierCredentials(supplierCode)
		if err != nil {
			log.Printf("获取供应商凭证失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      500,
				Message:   "内部服务器错误",
			})
			c.Abort()
			return
		}

		if credentials.APIKey != apiKey {
			c.JSON(http.StatusUnauthorized, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      401,
				Message:   "API Key无效",
			})
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      400,
				Message:   "读取请求体失败",
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		method := c.Request.Method
		path := c.Request.URL.Path

		if !VerifySignature(credentials.SecretKey, method, path, timestamp, requestID, body, signature) {
			c.JSON(http.StatusUnauthorized, models.QiuguoPushResponse{
				RequestID: requestID,
				Success:   false,
				Code:      401,
				Message:   "签名验证失败",
			})
			c.Abort()
			return
		}

		c.Set("supplier_code", supplierCode)
		c.Set("request_id", requestID)
		c.Set("api_key", apiKey)

		c.Next()
	}
}

func GenerateTestCredentials(supplierCode string) (*SupplierCredentials, error) {
	apiKey := generateRandomKey(32)
	secretKey := generateRandomKey(64)

	db := database.GetDB()

	var existingID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&existingID)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		UPDATE suppliers SET api_key = ?, secret_key = ? WHERE code = ?`,
		apiKey, secretKey, supplierCode)

	if err != nil {
		return nil, err
	}

	log.Printf("已为供应商 %s 生成凭证:\n  API Key: %s\n  Secret Key: %s",
		supplierCode, apiKey, secretKey)

	return &SupplierCredentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}, nil
}

func generateRandomKey(length int) string {
	timestamp := time.Now().UnixNano()
	randomPart := fmt.Sprintf("%d%d%d", timestamp, time.Now().Nanosecond(), timestamp%10000)
	for len(randomPart) < length {
		randomPart += randomPart
	}
	return randomPart[:length]
}

func CreateTestPushRequest(credentials *SupplierCredentials, method string, path string, requestBody interface{}) (map[string]string, []byte, error) {
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, err
	}

	timestamp := time.Now().Unix()
	requestID := fmt.Sprintf("TEST-%d", timestamp)

	signature := GenerateSignature(
		credentials.SecretKey,
		method,
		path,
		timestamp,
		requestID,
		bodyBytes,
	)

	headers := map[string]string{
		HeaderAPIKey:    credentials.APIKey,
		HeaderSignature: signature,
		HeaderTimestamp: fmt.Sprintf("%d", timestamp),
		HeaderRequestID: requestID,
	}

	return headers, bodyBytes, nil
}

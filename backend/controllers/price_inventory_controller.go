package controllers

import (
	"hotel-booking/database"
	"hotel-booking/models"
	"hotel-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPriceInventory(c *gin.Context) {
	supplierIDStr := c.Query("supplier_id")
	supplierHotelID := c.Query("supplier_hotel_id")
	supplierRoomID := c.Query("supplier_room_id")
	date := c.Query("date")

	if supplierIDStr == "" || supplierHotelID == "" || supplierRoomID == "" || date == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少必要参数: supplier_id, supplier_hotel_id, supplier_room_id, date",
		})
		return
	}

	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "supplier_id参数无效",
		})
		return
	}

	piService := services.NewPriceInventoryService()
	pi, err := piService.GetPriceInventory(supplierID, supplierHotelID, supplierRoomID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取价格库存失败: " + err.Error(),
		})
		return
	}

	if pi == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "价格库存不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    pi,
	})
}

func GetPriceInventoryByDateRange(c *gin.Context) {
	supplierIDStr := c.Query("supplier_id")
	supplierHotelID := c.Query("supplier_hotel_id")
	supplierRoomID := c.Query("supplier_room_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if supplierIDStr == "" || supplierHotelID == "" || supplierRoomID == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少必要参数: supplier_id, supplier_hotel_id, supplier_room_id, start_date, end_date",
		})
		return
	}

	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "supplier_id参数无效",
		})
		return
	}

	piService := services.NewPriceInventoryService()
	priceInventories, err := piService.GetPriceInventoryByDateRange(supplierID, supplierHotelID, supplierRoomID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取价格库存失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"count":            len(priceInventories),
			"price_inventories": priceInventories,
			"start_date":       startDate,
			"end_date":         endDate,
		},
	})
}

func GetRoomCurrentPrice(c *gin.Context) {
	supplierIDStr := c.Query("supplier_id")
	supplierHotelID := c.Query("supplier_hotel_id")
	supplierRoomID := c.Query("supplier_room_id")

	if supplierIDStr == "" || supplierHotelID == "" || supplierRoomID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少必要参数: supplier_id, supplier_hotel_id, supplier_room_id",
		})
		return
	}

	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "supplier_id参数无效",
		})
		return
	}

	piService := services.NewPriceInventoryService()
	rcp, err := piService.GetRoomCurrentPrice(supplierID, supplierHotelID, supplierRoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取房间当前价格失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    rcp,
	})
}

func GetSupplierPriceInventorySummary(c *gin.Context) {
	supplierCode := c.Param("code")

	db := database.GetDB()

	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&supplierID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "供应商不存在",
		})
		return
	}

	piService := services.NewPriceInventoryService()
	summary, err := piService.GetSyncStatus(supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取价格库存摘要失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    summary,
	})
}

func ClearPriceInventory(c *gin.Context) {
	supplierCode := c.Param("code")

	db := database.GetDB()

	var supplierID int
	err := db.QueryRow("SELECT id FROM suppliers WHERE code = ?", supplierCode).Scan(&supplierID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "供应商不存在",
		})
		return
	}

	piService := services.NewPriceInventoryService()
	err = piService.ClearPriceInventoryBySupplier(supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "清除价格库存失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "价格库存已清除",
	})
}

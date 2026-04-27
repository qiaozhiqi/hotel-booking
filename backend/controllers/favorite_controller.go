package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateFavorite(c *gin.Context) {
	var req models.CreateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = ? AND hotel_id = ?", userID, req.HotelID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "检查收藏状态失败",
		})
		return
	}

	if exists > 0 {
		c.JSON(http.StatusOK, models.Response{
			Code:    200,
			Message: "该酒店已收藏",
		})
		return
	}

	var hotel models.Hotel
	err = db.QueryRow(`
		SELECT id, name, city, address, rating, image_url, price_range 
		FROM hotels WHERE id = ?`, req.HotelID).Scan(
		&hotel.ID, &hotel.Name, &hotel.City, &hotel.Address,
		&hotel.Rating, &hotel.ImageURL, &hotel.PriceRange)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "酒店不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询酒店信息失败",
		})
		return
	}

	result, err := db.Exec(`
		INSERT INTO favorites (user_id, hotel_id, hotel_name, city, address, rating, image_url, price_range)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, hotel.ID, hotel.Name, hotel.City, hotel.Address,
		hotel.Rating, hotel.ImageURL, hotel.PriceRange)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "收藏酒店失败",
		})
		return
	}

	favoriteID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "收藏成功",
		Data: map[string]interface{}{
			"favorite_id": favoriteID,
		},
	})
}

func DeleteFavorite(c *gin.Context) {
	hotelID := c.Param("hotel_id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少酒店ID",
		})
		return
	}

	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	result, err := db.Exec("DELETE FROM favorites WHERE user_id = ? AND hotel_id = ?", userID, hotelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "取消收藏失败",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "收藏记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "取消收藏成功",
	})
}

func GetFavoriteList(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	db := database.GetDB()

	rows, err := db.Query(`
		SELECT id, user_id, hotel_id, hotel_name, city, address, rating, image_url, price_range, created_at
		FROM favorites WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, pageSize, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取收藏列表失败",
		})
		return
	}
	defer rows.Close()

	var favorites []models.Favorite
	for rows.Next() {
		var favorite models.Favorite
		err := rows.Scan(
			&favorite.ID, &favorite.UserID, &favorite.HotelID, &favorite.HotelName,
			&favorite.City, &favorite.Address, &favorite.Rating, &favorite.ImageURL,
			&favorite.PriceRange, &favorite.CreatedAt)
		if err != nil {
			continue
		}
		favorites = append(favorites, favorite)
	}

	var total int
	db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = ?", userID).Scan(&total)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"favorites": favorites,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func CheckFavoriteStatus(c *gin.Context) {
	hotelID := c.Param("hotel_id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少酒店ID",
		})
		return
	}

	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = ? AND hotel_id = ?", userID, hotelID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "检查收藏状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]bool{
			"is_favorite": exists > 0,
		},
	})
}

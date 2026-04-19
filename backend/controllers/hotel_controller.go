package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHotelList(c *gin.Context) {
	city := c.Query("city")
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

	var rows *sql.Rows
	var err error
	var countRows *sql.Rows

	if city != "" {
		rows, err = db.Query(`
			SELECT id, name, address, city, description, rating, image_url, price_range 
			FROM hotels 
			WHERE city = ? 
			ORDER BY rating DESC 
			LIMIT ? OFFSET ?`, city, pageSize, offset)
		
		countRows, _ = db.Query("SELECT COUNT(*) FROM hotels WHERE city = ?", city)
	} else {
		rows, err = db.Query(`
			SELECT id, name, address, city, description, rating, image_url, price_range 
			FROM hotels 
			ORDER BY rating DESC 
			LIMIT ? OFFSET ?`, pageSize, offset)
		
		countRows, _ = db.Query("SELECT COUNT(*) FROM hotels")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取酒店列表失败",
		})
		return
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var hotel models.Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
			&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange)
		if err != nil {
			continue
		}
		hotels = append(hotels, hotel)
	}

	var total int
	if countRows != nil {
		defer countRows.Close()
		if countRows.Next() {
			countRows.Scan(&total)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"hotels":    hotels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetHotelDetail(c *gin.Context) {
	hotelID := c.Param("id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少酒店ID",
		})
		return
	}

	db := database.GetDB()
	var hotel models.Hotel
	err := db.QueryRow(`
		SELECT id, name, address, city, description, rating, image_url, price_range 
		FROM hotels WHERE id = ?`, hotelID).Scan(
		&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, 
		&hotel.Description, &hotel.Rating, &hotel.ImageURL, &hotel.PriceRange)

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
			Message: "获取酒店详情失败",
		})
		return
	}

	rows, err := db.Query(`
		SELECT id, hotel_id, name, description, price, capacity, area, bed_type, 
		amenities, image_url, total_count, available_count 
		FROM rooms WHERE hotel_id = ?`, hotelID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取房型列表失败",
		})
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.HotelID, &room.Name, &room.Description, &room.Price,
			&room.Capacity, &room.Area, &room.BedType, &room.Amenities, &room.ImageURL,
			&room.TotalCount, &room.AvailableCount)
		if err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"hotel": hotel,
			"rooms": rooms,
		},
	})
}

func GetCities(c *gin.Context) {
	db := database.GetDB()
	rows, err := db.Query("SELECT DISTINCT city FROM hotels ORDER BY city")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取城市列表失败",
		})
		return
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err == nil {
			cities = append(cities, city)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    cities,
	})
}

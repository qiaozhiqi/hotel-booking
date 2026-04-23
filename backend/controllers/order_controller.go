package controllers

import (
	"database/sql"
	"fmt"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
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

	if req.GuestID > 0 {
		db := database.GetDB()
		var guestName, guestPhone string
		err := db.QueryRow(`
			SELECT name, phone FROM guests WHERE id = ? AND user_id = ?`,
			req.GuestID, userID).Scan(&guestName, &guestPhone)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, models.Response{
					Code:    400,
					Message: "选择的入住人不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "查询入住人信息失败",
			})
			return
		}
		req.GuestName = guestName
		req.GuestPhone = guestPhone
	}

	if req.GuestName == "" || req.GuestPhone == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请提供入住人信息",
		})
		return
	}

	db := database.GetDB()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "创建订单失败",
		})
		return
	}

	var room models.Room
	err = tx.QueryRow(`
		SELECT id, hotel_id, name, price, available_count FROM rooms WHERE id = ? AND hotel_id = ?`,
		req.RoomID, req.HotelID).Scan(&room.ID, &room.HotelID, &room.Name, &room.Price, &room.AvailableCount)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "房型不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询房型失败",
		})
		return
	}

	if room.AvailableCount <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "该房型暂无可用房间",
		})
		return
	}

	checkIn, _ := time.Parse("2006-01-02", req.CheckIn)
	checkOut, _ := time.Parse("2006-01-02", req.CheckOut)
	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	if nights <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "离店日期必须晚于入住日期",
		})
		return
	}

	totalAmount := room.Price * float64(nights)

	orderNo := fmt.Sprintf("HB%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)

	result, err := tx.Exec(`
		INSERT INTO orders (order_no, user_id, hotel_id, room_id, check_in, check_out, 
		guest_name, guest_phone, total_amount, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		orderNo, userID, req.HotelID, req.RoomID, req.CheckIn, req.CheckOut,
		req.GuestName, req.GuestPhone, totalAmount, "confirmed")

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "创建订单失败",
		})
		return
	}

	_, err = tx.Exec("UPDATE rooms SET available_count = available_count - 1 WHERE id = ?", req.RoomID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "更新房间库存失败",
		})
		return
	}

	tx.Commit()

	orderID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "订单创建成功",
		Data: map[string]interface{}{
			"order_id":     orderID,
			"order_no":     orderNo,
			"total_amount": totalAmount,
		},
	})
}

func GetOrderList(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	status := c.Query("status")
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
	var countRows *sql.Rows

	baseQuery := `
		SELECT o.id, o.order_no, o.user_id, o.hotel_id, o.room_id, o.check_in, o.check_out,
		o.guest_name, o.guest_phone, o.total_amount, o.status, o.created_at,
		h.name as hotel_name, r.name as room_name
		FROM orders o
		LEFT JOIN hotels h ON o.hotel_id = h.id
		LEFT JOIN rooms r ON o.room_id = r.id
		WHERE o.user_id = ?`

	countQuery := "SELECT COUNT(*) FROM orders WHERE user_id = ?"

	if status != "" {
		baseQuery += " AND o.status = ?"
		countQuery += " AND status = ?"
		baseQuery += " ORDER BY o.created_at DESC LIMIT ? OFFSET ?"
		rows, _ = db.Query(baseQuery, userID, status, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID, status)
	} else {
		baseQuery += " ORDER BY o.created_at DESC LIMIT ? OFFSET ?"
		rows, _ = db.Query(baseQuery, userID, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID)
	}

	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.OrderNo, &order.UserID, &order.HotelID, &order.RoomID,
			&order.CheckIn, &order.CheckOut, &order.GuestName, &order.GuestPhone,
			&order.TotalAmount, &order.Status, &order.CreatedAt, &order.HotelName, &order.RoomName)
		if err != nil {
			continue
		}
		orders = append(orders, order)
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
			"orders":    orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetOrderDetail(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少订单ID",
		})
		return
	}

	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 1
	}

	db := database.GetDB()
	var order models.Order
	err := db.QueryRow(`
		SELECT o.id, o.order_no, o.user_id, o.hotel_id, o.room_id, o.check_in, o.check_out,
		o.guest_name, o.guest_phone, o.total_amount, o.status, o.created_at,
		h.name as hotel_name, r.name as room_name, h.address as hotel_address, h.city as hotel_city
		FROM orders o
		LEFT JOIN hotels h ON o.hotel_id = h.id
		LEFT JOIN rooms r ON o.room_id = r.id
		WHERE o.id = ? AND o.user_id = ?`, orderID, userID).Scan(
		&order.ID, &order.OrderNo, &order.UserID, &order.HotelID, &order.RoomID,
		&order.CheckIn, &order.CheckOut, &order.GuestName, &order.GuestPhone,
		&order.TotalAmount, &order.Status, &order.CreatedAt, &order.HotelName, &order.RoomName)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "订单不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取订单详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    order,
	})
}

func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少订单ID",
		})
		return
	}

	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "取消订单失败",
		})
		return
	}

	var roomID int
	var status string
	err = tx.QueryRow("SELECT room_id, status FROM orders WHERE id = ? AND user_id = ?", orderID, userID).Scan(&roomID, &status)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "订单不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询订单失败",
		})
		return
	}

	if status == "cancelled" || status == "checked_out" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "该订单无法取消",
		})
		return
	}

	_, err = tx.Exec("UPDATE orders SET status = 'cancelled' WHERE id = ?", orderID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "取消订单失败",
		})
		return
	}

	_, err = tx.Exec("UPDATE rooms SET available_count = available_count + 1 WHERE id = ?", roomID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "恢复房间库存失败",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "订单取消成功",
	})
}

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

func CreateInvoice(c *gin.Context) {
	var req models.CreateInvoiceRequest
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

	var orderExists int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM orders 
		WHERE id = ? AND user_id = ? AND status = 'confirmed'`,
		req.OrderID, userID).Scan(&orderExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "检查订单状态失败",
		})
		return
	}

	if orderExists == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "订单不存在或无法开票",
		})
		return
	}

	var invoiceExists int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM invoices WHERE order_id = ? AND user_id = ?`,
		req.OrderID, userID).Scan(&invoiceExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "检查发票状态失败",
		})
		return
	}

	if invoiceExists > 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "该订单已申请过发票",
		})
		return
	}

	var orderNo string
	var amount float64
	err = db.QueryRow(`
		SELECT order_no, total_amount FROM orders WHERE id = ?`,
		req.OrderID).Scan(&orderNo, &amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取订单信息失败",
		})
		return
	}

	invoiceNo := fmt.Sprintf("INV%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)

	result, err := db.Exec(`
		INSERT INTO invoices (
			user_id, order_id, order_no, invoice_type, invoice_title, 
			tax_number, bank_name, bank_account, address, phone, 
			email, amount, status, invoice_no
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.OrderID, orderNo, req.InvoiceType, req.InvoiceTitle,
		req.TaxNumber, req.BankName, req.BankAccount, req.Address, req.Phone,
		req.Email, amount, "pending", invoiceNo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "申请发票失败",
		})
		return
	}

	invoiceID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "发票申请成功",
		Data: map[string]interface{}{
			"invoice_id": invoiceID,
			"invoice_no": invoiceNo,
		},
	})
}

func GetInvoiceList(c *gin.Context) {
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
		SELECT id, user_id, order_id, order_no, invoice_type, invoice_title, 
		       tax_number, bank_name, bank_account, address, phone, 
		       email, amount, status, invoice_no, created_at, updated_at
		FROM invoices WHERE user_id = ?`

	countQuery := "SELECT COUNT(*) FROM invoices WHERE user_id = ?"

	if status != "" {
		baseQuery += " AND status = ?"
		countQuery += " AND status = ?"
		baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
		rows, _ = db.Query(baseQuery, userID, status, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID, status)
	} else {
		baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
		rows, _ = db.Query(baseQuery, userID, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID)
	}

	if rows == nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取发票列表失败",
		})
		return
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.UserID, &invoice.OrderID, &invoice.OrderNo,
			&invoice.InvoiceType, &invoice.InvoiceTitle, &invoice.TaxNumber,
			&invoice.BankName, &invoice.BankAccount, &invoice.Address, &invoice.Phone,
			&invoice.Email, &invoice.Amount, &invoice.Status, &invoice.InvoiceNo,
			&invoice.CreatedAt, &invoice.UpdatedAt)
		if err != nil {
			continue
		}
		invoices = append(invoices, invoice)
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
			"invoices":  invoices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetInvoiceDetail(c *gin.Context) {
	invoiceID := c.Param("id")
	if invoiceID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少发票ID",
		})
		return
	}

	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 1
	}

	db := database.GetDB()
	var invoice models.Invoice

	err := db.QueryRow(`
		SELECT id, user_id, order_id, order_no, invoice_type, invoice_title, 
		       tax_number, bank_name, bank_account, address, phone, 
		       email, amount, status, invoice_no, created_at, updated_at
		FROM invoices WHERE id = ? AND user_id = ?`, invoiceID, userID).Scan(
		&invoice.ID, &invoice.UserID, &invoice.OrderID, &invoice.OrderNo,
		&invoice.InvoiceType, &invoice.InvoiceTitle, &invoice.TaxNumber,
		&invoice.BankName, &invoice.BankAccount, &invoice.Address, &invoice.Phone,
		&invoice.Email, &invoice.Amount, &invoice.Status, &invoice.InvoiceNo,
		&invoice.CreatedAt, &invoice.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "发票不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取发票详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    invoice,
	})
}

func GetInvoiceableOrders(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	rows, err := db.Query(`
		SELECT o.id, o.order_no, o.hotel_id, o.room_id, o.check_in, o.check_out,
		       o.guest_name, o.total_amount, o.status, o.created_at,
		       h.name as hotel_name, r.name as room_name
		FROM orders o
		LEFT JOIN hotels h ON o.hotel_id = h.id
		LEFT JOIN rooms r ON o.room_id = r.id
		WHERE o.user_id = ? AND o.status = 'confirmed' 
		AND o.id NOT IN (SELECT order_id FROM invoices WHERE user_id = ?)
		ORDER BY o.created_at DESC`, userID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取可开票订单失败",
		})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.OrderNo, &order.HotelID, &order.RoomID,
			&order.CheckIn, &order.CheckOut, &order.GuestName,
			&order.TotalAmount, &order.Status, &order.CreatedAt,
			&order.HotelName, &order.RoomName)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    orders,
	})
}

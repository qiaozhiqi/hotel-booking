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

	if req.InvoiceType != "personal" && req.InvoiceType != "company" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "发票类型错误，只能是 personal 或 company",
		})
		return
	}

	if req.InvoiceType == "company" && req.TaxNumber == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "企业发票需要提供税号",
		})
		return
	}

	db := database.GetDB()

	var order models.Order
	var totalAmount float64
	err = db.QueryRow(`
		SELECT id, order_no, total_amount, status, hotel_name, room_name, check_in, check_out
		FROM orders WHERE id = ? AND user_id = ?`, req.OrderID, userID).Scan(
		&order.ID, &order.OrderNo, &totalAmount, &order.Status,
		&order.HotelName, &order.RoomName, &order.CheckIn, &order.CheckOut)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "订单不存在或不属于当前用户",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "查询订单信息失败",
		})
		return
	}

	var existingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM invoices WHERE order_id = ?", req.OrderID).Scan(&existingCount)
	if err == nil && existingCount > 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "该订单已申请过发票",
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
		userID, order.ID, order.OrderNo, req.InvoiceType, req.InvoiceTitle,
		req.TaxNumber, req.BankName, req.BankAccount, req.Address, req.Phone,
		req.Email, totalAmount, "pending", invoiceNo)

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
		       email, amount, status, invoice_no, invoice_url, created_at, updated_at
		FROM invoices WHERE user_id = ?`

	countQuery := "SELECT COUNT(*) FROM invoices WHERE user_id = ?"

	if status != "" {
		baseQuery += " AND status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
		countQuery += " AND status = ?"
		rows, _ = db.Query(baseQuery, userID, status, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID, status)
	} else {
		baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
		rows, _ = db.Query(baseQuery, userID, pageSize, offset)
		countRows, _ = db.Query(countQuery, userID)
	}

	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		var taxNumber, bankName, bankAccount, address, phone, invoiceURL sql.NullString
		err := rows.Scan(
			&invoice.ID, &invoice.UserID, &invoice.OrderID, &invoice.OrderNo,
			&invoice.InvoiceType, &invoice.InvoiceTitle, &taxNumber, &bankName,
			&bankAccount, &address, &phone, &invoice.Email, &invoice.Amount,
			&invoice.Status, &invoice.InvoiceNo, &invoiceURL, &invoice.CreatedAt, &invoice.UpdatedAt)
		if err != nil {
			continue
		}
		if taxNumber.Valid {
			invoice.TaxNumber = taxNumber.String
		}
		if bankName.Valid {
			invoice.BankName = bankName.String
		}
		if bankAccount.Valid {
			invoice.BankAccount = bankAccount.String
		}
		if address.Valid {
			invoice.Address = address.String
		}
		if phone.Valid {
			invoice.Phone = phone.String
		}
		if invoiceURL.Valid {
			invoice.InvoiceURL = invoiceURL.String
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

	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	var invoice models.Invoice
	var taxNumber, bankName, bankAccount, address, phone, invoiceURL sql.NullString
	err = db.QueryRow(`
		SELECT id, user_id, order_id, order_no, invoice_type, invoice_title, 
		       tax_number, bank_name, bank_account, address, phone, 
		       email, amount, status, invoice_no, invoice_url, created_at, updated_at
		FROM invoices WHERE id = ? AND user_id = ?`, invoiceID, userID).Scan(
		&invoice.ID, &invoice.UserID, &invoice.OrderID, &invoice.OrderNo,
		&invoice.InvoiceType, &invoice.InvoiceTitle, &taxNumber, &bankName,
		&bankAccount, &address, &phone, &invoice.Email, &invoice.Amount,
		&invoice.Status, &invoice.InvoiceNo, &invoiceURL, &invoice.CreatedAt, &invoice.UpdatedAt)

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

	if taxNumber.Valid {
		invoice.TaxNumber = taxNumber.String
	}
	if bankName.Valid {
		invoice.BankName = bankName.String
	}
	if bankAccount.Valid {
		invoice.BankAccount = bankAccount.String
	}
	if address.Valid {
		invoice.Address = address.String
	}
	if phone.Valid {
		invoice.Phone = phone.String
	}
	if invoiceURL.Valid {
		invoice.InvoiceURL = invoiceURL.String
	}

	var order models.Order
	db.QueryRow(`
		SELECT id, order_no, hotel_name, room_name, check_in, check_out, total_amount, status
		FROM orders WHERE id = ?`, invoice.OrderID).Scan(
		&order.ID, &order.OrderNo, &order.HotelName, &order.RoomName,
		&order.CheckIn, &order.CheckOut, &order.TotalAmount, &order.Status)

	result := map[string]interface{}{
		"invoice": invoice,
		"order":   order,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

func GetInvoiceableOrders(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	rows, err := db.Query(`
		SELECT o.id, o.order_no, o.hotel_name, o.room_name, o.check_in, o.check_out, 
		       o.total_amount, o.status, o.created_at
		FROM orders o
		LEFT JOIN invoices i ON o.id = i.order_id
		WHERE o.user_id = ? AND i.id IS NULL AND o.status != 'cancelled'
		ORDER BY o.created_at DESC`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取可开票订单列表失败",
		})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.OrderNo, &order.HotelName, &order.RoomName,
			&order.CheckIn, &order.CheckOut, &order.TotalAmount, &order.Status,
			&order.CreatedAt)
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

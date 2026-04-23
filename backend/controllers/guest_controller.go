package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetGuestList(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetHeader("X-User-ID"))
	if err != nil || userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	rows, err := db.Query(`
		SELECT id, user_id, name, phone, id_type, id_number, is_default, created_at, updated_at
		FROM guests WHERE user_id = ? ORDER BY is_default DESC, created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取常用入住人列表失败",
		})
		return
	}
	defer rows.Close()

	var guests []models.Guest
	for rows.Next() {
		var guest models.Guest
		var idType, idNumber sql.NullString
		err := rows.Scan(&guest.ID, &guest.UserID, &guest.Name, &guest.Phone,
			&idType, &idNumber, &guest.IsDefault, &guest.CreatedAt, &guest.UpdatedAt)
		if err != nil {
			continue
		}
		if idType.Valid {
			guest.IDType = idType.String
		}
		if idNumber.Valid {
			guest.IDNumber = idNumber.String
		}
		guests = append(guests, guest)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    guests,
	})
}

func GetGuestDetail(c *gin.Context) {
	guestID := c.Param("id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少入住人ID",
		})
		return
	}

	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 1
	}

	db := database.GetDB()
	var guest models.Guest
	var idType, idNumber sql.NullString

	err := db.QueryRow(`
		SELECT id, user_id, name, phone, id_type, id_number, is_default, created_at, updated_at
		FROM guests WHERE id = ? AND user_id = ?`, guestID, userID).Scan(
		&guest.ID, &guest.UserID, &guest.Name, &guest.Phone,
		&idType, &idNumber, &guest.IsDefault, &guest.CreatedAt, &guest.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "入住人不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取入住人详情失败",
		})
		return
	}

	if idType.Valid {
		guest.IDType = idType.String
	}
	if idNumber.Valid {
		guest.IDNumber = idNumber.String
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    guest,
	})
}

func CreateGuest(c *gin.Context) {
	var req models.CreateGuestRequest
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

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "创建入住人失败",
		})
		return
	}

	if req.IsDefault {
		_, err = tx.Exec("UPDATE guests SET is_default = 0 WHERE user_id = ?", userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "创建入住人失败",
			})
			return
		}
	}

	if !req.IsDefault {
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM guests WHERE user_id = ?", userID).Scan(&count)
		if err != nil {
			req.IsDefault = true
		} else if count == 0 {
			req.IsDefault = true
		}
	}

	result, err := tx.Exec(`
		INSERT INTO guests (user_id, name, phone, id_type, id_number, is_default)
		VALUES (?, ?, ?, ?, ?, ?)`,
		userID, req.Name, req.Phone, req.IDType, req.IDNumber, req.IsDefault)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "创建入住人失败",
		})
		return
	}

	tx.Commit()

	guestID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "创建成功",
		Data: map[string]interface{}{
			"guest_id": guestID,
		},
	})
}

func UpdateGuest(c *gin.Context) {
	guestID := c.Param("id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少入住人ID",
		})
		return
	}

	var req models.UpdateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误",
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
			Message: "更新入住人失败",
		})
		return
	}

	var existingGuest models.Guest
	err = tx.QueryRow("SELECT id, user_id, is_default FROM guests WHERE id = ? AND user_id = ?", guestID, userID).Scan(
		&existingGuest.ID, &existingGuest.UserID, &existingGuest.IsDefault)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "入住人不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "更新入住人失败",
		})
		return
	}

	if req.IsDefault != nil && *req.IsDefault {
		_, err = tx.Exec("UPDATE guests SET is_default = 0 WHERE user_id = ?", userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "更新入住人失败",
			})
			return
		}
	}

	query := "UPDATE guests SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Name != "" {
		query += ", name = ?"
		args = append(args, req.Name)
	}
	if req.Phone != "" {
		query += ", phone = ?"
		args = append(args, req.Phone)
	}
	if req.IDType != "" {
		query += ", id_type = ?"
		args = append(args, req.IDType)
	}
	if req.IDNumber != "" {
		query += ", id_number = ?"
		args = append(args, req.IDNumber)
	}
	if req.IsDefault != nil {
		query += ", is_default = ?"
		args = append(args, *req.IsDefault)
	}

	query += " WHERE id = ? AND user_id = ?"
	args = append(args, guestID, userID)

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "更新入住人失败",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "更新成功",
	})
}

func DeleteGuest(c *gin.Context) {
	guestID := c.Param("id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少入住人ID",
		})
		return
	}

	userID, _ := strconv.Atoi(c.GetHeader("X-User-ID"))
	if userID == 0 {
		userID = 1
	}

	db := database.GetDB()

	result, err := db.Exec("DELETE FROM guests WHERE id = ? AND user_id = ?", guestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "删除入住人失败",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "入住人不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "删除成功",
	})
}

func SetDefaultGuest(c *gin.Context) {
	guestID := c.Param("id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少入住人ID",
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
			Message: "设置默认入住人失败",
		})
		return
	}

	_, err = tx.Exec("UPDATE guests SET is_default = 0 WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "设置默认入住人失败",
		})
		return
	}

	result, err := tx.Exec("UPDATE guests SET is_default = 1 WHERE id = ? AND user_id = ?", guestID, userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "设置默认入住人失败",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "入住人不存在",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "设置成功",
	})
}

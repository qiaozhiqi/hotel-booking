package controllers

import (
	"database/sql"
	"hotel-booking/database"
	"hotel-booking/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	db := database.GetDB()
	var user models.User
	err := db.QueryRow("SELECT id, username, password, email, phone FROM users WHERE username = ?", req.Username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "服务器内部错误",
		})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "登录成功",
		Data:    user,
	})
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	db := database.GetDB()

	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "服务器内部错误",
		})
		return
	}

	if exists > 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "用户名已存在",
		})
		return
	}

	result, err := db.Exec("INSERT INTO users (username, password, email, phone) VALUES (?, ?, ?, ?)",
		req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "注册失败",
		})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "注册成功",
		Data:    map[string]int64{"user_id": id},
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "缺少用户ID",
		})
		return
	}

	db := database.GetDB()
	var user models.User
	err := db.QueryRow("SELECT id, username, email, phone, created_at FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    404,
				Message: "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    user,
	})
}

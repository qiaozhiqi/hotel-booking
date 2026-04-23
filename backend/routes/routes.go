package routes

import (
	"hotel-booking/controllers"
	"hotel-booking/security"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-ID", "X-API-Key", "X-Signature", "X-Timestamp", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)
		api.GET("/user", controllers.GetUserInfo)

		api.GET("/hotels", controllers.GetHotelList)
		api.GET("/hotels/:id", controllers.GetHotelDetail)
		api.GET("/cities", controllers.GetCities)
		api.GET("/filter-options", controllers.GetFilterOptions)

		api.POST("/orders", controllers.CreateOrder)
		api.GET("/orders", controllers.GetOrderList)
		api.GET("/orders/:id", controllers.GetOrderDetail)
		api.POST("/orders/:id/cancel", controllers.CancelOrder)

		api.GET("/guests", controllers.GetGuestList)
		api.GET("/guests/:id", controllers.GetGuestDetail)
		api.POST("/guests", controllers.CreateGuest)
		api.PUT("/guests/:id", controllers.UpdateGuest)
		api.DELETE("/guests/:id", controllers.DeleteGuest)
		api.POST("/guests/:id/default", controllers.SetDefaultGuest)

		api.GET("/suppliers", controllers.GetSupplierList)
		api.POST("/suppliers/:code/pull", controllers.PullSupplierData)
		api.GET("/suppliers/:code/status", controllers.GetSyncStatus)

		priceInventory := api.Group("/price-inventory")
		{
			priceInventory.GET("", controllers.GetPriceInventory)
			priceInventory.GET("/range", controllers.GetPriceInventoryByDateRange)
			priceInventory.GET("/room-price", controllers.GetRoomCurrentPrice)
			priceInventory.GET("/summary/:code", controllers.GetSupplierPriceInventorySummary)
			priceInventory.DELETE("/clear/:code", controllers.ClearPriceInventory)
		}

		shiji := api.Group("/shiji")
		{
			qiuguo := shiji.Group("/qiuguo")
			{
				qiuguo.GET("/status", controllers.GetQiuguoSyncStatus)
				
				qiuguoProtected := qiuguo.Group("")
				qiuguoProtected.Use(security.SupplierPushAuthMiddleware("shiji_qiuguo"))
				{
					qiuguoProtected.POST("/push", controllers.HandleQiuguoPush)
				}
			}
		}

		mock := api.Group("/mock")
		{
			mock.GET("/huazhu/hotels", controllers.MockHuazhuGetHotels)
			mock.GET("/huazhu/hotels/:id", controllers.MockHuazhuGetHotelDetail)
		}
	}

	return r
}

package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.POST("/adminlogin", handlers.LoginHandler)

	r.Use(middleware.AdminAuthMiddleware())
	{

		r.GET("/dashboard", handlers.DashBoard)
		r.GET("/sales-report", handlers.FilteredSalesReport)

		//user management
		users := r.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:page", handlers.GetUsers)
			users.GET("/block", handlers.BlockUser)
			users.GET("/unblock", handlers.UnBlockUser)
		}

		//products management
		products := r.Group("/products")
		{
			products.GET("", handlers.ShowAllProductsFromAdmin)
			products.POST("", handlers.AddProducts)
			products.PUT("", handlers.UpdateProduct) //update the product quantity
			products.DELETE("", handlers.DeleteProducts)
			products.POST("/upload-image", handlers.UploadImage)

		}

		//category management
		category := r.Group("/category")
		{
			category.GET("", handlers.GetCategory)
			category.POST("", handlers.AddCategory)
			category.PUT("", handlers.UpdateCategory)
			category.DELETE("", handlers.DeleteCategory)

		}

		//order
		order := r.Group("/order")
		{
			order.GET("", handlers.GetAllOrderDetailsForAdmin)
			order.GET("/approve", handlers.ApproveOrder)
			order.GET("/cancel", handlers.CancelOrderFromAdmin)
		}

		//image cropping issuess
		r.POST("/image-crop", handlers.CropImage)

	}
	return r
}

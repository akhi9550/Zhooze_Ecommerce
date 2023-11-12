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

		users := r.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.PUT("/block", handlers.BlockUser)
			users.PUT("/unblock", handlers.UnBlockUser)
		}

		products := r.Group("/products")
		{
			products.GET("", handlers.ShowAllProductsFromAdmin)
			products.POST("", handlers.AddProducts)
			products.PUT("", handlers.UpdateProduct) //update the product quantity
			products.DELETE("", handlers.DeleteProducts)
			products.POST("/upload-image", handlers.UploadImage)

		}

		category := r.Group("/category")
		{
			category.GET("", handlers.GetCategory)
			category.POST("", handlers.AddCategory)
			category.PUT("", handlers.UpdateCategory)
			category.DELETE("", handlers.DeleteCategory)

		}

		order := r.Group("/order")
		{
			order.GET("", handlers.GetAllOrderDetailsForAdmin)
			order.GET("/approve", handlers.ApproveOrder)
			order.GET("/cancel", handlers.CancelOrderFromAdmin)
		}

		offer := r.Group("/offer")
		{
			coupons := offer.Group("/coupons")
			{
				coupons.POST("", handlers.AddCoupon)
				coupons.GET("", handlers.GetCoupon)
				// 	coupons.PATCH("/expire/:id", handlers.ExpireCoupon)
				// }
				// offer.POST("/product-offer", handlers.AddProdcutOffer)
				// offer.POST("/category-offer", handlers.AddCategoryOffer)
			}
			r.POST("/image-crop", handlers.CropImage)

		}
		return r
	}

	////////////////////payment methodddddddddddddddd
}

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
		r.GET("/sales-report-date", handlers.SalesReportByDate)

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
			products.PUT("", handlers.UpdateProduct) 
			products.DELETE("", handlers.DeleteProducts)
			products.GET("/search", handlers.SearchProducts)
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

		payment := r.Group("/payment-method")
		{
			payment.POST("", handlers.AddPaymentMethod)
			payment.GET("", handlers.ListPaymentMethods)
			payment.DELETE("", handlers.DeletePaymentMethod)
		}
		coupons := r.Group("/coupons")
		{
			coupons.POST("", handlers.AddCoupon)
			coupons.GET("", handlers.GetCoupon)
			coupons.PATCH("", handlers.ExpireCoupon)
		}
		Productoffer := r.Group("/productoffer")
		{
			Productoffer.POST("", handlers.AddProdcutOffer)
			Productoffer.GET("", handlers.GetProductOffer)
		}
		Categoryoffer := r.Group("/categoryoffer")
		{
			Categoryoffer.POST("", handlers.AddCategoryOffer)
			Categoryoffer.GET("", handlers.GetCategoryOffer)
		}

	}
	return r
}

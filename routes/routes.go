package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AllRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	//user
	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	r.POST("/send-otp", handlers.SendOtp)
	r.POST("verify-otp", handlers.VerifyOtp)

	r.GET("/products", handlers.AllProducts)
	r.GET("/page/:page", handlers.ShowAllProducts) //arranging order
	r.POST("/filter", handlers.FilterCategory)

	r.GET("/address", middleware.UserAuthMiddleware, handlers.GetAllAddress)
	r.POST("/add-address", middleware.UserAuthMiddleware, handlers.AddAddress)
	r.GET("/user-details", middleware.UserAuthMiddleware, handlers.UserDetails)
	r.PATCH("/edit-user-profile", middleware.UserAuthMiddleware, handlers.UpdateUserDetails)
	// r.POST("/update-password",handlers.UpdatePassword)

	//admin
	r.POST("/adminlogin", handlers.LoginHandler)
	// r.Use(middleware.AdminAuthMiddleware)
	// {
	r.GET("/dashboard", handlers.DashBoard)

	r.GET("/getusers", handlers.GetUsers)
	r.POST("/block", handlers.BlockUser)
	r.POST("/unblock", handlers.UnBlockUser)

	r.POST("/add-category", handlers.AddCategory)
	r.PUT("/update-category", handlers.UpdateCategory)
	r.DELETE("/delete-category", handlers.DeleteCategory)

	r.GET("/products-ad", handlers.AllProducts)
	r.POST("add-product", handlers.AddProducts)
	// r.PUT("/update-product", handlers.UpdateProduct)
	r.DELETE("/delete-product", handlers.DeleteProducts)

	r.GET("/order/:page", handlers.GetAllOrderDetailsForAdmin)
	r.GET("/approve-order/:order_id", handlers.ApproveOrder)
	r.GET("/cancel-order/:order_id", handlers.CancelOrderFromAdmin)
	r.PUT("/refund-order/:order_id", handlers.RefundUser)
	// }

	return r
}

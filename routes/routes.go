package routes

import (
	"Zhooze/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AllRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	r.POST("/send-otp", handlers.SendOtp)
	r.POST("verify-otp", handlers.VerifyOtp)

	r.GET("/page/:page", handlers.ShowAllProducts)
	r.POST("/filter", handlers.FilterCategory)

	r.POST("/adminlogin", handlers.LoginHandler)
	r.GET("/dashboard", handlers.DashBoard)

	r.GET("/getusers", handlers.GetUsers)
	r.POST("/block", handlers.BlockUser)
	r.POST("/unblock", handlers.UnBlockUser)

	r.GET("/products", handlers.AllProducts)
	r.POST("add-product", handlers.AddProducts)
	// r.PUT("/update-product",handlers.UpdateProduct)
	r.DELETE("delete-product", handlers.DeleteProducts)

	r.POST("/add-category", handlers.AddCategory)
	r.PUT("/update-category", handlers.UpdateCategory)
	r.DELETE("/delete-category", handlers.DeleteCategory)
	return r
}

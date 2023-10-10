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

	return r
}

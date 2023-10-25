package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AllRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	//admin
	r.POST("/adminlogin", handlers.LoginHandler)
	r.GET("/dashboard", middleware.AdminAuthMiddleware(), handlers.DashBoard)
	//usermanagefrom
	r.GET("/getusers", middleware.AdminAuthMiddleware(), handlers.GetUsers)
	r.POST("/block", middleware.AdminAuthMiddleware(), handlers.BlockUser)
	r.POST("/unblock", middleware.AdminAuthMiddleware(), handlers.UnBlockUser)
	//Category
	r.POST("/add-category", middleware.AdminAuthMiddleware(), handlers.AddCategory)
	r.PUT("/update-category", middleware.AdminAuthMiddleware(), handlers.UpdateCategory)
	r.DELETE("/delete-category", middleware.AdminAuthMiddleware(), handlers.DeleteCategory)
	//Product
	r.GET("/products-ad", middleware.AdminAuthMiddleware(), handlers.AllProducts)
	r.POST("add-product", middleware.AdminAuthMiddleware(), handlers.AddProducts)
	r.PATCH("/update-product", middleware.AdminAuthMiddleware(), handlers.UpdateProduct)
	r.DELETE("/delete-product", middleware.AdminAuthMiddleware(), handlers.DeleteProducts)
	//Order
	r.GET("/order/:page", middleware.AdminAuthMiddleware(), handlers.GetAllOrderDetailsForAdmin)
	r.GET("/approve-order/:order_id", middleware.AdminAuthMiddleware(), handlers.ApproveOrder)
	r.GET("/cancel-order/:order_id", middleware.AdminAuthMiddleware(), handlers.CancelOrderFromAdmin)
	r.PUT("/refund-order/:order_id", middleware.AdminAuthMiddleware(), handlers.RefundUser)

	//user

	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	r.POST("/send-otp", handlers.SendOtp)
	r.POST("verify-otp", handlers.VerifyOtp)

	r.PUT("/changepassword", handlers.ChangePassword)
	r.GET("/forgot-password", handlers.ForgotPasswordSend)
	// r.POST("/forgot-password", handlers.ForgotPasswordVerifyAndChange)

	r.GET("/products", handlers.AllProducts)
	r.GET("/page/:page", handlers.ShowAllProducts) //arranging page and count
	r.POST("/filter", handlers.FilterCategory)

	//address
	r.GET("/address", middleware.UserAuthMiddleware(), handlers.GetAllAddress)
	r.POST("/add-address", middleware.UserAuthMiddleware(), handlers.AddAddress)
	r.GET("/user-details", middleware.UserAuthMiddleware(), handlers.UserDetails)
	r.PATCH("/edit-user-profile", middleware.UserAuthMiddleware(), handlers.UpdateUserDetails)
	r.PATCH("/edit-address/:address_id", middleware.UserAuthMiddleware(), handlers.UpdateAddress)

	//ORDERS
	r.GET("/orders/:page", middleware.UserAuthMiddleware(), handlers.GetOrderDetails)
	r.PUT("/cancel-orders/:id", middleware.UserAuthMiddleware(), handlers.CancelOrder)
	r.GET("/checkout", middleware.UserAuthMiddleware(), handlers.CheckOut)
	r.GET("/place-order/:order_id/:payment", middleware.UserAuthMiddleware(), handlers.PlaceOrder)

	//CART
	r.POST("/addtocart/:id", middleware.UserAuthMiddleware(), handlers.AddToCart)
	r.DELETE("/removefromcart/:id", middleware.UserAuthMiddleware(), handlers.RemoveFromCart)
	r.POST("/updatequantityadd", middleware.UserAuthMiddleware(), handlers.UpdateQuantityAdd)
	r.POST("/updatequantityless", middleware.UserAuthMiddleware(), handlers.UpdateQuantityless)
	r.GET("/displaycart", middleware.UserAuthMiddleware(), handlers.DisplayCart)
	r.DELETE("/emptycart", middleware.UserAuthMiddleware(), handlers.EmptyCart)

	return r
}

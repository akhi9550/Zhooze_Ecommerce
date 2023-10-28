package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func AllRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	//ADMIN
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	r.GET("/products-ad", middleware.AdminAuthMiddleware(), handlers.ShowAllProductsFromAdmin)
	r.POST("add-product", middleware.AdminAuthMiddleware(), handlers.AddProducts)
	r.PATCH("/update-product", middleware.AdminAuthMiddleware(), handlers.UpdateProduct)
	r.DELETE("/delete-product", middleware.AdminAuthMiddleware(), handlers.DeleteProducts)

	//Order
	r.GET("/order", middleware.AdminAuthMiddleware(), handlers.GetAllOrderDetailsForAdmin)
	r.GET("/approve-order", middleware.AdminAuthMiddleware(), handlers.ApproveOrder)
	r.GET("/cancel-order", middleware.AdminAuthMiddleware(), handlers.CancelOrderFromAdmin)
	r.PUT("/refund-order", middleware.AdminAuthMiddleware(), handlers.RefundUser)
	//IMAGE CROPPING
	r.POST("/image-crop",middleware.AdminAuthMiddleware(),handlers.CropImage)

	//USER//

	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	r.POST("/send-otp", handlers.SendOtp)
	r.POST("verify-otp", handlers.VerifyOtp)

	//security
	r.GET("/forgot-password", handlers.ForgotPasswordSend)
	r.POST("/forgot-password", handlers.ForgotPasswordVerifyAndChange)
	r.PUT("/changepassword", middleware.UserAuthMiddleware(), handlers.ChangePassword)

	//products
	r.GET("/products", handlers.AllProducts)
	r.GET("/page", handlers.ShowAllProducts) //arranging page and count
	r.POST("/filter", handlers.FilterCategory)

	//profile
	r.GET("/address", middleware.UserAuthMiddleware(), handlers.GetAllAddress)
	r.POST("/add-address", middleware.UserAuthMiddleware(), handlers.AddAddress)
	r.GET("/user-details", middleware.UserAuthMiddleware(), handlers.UserDetails)
	r.PATCH("/edit-user-profile", middleware.UserAuthMiddleware(), handlers.UpdateUserDetails)
	r.PATCH("/edit-address", middleware.UserAuthMiddleware(), handlers.UpdateAddress)

	//ORDERS
	r.GET("/orders", middleware.UserAuthMiddleware(), handlers.GetOrderDetails)
	r.PUT("/cancel-orders", middleware.UserAuthMiddleware(), handlers.CancelOrder)
	r.GET("/checkout", middleware.UserAuthMiddleware(), handlers.CheckOut)
	r.GET("/place-order", middleware.UserAuthMiddleware(), handlers.PlaceOrder)

	//CART
	r.POST("/addtocart", middleware.UserAuthMiddleware(), handlers.AddToCart)
	r.DELETE("/removefromcart", middleware.UserAuthMiddleware(), handlers.RemoveFromCart)
	r.GET("/displaycart", middleware.UserAuthMiddleware(), handlers.DisplayCart)
	r.DELETE("/emptycart", middleware.UserAuthMiddleware(), handlers.EmptyCart)
	///////
	r.PUT("/updatequantityadd", middleware.UserAuthMiddleware(), handlers.UpdateQuantityAdd)
	r.PUT("/updatequantityless", middleware.UserAuthMiddleware(), handlers.UpdateQuantityless)

	//PAYMENT
	r.GET("/razorpay",handlers.MakePaymentRazorPay)
	r.GET("/update_status",handlers.VerifyPayment)
	return r
}

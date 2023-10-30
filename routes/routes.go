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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//***********************************ADMIN***********************************//

	r.POST("/adminlogin", handlers.LoginHandler)
	r.GET("/dashboard", middleware.AdminAuthMiddleware(), handlers.DashBoard)
	r.GET("/sales-report", middleware.AdminAuthMiddleware(), handlers.FilteredSalesReport)

	//USERMANAGEMENT
	r.GET("/getusers", middleware.AdminAuthMiddleware(), handlers.GetUsers)
	r.POST("/block", middleware.AdminAuthMiddleware(), handlers.BlockUser)
	r.POST("/unblock", middleware.AdminAuthMiddleware(), handlers.UnBlockUser)

	//CATEGORY
	r.POST("/add-category", middleware.AdminAuthMiddleware(), handlers.AddCategory)
	r.PUT("/update-category", middleware.AdminAuthMiddleware(), handlers.UpdateCategory)
	r.DELETE("/delete-category", middleware.AdminAuthMiddleware(), handlers.DeleteCategory)

	//PRODUCT
	r.GET("/products-ad", middleware.AdminAuthMiddleware(), handlers.ShowAllProductsFromAdmin)
	r.POST("add-product", middleware.AdminAuthMiddleware(), handlers.AddProducts)
	r.PATCH("/update-product", middleware.AdminAuthMiddleware(), handlers.UpdateProduct)
	r.DELETE("/delete-product", middleware.AdminAuthMiddleware(), handlers.DeleteProducts)

	//ORDER
	r.GET("/order", middleware.AdminAuthMiddleware(), handlers.GetAllOrderDetailsForAdmin)
	r.GET("/approve-order", middleware.AdminAuthMiddleware(), handlers.ApproveOrder)
	r.GET("/cancel-order", middleware.AdminAuthMiddleware(), handlers.CancelOrderFromAdmin)
	r.PUT("/refund-order", middleware.AdminAuthMiddleware(), handlers.RefundUser)

	//IMAGE CROPPING
	r.POST("/image-crop", middleware.AdminAuthMiddleware(), handlers.CropImage)

	//***********************************USER***********************************//

	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	//OTP
	r.POST("/send-otp", handlers.SendOtp)
	r.POST("verify-otp", handlers.VerifyOtp)

	//SECURITY
	r.GET("/forgot-password", handlers.ForgotPasswordSend)
	r.POST("/forgot-password", handlers.ForgotPasswordVerifyAndChange)

	//PRODUCT
	r.GET("/products", handlers.AllProducts)
	r.GET("/page", handlers.ShowAllProducts) //arranging page and count
	r.POST("/filter", handlers.FilterCategory)

	//PROFILE
	r.GET("/address", middleware.UserAuthMiddleware(), handlers.GetAllAddress)
	r.POST("/add-address", middleware.UserAuthMiddleware(), handlers.AddAddress)
	r.GET("/user-details", middleware.UserAuthMiddleware(), handlers.UserDetails)
	r.PATCH("/edit-user-profile", middleware.UserAuthMiddleware(), handlers.UpdateUserDetails)
	r.PATCH("/edit-address", middleware.UserAuthMiddleware(), handlers.UpdateAddress)
	r.PUT("/changepassword", middleware.UserAuthMiddleware(), handlers.ChangePassword)

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
	r.PUT("/updatequantityadd", middleware.UserAuthMiddleware(), handlers.UpdateQuantityAdd)
	r.PUT("/updatequantityless", middleware.UserAuthMiddleware(), handlers.UpdateQuantityless)

	//PAYMENT
	r.GET("/razorpay", handlers.MakePaymentRazorPay)
	r.GET("/update_status", handlers.VerifyPayment)

	//WISHLIST
	r.GET("/wishlist", middleware.UserAuthMiddleware(), handlers.GetWishList)
	r.POST("/wishlist-add", middleware.UserAuthMiddleware(), handlers.AddToWishlist)
	r.DELETE("/wishlist-remove", middleware.UserAuthMiddleware(), handlers.RemoveFromWishlist)

	return r
}

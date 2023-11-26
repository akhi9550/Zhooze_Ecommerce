package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.POST("/signup", handlers.UserSignup)
	r.POST("/userlogin", handlers.Userlogin)

	r.POST("/send-otp", handlers.SendOtp)
	r.POST("/verify-otp", handlers.VerifyOtp)

	r.POST("/forgot-password", handlers.ForgotPasswordSend)
	r.POST("/forgot-password-verify", handlers.ForgotPasswordVerifyAndChange)

	r.GET("/razorpay", handlers.MakePaymentRazorPay)
	r.GET("/update_status", handlers.VerifyPayment)

	products := r.Group("/products")
	{
		products.GET("", handlers.ShowAllProducts)
		products.POST("/filter", handlers.FilterCategory)
	

	}

	r.Use(middleware.UserAuthMiddleware())
	{

		address := r.Group("/address")
		{
			address.GET("", handlers.GetAllAddress)
			address.POST("", handlers.AddAddress)
			address.PUT("", handlers.UpdateAddress)
			address.DELETE("", handlers.DeleteAddressByID)
		}
		users := r.Group("/users")
		{
			users.GET("", handlers.UserDetails)
			users.PUT("", handlers.UpdateUserDetails)
			users.PUT("/changepassword", handlers.ChangePassword)
		}

		wishlist := r.Group("/wishlist")
		{
			wishlist.POST("", handlers.AddToWishlist)
			wishlist.GET("", handlers.GetWishList)
			wishlist.DELETE("", handlers.RemoveFromWishlist)
		}

		cart := r.Group("/cart")
		{
			cart.POST("", handlers.AddToCart)
			cart.DELETE("", handlers.RemoveFromCart)
			cart.GET("", handlers.DisplayCart)
			cart.DELETE("/empty", handlers.EmptyCart)
			cart.PUT("/updatequantityadd", handlers.UpdateQuantityAdd)
			cart.PUT("/updatequantityless", handlers.UpdateQuantityless)
		}

		order := r.Group("/order")
		{
			order.GET("", handlers.GetOrderDetails)
			order.PUT("", handlers.CancelOrder)
			order.GET("/place-order", handlers.PlaceOrderCOD)

		}
		checkout := r.Group("/checkout")
		{
			checkout.GET("", handlers.CheckOut)
			checkout.POST("", handlers.OrderItemsFromCart)

		}
		coupon := r.Group("/coupon")
		{
			coupon.POST("", handlers.ApplyCoupon)
		}
		wallet := r.Group("/wallet")
		{
			wallet.GET("", handlers.GetWallet)
			wallet.GET("/history", handlers.WalletHistory)
		}
	}

	return r

}

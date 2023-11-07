package routes

import (
	"Zhooze/handlers"
	"Zhooze/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	
	r.POST("/signup", handlers.UserSignup)
	r.POST("/loginlogin", handlers.Userlogin)

	//OTP
	r.POST("/send-otp", handlers.SendOtp)
	r.POST("/verify-otp", handlers.VerifyOtp)

	//SECURITY
	r.POST("/forgot-password", handlers.ForgotPasswordSend)
	r.POST("/forgot-password-verify", handlers.ForgotPasswordVerifyAndChange)

	products := r.Group("/products")
	{
		products.GET("", handlers.ShowAllProducts)
		products.GET("/page/:page", handlers.ShowAllProducts) //TO ARRANGE PAGE WITH COUNT
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

		//wishlist
		wishlist := r.Group("/wishlist")
		{

			wishlist.POST("", handlers.AddToWishlist)
			wishlist.GET("", handlers.GetWishList)
			wishlist.DELETE("", handlers.RemoveFromWishlist)
		}

		//cart
		cart := r.Group("/cart")
		{
			cart.POST("/:id", handlers.AddToCart)
			cart.DELETE("/:id", handlers.RemoveFromCart)
			cart.GET("", handlers.DisplayCart)
			cart.DELETE("", handlers.EmptyCart)
			r.PUT("updatequantityadd", handlers.UpdateQuantityAdd)
			r.PUT("updatequantityless", handlers.UpdateQuantityAdd)

		}

		//order
		order := r.Group("/order")
		{

			order.POST("", handlers.OrderItemsFromCart)
			order.GET("", handlers.GetOrderDetails)
			order.GET("/:page", handlers.GetOrderDetails)
			order.PUT("/:id", handlers.CancelOrder)
		}
		r.GET("/checkout", handlers.CheckOut)
		r.GET("/place-order", handlers.PlaceOrder)
	}

	//payment
	r.GET("/payment", handlers.MakePaymentRazorPay)
	r.GET("/payment-success", handlers.VerifyPayment)
	return r

}

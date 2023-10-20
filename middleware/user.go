package middleware

import (
	"Zhooze/config"
	"Zhooze/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// // func AuthMiddleware() gin.HandlerFunc {
// // 	return func(c *gin.Context) {
// // 		// Retrieve the JWT token from the Authorization header
// // 		authHeader := c.GetHeader("Authorization")
// // tokenString := helper.GetTokenFromHeader(authHeader)
// // 		fmt.Println(tokenString)
// // 		// Validate the token and extract the user ID
// // 		if tokenString == "" {
// // 			var err error
// // 			tokenString, err = c.Cookie("Authorization")
// // 			if err != nil {
// // 				c.AbortWithStatus(http.StatusUnauthorized)
// // 				return
// // 			}
// // 		}
// // 		userID, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
// // 		if err != nil {
// // 			fmt.Println("error is 👺👺👺👺👺👺", err)
// // 			c.AbortWithStatus(http.StatusUnauthorized)
// // 			return
// // 		}
// // 		// Add the user ID to the Gin context

// // 		c.Set("user_id", userID)
// // 		c.Set("user_email", userEmail)
// // 		// c.Set("is_admin",)

// //			// Call the next handler
// //			c.Next()
// //		}
// //	}
// func UserAuthMiddleware(c *gin.Context) {
// 	cfg, _ := config.LoadConfig()
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
// 		c.Abort()
// 		return
// 	}
// 	tokenString = strings.TrimPrefix(tokenString, "Bearer")
// 	fmt.Println(tokenString, "😒😒😒😒😒😒😒😒")
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(cfg.KEY), nil
// 	})
// 	fmt.Println(token, "👺👺👺👺👺👺👺")
// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token😎"})
// 		c.Abort()
// 		return
// 	}

// 	// Extract claims
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
// 		c.Abort()
// 		return
// 	}

// 	role, ok := claims["role"].(string)
// 	if !ok || role != "client" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
// 		c.Abort()
// 		return
// 	}

// 	userID, err := helper.ExtractUserIDFromToken(tokenString)
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "some internal problem"})
// 		c.Abort()
// 		return
// 	}

//		// Set both "role" and "user_id" in the Gin context.
//		c.Set("role", role)
//		c.Set("id", userID)
//		c.Next()
//	}
func UserAuthMiddleware(c *gin.Context) {
	cfg, _ := config.LoadConfig()
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		c.Abort()
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte(cfg.KEY), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "client" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}
	userID, err := helper.ExtractUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "some internal problem"})
		c.Abort()
		return
	}
	c.Set("role", role)
	c.Set("user_id", userID)

	c.Next()
}

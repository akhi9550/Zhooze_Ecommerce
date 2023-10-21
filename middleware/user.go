package middleware

import (
	"Zhooze/helper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		tokenString := helper.GetTokenFromHeader(authHeader)

		// Validate the token and extract the user ID
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie("Authorization")
			if err != nil {

				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		userID, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			fmt.Println("error is ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Add the user ID to the Gin context
		c.Set("user_id", userID)
		c.Set("user_email", userEmail)

		// Call the next handler
		c.Next()
	}
}

// func UserAuthMiddleware(c *gin.Context) {
// 	cfg, _ := config.LoadConfig()
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
// 		c.Abort()
// 		return
// 	}
// 	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

// 		return []byte(cfg.KEY), nil
// 	})

// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
// 		c.Abort()
// 		return
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok || !token.Valid {
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
// 	c.Set("role", role)
// 	c.Set("user_id", userID)

// 	c.Next()
// }

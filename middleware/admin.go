package middleware

import (
	"Zhooze/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)
func AdminAuthMiddleware(c *gin.Context) {
	cfg,_:=config.LoadConfig()
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		c.Abort()
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.KEY_ADMIN), nil
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
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}

	c.Set("role", role)

	c.Next()
}

// func AuthorizationMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenHeader := c.GetHeader("authorization")
// 		fmt.Println(tokenHeader, "this is the token header")
// 		if tokenHeader == "" {
// 			response := response.ClientResponse(http.StatusUnauthorized, "No auth header provided", nil, nil)
// 			c.JSON(http.StatusUnauthorized, response)
// 			c.Abort()
// 			return
// 		}

// 		splitted := strings.Split(tokenHeader, " ")
// 		if len(splitted) != 2 {
// 			response := response.ClientResponse(http.StatusUnauthorized, "Invalid Token Format", nil, nil)
// 			c.JSON(http.StatusUnauthorized, response)
// 			c.Abort()
// 			return

// 		}
// 		tokenpart := splitted[1]
// 		tokenClaims, err := helper.ValidateToken(tokenpart)
// 		if err != nil {
// 			response := response.ClientResponse(http.StatusUnauthorized, "Invalid Token  ", nil, err.Error())
// 			c.JSON(http.StatusUnauthorized, response)
// 			c.Abort()
// 			return

// 		}
// 		c.Set("tokenClaims", tokenClaims)

// 		c.Next()

// 	}

// }

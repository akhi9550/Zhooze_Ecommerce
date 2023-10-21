package helper

import (
	"Zhooze/config"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// func ExtractUserIDFromToken(tokenString string) (int, error) {

// 	cfg, _ := config.LoadConfig()
// 	// Parse the token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(cfg.KEY), nil
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	// Check if the token is valid
// 	if !token.Valid {
// 		return 0, errors.New("internal  error")
// 	}

// 	// Extract the user ID from the claims
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return 0, errors.New("internal  error")
// 	}

// 	userIDfloat64, ok := claims["id"].(float64) // Assuming "Id" is the key for user ID
// 	if !ok {
// 		return 0, errors.New("internal error")
// 	}
// 	userID := int(userIDfloat64)

//		return userID, nil
//	}
func GetTokenFromHeader(header string) string {
	// Example header format: "Bearer <token>"
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}

	return header
}
func ExtractUserIDFromToken(tokenString string) (int, string, error) {
	cfg, _ := config.LoadConfig()

	token, err := jwt.ParseWithClaims(tokenString, &AuthUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(cfg.KEY), nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*AuthUserClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	return claims.Id, claims.Email, nil

}

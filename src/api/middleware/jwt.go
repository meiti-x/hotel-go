package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *fiber.Ctx) error {
	token := c.Get("Token")
	if token == "" {
		return fmt.Errorf("unauthorized")
	}
	if err := parseToken(token); err != nil {
		return fmt.Errorf("unauthorized")
	}

	fmt.Println(token)
	return nil
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable not set")
		}
		fmt.Println("secret", secret)
		return []byte(secret), nil
	})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return nil
	}
	return fmt.Errorf("unauthorized")
}

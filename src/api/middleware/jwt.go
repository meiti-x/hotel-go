package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meiti-x/hotel-go/src/db"
)

func JWTAuth(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Token"]

		if !ok {
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(strings.Join(token, ""))
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		if !ok {
			return fmt.Errorf("invalid expires claim")
		}

		if err != nil {
			return fmt.Errorf("invalid expires format: %v", err)
		}
		expiresFloat, ok := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("Token expired")
		}
		fmt.Println("expires", expires)
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		// set the current authenicated user to context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	fmt.Println("secret", secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("secret", secret)
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable not set")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}

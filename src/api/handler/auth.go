package handler

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userstore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userstore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("Invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("Invalid credentials pa")
	}
	fmt.Println("authenicated")

	return c.JSON(&AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	})
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secert := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secert))
	if err != nil {
		fmt.Println("faild to sign token with secert")
	}
	return tokenString
}

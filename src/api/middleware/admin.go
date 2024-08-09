package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok || !user.IsAdmin {
		return fmt.Errorf("not Authorized")
	}
	return nil

}

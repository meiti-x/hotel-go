package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userstore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userstore,
	}
}

func (h *UserHandler) GetUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) GetUsersHandler(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]interface{}{
		"deleted": userId,
	})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParmas
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return nil
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	insertdUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertdUser)
}

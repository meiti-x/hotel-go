package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/api"
	"github.com/meiti-x/hotel-go/src/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store db.Store
}

func NewBookingHanlder(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: *store,
	}
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := api.GetAuthenicatedUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(&genericResp{
			Type: "error",
			Msg:  "unathorized",
		})
	}

	if _, err := h.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "message",
		Msg:  "okay",
	})
}

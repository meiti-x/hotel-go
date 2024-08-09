package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	store db.Store
}

func NewHotelHanlder(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: *store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")

	rooms, err := h.store.Hotel.GetById(c.Context(), hotelID)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oid, err := db.ToObjectID(hotelID)

	filter := bson.D{{"hotelID", oid}}
	rooms, err := h.store.Room.GetAllRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

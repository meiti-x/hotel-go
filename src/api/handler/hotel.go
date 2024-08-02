package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHanlder(hotelStore db.HotelStore, roomStore db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oid, err := db.ToObjectID(hotelID)

	filter := bson.D{{"hotelID", oid}}
	rooms, err := h.roomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

// func (h *HotelHandler) GetHotelHandler(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	user, err := h.hotelStore.GetById(c.Context(), id)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return c.JSON(map[string]string{"error": "not found"})
// 		}
// 		return err
// 	}
// 	return c.JSON(user)
// }

// func (h *HotelHandler) HandleDelete(c *fiber.Ctx) error {
// 	hotelId := c.Params("id")
// 	if err := h.hotelStore.Delete(c.Context(), hotelId); err != nil {
// 		return err
// 	}
// 	return c.JSON(map[string]interface{}{
// 		"deleted": hotelId,
// 	})
// }

// func (h *HotelHandler) HandleUpdateHotel(c *fiber.Ctx) error {
// 	var params bson.M
// 	hotelId := c.Params("id")
// 	if err := c.BodyParser(&params); err != nil {
// 		return err
// 	}
// 	oid, err := db.ToObjectID(hotelId)
// 	if err != nil {
// 		return err
// 	}

// 	filter := bson.D{{"_id", oid}}
// 	update := bson.D{{"$set", params}}

// 	if err := h.hotelStore.Update(c.Context(), filter, update); err != nil {
// 		return err
// 	}
// 	return c.JSON(map[string]interface{}{
// 		"params": params,
// 		"filter": filter,
// 		"update": update,
// 	})
// }

// func (h *HotelHandler) HandleCreateHotel(c *fiber.Ctx) error {
// 	var params types.CreateHotelParmas
// 	if err := c.BodyParser(&params); err != nil {
// 		return err
// 	}
// 	user, err := types.NewUserFromParams(params)
// 	if err != nil {
// 		return nil
// 	}

// 	if errors := params.Validate(); len(errors) > 0 {
// 		return c.JSON(errors)
// 	}
// 	insertdUser, err := h.hotelStore.Insert(c.Context(), user)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(insertdUser)
// }

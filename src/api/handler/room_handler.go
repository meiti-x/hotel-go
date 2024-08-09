package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store db.Store
}

func NewRoomHanlder(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: *store,
	}
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()

	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannt book a room in the past")
	}
	return nil
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}
	roomID := c.Params("id")
	roomOID, err := db.ToObjectID(roomID)
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internal server error",
		})
	}

	ok, err = h.isRoomAvailbeForBook(c.Context(), roomOID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "Room has Aleardy taken by someone else",
		})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", booking)
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailbeForBook(ctx context.Context, roomId primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomId,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$gte": params.FromDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	for _, cur := range bookings {
		fmt.Println(cur)
	}

	ok := len(bookings) == 0
	if len(bookings) > 0 {
		return false, nil
	}
	return ok, nil
}

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
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

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams

	if err := c.BodyParser(&params); err != nil {
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

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	fmt.Printf("%+v\n", booking)
	return nil
}

package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/api/handler"
	"github.com/meiti-x/hotel-go/src/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userHandler  = handler.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelHandler = handler.NewHotelHanlder(hotelStore, roomStore)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
	)

	// user handler
	apiv1.Get("/users", userHandler.GetUsersHandler)
	apiv1.Get("/users/:id", userHandler.GetUserHandler)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/users", userHandler.HandleCreateUser)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)

	// hotel handler
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetHotelRooms)
	app.Listen(":5000")
}

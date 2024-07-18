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

const (
	dbUri    = "mongodb://localhost:27017"
	dbName   = "hotel_reservation"
	userColl = "users"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := handler.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	api := app.Group("/api/v1")

	api.Get("/users", userHandler.GetUsersHandler)
	api.Get("/users/:id", userHandler.GetUserHandler)
	api.Delete("/users/:id", userHandler.HandleDeleteUser)
	api.Post("/users", userHandler.HandleCreateUser)
	api.Put("/users/:id", userHandler.HandleUpdateUser)
	app.Listen(":5000")
}

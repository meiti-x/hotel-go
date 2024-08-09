package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/api/handler"
	"github.com/meiti-x/hotel-go/src/api/middleware"
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
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Room:    roomStore,
			Hotel:   hotelStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = handler.NewUserHandler(db.NewMongoUserStore(client))
		hotelHandler   = handler.NewHotelHanlder(store)
		authHandler    = handler.NewAuthHandler(userStore)
		roomHandler    = handler.NewRoomHanlder(store)
		BookingHandler = handler.NewBookingHanlder(store)
		app            = fiber.New(config)
		apiv1NoAuth    = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuth(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// auth handlers
	apiv1NoAuth.Post("/auth/login", authHandler.HandleAuth)

	// user handlers
	apiv1.Get("/users", userHandler.GetUsersHandler)
	apiv1.Get("/users/:id", userHandler.GetUserHandler)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/users", userHandler.HandleCreateUser)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)

	// hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetHotelRooms)

	// room handler
	apiv1.Post("/rooms", roomHandler.HandleGetRooms)
	apiv1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)

	// TODO cancel booking

	// booking
	apiv1.Get("/bookings/:id", BookingHandler.HandleGetBooking)
	apiv1.Get("/bookings/:id/cancel", BookingHandler.HandleCancelBooking)

	// admin routes
	admin.Get("/bookings", BookingHandler.HandleGetBookings)

	app.Listen(":5000")
}

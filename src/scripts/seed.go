package main

import (
	"context"
	"log"

	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	hotelStore db.HotelStore
	userStore  db.UserStore
	roomStore  db.RoomStore
	ctx        = context.Background()
)

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	userStore = db.NewMongoUserStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}

func seedUser(isAdmin bool, user *types.User) {
	user, err := types.NewUserFromParams(types.CreateUserParmas{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  "Pass",
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rate int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rate,
	}
	rooms := []types.Room{
		{
			Type:      types.DeluxRoomType,
			BasePrice: 200,
			Size:      "small",
		},
		{
			Type:      types.SinglePersonRoomType,
			BasePrice: 99.9,
			Size:      "large",
		},
		{
			Type:      types.DoublePersonRoomType,
			BasePrice: 120,
			Size:      "kingsize",
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Ilam", "The Crown of Zagros", 5)
	seedUser(true, &types.User{
		Email:     "admin@mahdi.com",
		FirstName: "Mahdi",
		LastName:  "M",
	})
	seedUser(false, &types.User{
		Email:     "user@mahdi.com",
		FirstName: "Mahdi",
		LastName:  "M",
	})
}

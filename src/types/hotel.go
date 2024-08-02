package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type RoomType int

const (
	_ RoomType = iota
	SinglePersonRoomType
	DoublePersonRoomType
	SeaSideRoomType
	DeluxRoomType
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

type CreateHotelParmas struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rooms    string `json:"rooms"`
}

// func NewHotelFromParams(params CreateHotelParmas) (*Hotel, error) {
// 	return &Hotel{
// 		Name:     params.Name,
// 		Location: params.Location,
// 		Rooms:    params.Rooms.(primitive.ObjectID),
// 	}, nil
// }

package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DBNAME           = "hotel-reservation"
	TestDBNAME       = "hotel_reservation_test"
	USER_COLLECTION  = "users"
	HOTEL_COLLECTION = "hotels"
	ROOM_COLLECTION  = "rooms"
	DB_URI           = "mongodb://localhost:27017"
)

func ToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}

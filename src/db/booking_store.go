package db

import (
	"context"
	"fmt"

	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error)
	GetBookingByID(ctx context.Context, id string) (*types.Booking, error)
	UpdateBooking(ctx context.Context, id string, data bson.M) (*types.Booking, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(BOOKING_COLLECTION),
	}
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	var booking types.Booking
	oid, err := ToObjectID(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) (*types.Booking, error) {
	var booking types.Booking
	fmt.Println(id)
	oid, err := ToObjectID(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(oid)
	_, err = s.coll.UpdateByID(ctx, oid, bson.M{"$set": update})

	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := curr.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

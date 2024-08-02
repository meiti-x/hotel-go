package db

import (
	"context"
	"fmt"

	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	Dropper

	InsertHotel(ctx context.Context, user *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filters bson.D, data bson.D) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	return nil
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(HOTEL_COLLECTION),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filters bson.D, data bson.D) error {
	z, err := s.coll.UpdateOne(ctx, filters, data)
	fmt.Sprintln("z")
	fmt.Sprintln(z)
	fmt.Sprintln("z")
	if err != nil {
		return err
	}
	return nil
}

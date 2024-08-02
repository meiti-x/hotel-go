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

	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filters bson.D, data bson.D) error
	// GetById(ctx context.Context, id string) (*types.Hotel, error)
	GetAll(ctx context.Context) ([]*types.Hotel, error)
	GetAllRooms(ctx context.Context, hotelID string) ([]*types.Hotel, error)
	// Delete(ctx context.Context, id string) error
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

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filters bson.D, data bson.D) error {
	z, err := s.coll.UpdateOne(ctx, filters, data)
	fmt.Sprintln("z")
	fmt.Sprintln(z)
	fmt.Sprintln("z")
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelStore) GetAll(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

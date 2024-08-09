package db

import (
	"context"

	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	Dropper

	InsertRoom(ctx context.Context, user *types.Room) (*types.Room, error)
	GetAllRooms(ctx context.Context, filter bson.D) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(ROOM_COLLECTION),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetAllRooms(ctx context.Context, filter bson.D) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)
	println(room.HotelID.Hex())
	// update the hotel with this room id
	filter := bson.D{{"_id", room.HotelID}}
	update := bson.D{
		{"$push", bson.D{
			{"rooms", room.ID},
		}},
	}
	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

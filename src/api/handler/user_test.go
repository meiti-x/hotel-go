package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	db.UserStore
}

const (
	dbUri = "mongodb://localhost:27017"
)

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testDB{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func (db *testDB) teardown(t *testing.T) {
	if err := db.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func assertEqual(t *testing.T, fieldName, expected, actual string) {
	if expected != actual {
		t.Errorf("expected %s %s but got %s", fieldName, expected, actual)
	}
}

func TestPostUsesr(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(db.UserStore)
	app.Post("/", UserHandler.HandleCreateUser)

	params := types.CreateUserParmas{
		Email:     "justfortest@mail.com",
		FirstName: "string",
		LastName:  "string",
		Password:  "12345678",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting user id to be set")
	}

	assertEqual(t, "FirstName", params.FirstName, user.FirstName)
	assertEqual(t, "LastName", params.LastName, user.LastName)
	assertEqual(t, "Email", params.Email, user.Email)
}

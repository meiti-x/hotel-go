package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/meiti-x/hotel-go/src/db"
	"github.com/meiti-x/hotel-go/src/types"
)

func InsertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParmas{
		Email:     "test@test.com",
		FirstName: "test",
		LastName:  "testy",
		Password:  "test",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	insertedUser := InsertTestUser(t, db)

	app := fiber.New()
	authHandler := NewAuthHandler(db.UserStore)

	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "test@test.com",
		Password: "test",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	if authResp.Token == "" {
		t.Fatalf("Expected the jwt token to be in the out response")
	}
	// we omit the encrypted password from user object, so we compare it with auth resp
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected user to be same inserted user")
	}
}

func TestAuthWithWrongPass(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	InsertTestUser(t, db)

	app := fiber.New()
	authHandler := NewAuthHandler(db.UserStore)

	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "test@test.com",
		Password: "NOTCORRECT",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("Exptected gen response type to get error but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("Exptected gen response Msg to get 'invalid credentials' but got %s", genResp.Msg)
	}
}

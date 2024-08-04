package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	firstNameMinLen = 3
	lastNameMinLen  = 4
	passwordMinLen  = 7
)

type CreateUserParmas struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c *CreateUserParmas) Validate() map[string]string {
	errors := map[string]string{}
	if len(c.FirstName) < firstNameMinLen {
		errors["firstName"] = fmt.Sprintf("first name length most be greater than %d", firstNameMinLen)
	}

	if len(c.LastName) < lastNameMinLen {
		errors["lastName"] = fmt.Sprintf("last name length most be greater than %d", lastNameMinLen)
	}

	if len(c.Password) < passwordMinLen {
		errors["password"] = fmt.Sprintf("password length most be greater than %d", passwordMinLen)
	}

	if !isValidEmail(c.Email) {
		errors["email"] = fmt.Sprintf("PLEASE ENTER A VALID EMAIL")
	}

	return errors
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParmas) (*User, error) {
	encbw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encbw),
	}, nil
}

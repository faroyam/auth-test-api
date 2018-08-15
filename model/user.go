package model

import (
	"strings"
	"unicode/utf8"

	"gopkg.in/mgo.v2/bson"
)

// User represents user model
type User struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Login          string        `json:"login,omitempty" bson:"login,omitempty"`
	HashedPassword string        `json:"password,omitempty" bson:"password,omitempty"`
	Email          string        `json:"email" bson:"email"`
}

// Validate validates
func (u *User) Validate() bool {
	username := utf8.RuneCountInString(u.Login) >= 8
	email := strings.Contains(u.Email, "@")
	password := utf8.RuneCountInString(u.HashedPassword) >= 8

	return username && email && password
}

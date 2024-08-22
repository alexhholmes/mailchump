package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

func (u *User) Validate() error {
	// TODO implement
	return nil
}

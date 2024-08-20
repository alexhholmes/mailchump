package model

import "github.com/google/uuid"

type Author struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

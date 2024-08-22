package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (i SignUpInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&i.Password, validation.Required, validation.Length(8, 100), is.PrintableASCII),
	)
}

func (i SignInInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Username, validation.Required),
		validation.Field(&i.Password, validation.Required),
	)
}
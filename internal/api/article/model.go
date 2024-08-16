package article

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Article struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Author    string    `db:"author" json:"author"`
	Content   string    `db:"content" json:"content"`
	ImageURL  *string   `db:"image_url" json:"image_url,omitempty"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateArticleInput struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url,omitempty"`
}

type UpdateArticleInput struct {
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

func (i CreateArticleInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&i.Content, validation.Required),
		validation.Field(&i.ImageURL, validation.NilOrNotEmpty, is.URL),
	)
}

func (i UpdateArticleInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&i.Content, validation.NilOrNotEmpty),
		validation.Field(&i.ImageURL, validation.NilOrNotEmpty, is.URL),
	)
}
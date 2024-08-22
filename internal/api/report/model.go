package report

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Report struct {
	ID           int64     `db:"id" json:"id"`
	ReporterName string    `db:"reporter_name" json:"reporter_name"`
	Email        string    `db:"email" json:"email"`
	Description  string    `db:"description" json:"description"`
	DataURL      *string   `db:"data_url" json:"data_url,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type CreateReportInput struct {
	ReporterName string  `json:"reporter_name"`
	Email        string  `json:"email"`
	Description  string  `json:"description"`
	DataURL      *string `json:"data_url,omitempty"`
}

type UpdateReportInput struct {
	ReporterName *string `json:"reporter_name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Description  *string `json:"description,omitempty"`
	DataURL      *string `json:"data_url,omitempty"`
}

func (i CreateReportInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ReporterName, validation.Required, validation.Length(1, 255)),
		validation.Field(&i.Email, validation.Required, is.Email),
		validation.Field(&i.Description, validation.Required),
		validation.Field(&i.DataURL, validation.NilOrNotEmpty, is.URL),
	)
}

func (i UpdateReportInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ReporterName, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&i.Email, validation.NilOrNotEmpty, is.Email),
		validation.Field(&i.Description, validation.NilOrNotEmpty),
		validation.Field(&i.DataURL, validation.NilOrNotEmpty, is.URL),
	)
}
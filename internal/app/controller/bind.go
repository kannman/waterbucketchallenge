package controller

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)

func Bind(r *http.Request, v render.Binder) error {
	if err := render.Bind(r, v); err != nil {
		return err
	}
	if validatable, ok := v.(validation.Validatable); ok {
		return validatable.Validate()
	}
	return nil
}

package api

import (
	"net/http"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/web"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,gte=5,lte=32"`
	Password string `json:"password" validate:"required,gte=7,lte=32"`
}

type SignupResponse struct{}

func handleSignup(d db.Executor, v *web.Validator) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		data, err := web.DecodeValidated[SignupRequest](v, r)
		if err != nil {
			return err
		}

		_, err = resource.CreateUser(r.Context(), d, resource.NewUser{
			Name:     data.Name,
			Email:    data.Email,
			Password: data.Password,
		})
		if err != nil {
			return err
		}

		return web.Encode(w, http.StatusCreated, &SignupResponse{})
	}
}

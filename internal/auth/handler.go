package auth

import (
	"fmt"
	"net/http"
	"project/configs"
	"project/pkg/request"
	"project/pkg/response"
)

type handler struct {
	authConfig *configs.AuthConfig
}

func (h *handler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		fmt.Println(body)
		resp := &LoginResponse{Token: h.authConfig.Secret}
		response.Json(w, resp, http.StatusOK)
	}
}

func (h *handler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}

		fmt.Println(body)
		resp := &RegisterResponse{"Ok"}
		response.Json(w, resp, http.StatusOK)
	}
}

func NewAuthHandler(r *http.ServeMux, ac *configs.AuthConfig) {
	h := handler{ac}
	r.HandleFunc("POST /auth/login", h.login())
	r.HandleFunc("POST /auth/register", h.register())
}

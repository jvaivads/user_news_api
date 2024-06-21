package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"user_news_api/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func (uc *UserController) registerRoutes(router chi.Router) {
	router.Post("/notifications", uc.handleNotifyUser)
}

// UserNotifier is an abstraction for services.UserNotifierService making it mockeable
type UserNotifier interface {
	Notify(context.Context, string, string) error
}

func SetUserController(router chi.Router, service UserNotifier) {
	controller := &UserController{service: service}

	controller.registerRoutes(router)
}

type UserController struct {
	service UserNotifier
}

type NotifyUserRequestPayload struct {
	UserEmail   string `json:"user_email" validate:"required,email"`
	MessageType string `json:"message_type" validate:"required"`
}

func (uc *UserController) handleNotifyUser(w http.ResponseWriter, r *http.Request) {
	var payload NotifyUserRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("error marshalling request body due to: %s", err.Error()), http.StatusBadRequest)

		return
	}

	if err := validator.New().Struct(payload); err != nil {
		http.Error(w, fmt.Sprintf("request validation fails due to: %s", err.Error()), http.StatusBadRequest)

		return
	}

	if err := uc.service.Notify(r.Context(), payload.UserEmail, payload.MessageType); err != nil {
		if errors.Is(err, services.ErrLimitExceeded) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)

			return
		}

		log.Printf("error notifying user: %s", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

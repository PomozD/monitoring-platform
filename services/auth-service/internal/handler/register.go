package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/application/auth"
	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type RegisterHandler struct {
	registerUser *auth.RegisterUser
}

func NewRegisterHandler(
	registerUser *auth.RegisterUser,
) *RegisterHandler {
	return &RegisterHandler{
		registerUser: registerUser,
	}
}

func (h *RegisterHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"Invalid request body",
		)
		return
	}

	user, err := h.registerUser.Execute(
		r.Context(),
		auth.RegisterInput{
			Email:    request.Email,
			Password: request.Password,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmail):
			writeError(
				w,
				http.StatusBadRequest,
				"INVALID_EMAIL",
				"Invalid email",
			)

		case errors.Is(err, domain.ErrInvalidPassword):
			writeError(
				w,
				http.StatusBadRequest,
				"INVALID_PASSWORD",
				"Invalid password",
			)

		case errors.Is(err, domain.ErrUserAlreadyExists):
			writeError(
				w,
				http.StatusConflict,
				"USER_ALREADY_EXISTS",
				"User already exists",
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				"Internal server error",
			)
		}

		return
	}

	response := RegisterResponse{
		ID:     user.ID.String(),
		Email:  user.Email,
		Status: string(user.Status),
	}

	writeJSON(
		w,
		http.StatusCreated,
		response,
	)
}

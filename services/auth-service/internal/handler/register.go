package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/application/auth"
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
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
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
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	response := RegisterResponse{
		ID:     user.ID.String(),
		Email:  user.Email,
		Status: string(user.Status),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

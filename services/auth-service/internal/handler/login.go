package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/application/auth"
	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type LoginHandler struct {
	loginUser *auth.LoginUser
}

func NewLoginHandler(loginUser *auth.LoginUser) *LoginHandler {
	return &LoginHandler{
		loginUser: loginUser,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}

	user, err := h.loginUser.Execute(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		ID:     user.ID.String(),
		Email:  user.Email,
		Status: string(user.Status),
	})
}

package controllers

import (
	"myapp/internal/users"
	"myapp/pkg/io"
	"net/http"
)

func (h *HandlerLayer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req users.LoginRequest
	err := io.ReadJSON(r, &req)
	if err != nil {
		io.SendError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user, err := h.LogicLayer.VerifyUserCredentials(req.Username, req.Password)
	if err != nil || user == nil {
		io.SendError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	token, err := h.LogicLayer.CreateSession(user.Username)
	if err != nil {
		io.SendError(w, "Failed to create session or session already exists", http.StatusUnauthorized)
		return
	}
	io.WriteJSON(w, http.StatusOK, users.LoginResponse{Token: token})
}

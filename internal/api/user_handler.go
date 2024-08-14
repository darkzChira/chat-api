package api

import (
	"chat-api/internal/models"
	"chat-api/internal/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.userService.RegisterUser(context.Background(), &user); err != nil {
		if err.Error() == "username already taken" {
			http.Error(w, "Username already taken", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginUser models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tokenString, user, err := h.userService.LoginUser(context.Background(), loginUser.Username, loginUser.Password)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  user.ConvertToResponse(),
	})
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := h.userService.LogoutUser(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to logout user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.userService.GetUser(context.Background(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user.ConvertToResponse())
}

func (h *UserHandler) GetOtherUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currentUserID := vars["id"]

	users, err := h.userService.GetOtherUsers(context.Background(), currentUserID)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

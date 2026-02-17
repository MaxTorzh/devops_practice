package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_microservices/internal/models"
	"go_microservices/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", h.listUsers)
	mux.HandleFunc("POST /users", h.createUser)
	mux.HandleFunc("GET /users/{id}", h.getUser)
	mux.HandleFunc("PUT /users/{id}", h.updateUser)
	mux.HandleFunc("DELETE /users/{id}", h.deleteUser)
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := h.userService.GetAll(ctx, page, limit)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	total, _ := h.userService.Count(ctx)

	response := map[string]interface{}{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userService.Create(ctx, &req)
	if err != nil {
		if err.Error() == "user with email already exists" {
			h.respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		h.respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userService.Update(ctx, id, &req)
	if err == sql.ErrNoRows {
		h.respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.userService.Delete(ctx, id)
	if err == sql.ErrNoRows {
		h.respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *UserHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Status:  code,
	})
}

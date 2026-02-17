package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go_api/internal/config"
	"go_api/internal/models"
	"go_api/internal/repository/postgres"
)

type UserHandler struct {
	repo *postgres.UserRepository
	cfg  *config.Config
}

func NewUserHandler(db *sql.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		repo: postgres.NewUserRepository(db),
		cfg:  cfg,
	}
}

func (h *UserHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.handleHome)
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/", h.handleUserByID)

	return mux
}

func (h *UserHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	count, _ := h.repo.Count()

	response := map[string]interface{}{
		"message":     "Connected to PostgreSQL!",
		"users_count": count,
		"version":     "v1.0.0",
		"endpoints": []string{
			"GET /users",
			"GET /users/{id}",
			"POST /users",
			"PUT /users/{id}",
			"DELETE /users/{id}",
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		h.respondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAll()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	h.respondWithJSON(w, http.StatusOK, users)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id int) {
	user, err := h.repo.GetByID(id)
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

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" {
		h.respondWithError(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	user, err := h.repo.Create(req.Name, req.Email)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.Update(id, req.Name, req.Email); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated"})
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.repo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
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

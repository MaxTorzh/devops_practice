package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go_volume/internal/config"
	"go_volume/internal/models"
	"go_volume/internal/service"
)

type FileHandler struct {
	service *service.FileService
	config  *config.Config
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		service: service.NewFileService(cfg.DataDir),
		config:  cfg,
	}
}

func (h *FileHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.handleHome)
	mux.HandleFunc("/files", h.handleFiles)
	mux.HandleFunc("/files/", h.handleFile)
	mux.HandleFunc("/info/", h.handleFileInfo)

	return mux
}

func (h *FileHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"message":  "File Service with Volume Mount",
		"data_dir": h.config.DataDir,
		"endpoints": []string{
			"GET    /files          - list all files",
			"GET    /files/{name}   - get file content",
			"POST   /files/{name}   - create file",
			"PUT    /files/{name}   - update file",
			"DELETE /files/{name}   - delete file",
			"GET    /info/{name}    - get file info",
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *FileHandler) handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	files, err := h.service.ListFiles()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, files)
}

func (h *FileHandler) handleFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/files/")
	if filename == "" {
		h.respondWithError(w, http.StatusBadRequest, "filename required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getFile(w, filename)
	case http.MethodPost:
		h.createFile(w, r, filename)
	case http.MethodPut:
		h.updateFile(w, r, filename)
	case http.MethodDelete:
		h.deleteFile(w, filename)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *FileHandler) handleFileInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/info/")
	if filename == "" {
		h.respondWithError(w, http.StatusBadRequest, "filename required")
		return
	}

	info, err := h.service.GetFileInfo(filename)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, info)
}

func (h *FileHandler) getFile(w http.ResponseWriter, filename string) {
	content, err := h.service.GetFile(filename)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func (h *FileHandler) createFile(w http.ResponseWriter, r *http.Request, filename string) {
	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxUploadSize)

	content, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "file too large")
		return
	}

	if err := h.service.CreateFile(filename, content); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("file %s created", filename),
	}
	h.respondWithJSON(w, http.StatusCreated, response)
}

func (h *FileHandler) updateFile(w http.ResponseWriter, r *http.Request, filename string) {
	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxUploadSize)

	content, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "file too large")
		return
	}

	if err := h.service.UpdateFile(filename, content); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("file %s updated", filename),
	}
	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *FileHandler) deleteFile(w http.ResponseWriter, filename string) {
	if err := h.service.DeleteFile(filename); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FileHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *FileHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Status:  code,
	})
}

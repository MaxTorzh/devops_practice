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

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.listProducts)
	mux.HandleFunc("POST /products", h.createProduct)
	mux.HandleFunc("GET /products/{id}", h.getProduct)
	mux.HandleFunc("PUT /products/{id}", h.updateProduct)
	mux.HandleFunc("DELETE /products/{id}", h.deleteProduct)
	mux.HandleFunc("PATCH /products/{id}/stock", h.updateStock)
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	products, err := h.productService.GetAll(ctx, page, limit)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	total, _ := h.productService.Count(ctx)

	response := map[string]interface{}{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productService.Create(ctx, &req)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productService.GetByID(ctx, id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if product == nil {
		h.respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productService.Update(ctx, id, &req)
	if err == sql.ErrNoRows {
		h.respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) updateStock(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.productService.UpdateStock(ctx, id, req.Quantity)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Stock updated"})
}

func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.productService.Delete(ctx, id)
	if err == sql.ErrNoRows {
		h.respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Status:  code,
	})
}

package service

import (
	"context"
	"database/sql"
	"fmt"

	"go_microservices/internal/models"
	"go_microservices/internal/repository/postgres"
	"go_microservices/internal/repository/redis"
)

type ProductService struct {
	productRepo *postgres.ProductRepository
	cacheRepo   *redis.CacheRepository
}

func NewProductService(productRepo *postgres.ProductRepository, cacheRepo *redis.CacheRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, redis.ProductListKey(1, 100))

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	cacheKey := redis.ProductKey(id)
	var product models.Product

	err := s.cacheRepo.Get(ctx, cacheKey, &product)
	if err == nil {
		return &product, nil
	}

	productPtr, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if productPtr == nil {
		return nil, nil
	}

	s.cacheRepo.Set(ctx, cacheKey, productPtr)

	return productPtr, nil
}

func (s *ProductService) GetAll(ctx context.Context, page, limit int) ([]models.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	cacheKey := redis.ProductListKey(page, limit)
	var products []models.Product

	err := s.cacheRepo.Get(ctx, cacheKey, &products)
	if err == nil {
		return products, nil
	}

	products, err = s.productRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.Set(ctx, cacheKey, products)

	return products, nil
}

func (s *ProductService) Update(ctx context.Context, id int, req *models.UpdateProductRequest) (*models.Product, error) {
	existing, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, sql.ErrNoRows
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Update(id, product); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, redis.ProductKey(id), redis.ProductListKey(1, 100))

	return s.productRepo.GetByID(id)
}

func (s *ProductService) UpdateStock(ctx context.Context, id, quantity int) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return sql.ErrNoRows
	}
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, quantity)
	}

	if err := s.productRepo.UpdateStock(id, quantity); err != nil {
		return err
	}

	s.cacheRepo.Delete(ctx, redis.ProductKey(id), redis.ProductListKey(1, 100))

	return nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	s.cacheRepo.Delete(ctx, redis.ProductKey(id), redis.ProductListKey(1, 100))

	return nil
}

func (s *ProductService) Count(ctx context.Context) (int, error) {
	return s.productRepo.Count()
}

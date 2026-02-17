package service

import (
	"context"
	"database/sql"
	"fmt"

	"go_microservices/internal/models"
	"go_microservices/internal/repository/postgres"
	"go_microservices/internal/repository/redis"
)

type UserService struct {
	userRepo  *postgres.UserRepository
	cacheRepo *redis.CacheRepository
}

func NewUserService(userRepo *postgres.UserRepository, cacheRepo *redis.CacheRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, redis.UserListKey(1, 100))

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	cacheKey := redis.UserKey(id)
	var user models.User

	err := s.cacheRepo.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	userPtr, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if userPtr == nil {
		return nil, nil
	}

	s.cacheRepo.Set(ctx, cacheKey, userPtr)

	return userPtr, nil
}

func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	cacheKey := redis.UserListKey(page, limit)
	var users []models.User

	err := s.cacheRepo.Get(ctx, cacheKey, &users)
	if err == nil {
		return users, nil
	}

	users, err = s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.Set(ctx, cacheKey, users)

	return users, nil
}

func (s *UserService) Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	existing, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, sql.ErrNoRows
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userRepo.Update(id, user); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, redis.UserKey(id), redis.UserListKey(1, 100))

	return s.userRepo.GetByID(id)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	s.cacheRepo.Delete(ctx, redis.UserKey(id), redis.UserListKey(1, 100))

	return nil
}

func (s *UserService) Count(ctx context.Context) (int, error) {
	return s.userRepo.Count()
}

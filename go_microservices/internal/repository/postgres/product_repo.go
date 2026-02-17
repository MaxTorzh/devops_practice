package postgres

import (
	"database/sql"
	"fmt"

	"go_microservices/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `
        INSERT INTO products (name, description, price, stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	return r.db.QueryRow(
		query, product.Name, product.Description, product.Price, product.Stock,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        WHERE id = $1
    `

	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description,
		&product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &product, err
}

func (r *ProductRepository) GetAll(limit, offset int) ([]models.Product, error) {
	query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if products == nil {
		return []models.Product{}, nil
	}
	return products, nil
}

func (r *ProductRepository) Update(id int, product *models.Product) error {
	query := `
        UPDATE products
        SET name = COALESCE($1, name),
            description = COALESCE($2, description),
            price = COALESCE($3, price),
            stock = COALESCE($4, stock),
            updated_at = NOW()
        WHERE id = $5
        RETURNING updated_at
    `

	var name, description *string
	var price *float64
	var stock *int

	if product.Name != "" {
		name = &product.Name
	}
	if product.Description != "" {
		description = &product.Description
	}
	if product.Price > 0 {
		price = &product.Price
	}
	if product.Stock >= 0 {
		stock = &product.Stock
	}

	return r.db.QueryRow(query, name, description, price, stock, id).Scan(&product.UpdatedAt)
}

func (r *ProductRepository) UpdateStock(id, quantity int) error {
	query := `
        UPDATE products
        SET stock = stock - $1,
            updated_at = NOW()
        WHERE id = $2 AND stock >= $1
        RETURNING stock
    `

	var newStock int
	err := r.db.QueryRow(query, quantity, id).Scan(&newStock)
	if err == sql.ErrNoRows {
		return fmt.Errorf("insufficient stock or product not found")
	}
	return err
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	return count, err
}

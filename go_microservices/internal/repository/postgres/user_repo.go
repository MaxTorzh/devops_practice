package postgres

import (
	"database/sql"

	"go_microservices/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	return r.db.QueryRow(query, user.Name, user.Email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetAll(limit, offset int) ([]models.User, error) {
	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if users == nil {
		return []models.User{}, nil
	}
	return users, nil
}

func (r *UserRepository) Update(id int, user *models.User) error {
	query := `
        UPDATE users
        SET name = COALESCE($1, name),
            email = COALESCE($2, email),
            updated_at = NOW()
        WHERE id = $3
        RETURNING updated_at
    `

	var name, email *string
	if user.Name != "" {
		name = &user.Name
	}
	if user.Email != "" {
		email = &user.Email
	}

	return r.db.QueryRow(query, name, email, id).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

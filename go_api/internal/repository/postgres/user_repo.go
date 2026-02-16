package postgres

import (
	"database/sql"

	"go_api/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
    rows, err := r.db.Query("SELECT id, name, email, created_at FROM users ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
	if users == nil {
		return []models.User{}, nil
	}
    return users, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
    var u models.User
    err := r.db.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE id = $1", id,
    ).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &u, err
}

func (r *UserRepository) Create(name, email string) (*models.User, error) {
    var u models.User
    err := r.db.QueryRow(
        "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at",
        name, email,
    ).Scan(&u.ID, &u.CreatedAt)

    if err != nil {
        return nil, err
    }

    u.Name = name
    u.Email = email
    return &u, nil
}

func (r *UserRepository) Update(id int, name, email string) error {
    result, err := r.db.Exec(
        "UPDATE users SET name = COALESCE($1, name), email = COALESCE($2, email) WHERE id = $3",
        nullIfEmpty(name), nullIfEmpty(email), id,
    )
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return sql.ErrNoRows
    }
    return nil
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

func nullIfEmpty(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}
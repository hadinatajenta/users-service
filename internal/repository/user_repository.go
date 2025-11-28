package repository

import (
	"context"
	"database/sql"
	"fmt"

	"users-service/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password, created_at
	`

	var created entity.User
	if err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).
		Scan(&created.ID, &created.Name, &created.Email, &created.Password, &created.CreatedAt); err != nil {
		return entity.User{}, fmt.Errorf("insert user: %w", err)
	}

	return created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (entity.User, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	if err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user %d not found", id)
		}
		return entity.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]entity.User, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

package users

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository defines persistence behavior for users.
type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	List(ctx context.Context) ([]User, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, user User) (User, error) {
	const query = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password, created_at
	`

	var created User
	if err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).
		Scan(&created.ID, &created.Name, &created.Email, &created.Password, &created.CreatedAt); err != nil {
		return User{}, fmt.Errorf("insert user: %w", err)
	}

	return created, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (User, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	var user User
	if err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user %d not found", id)
		}
		return User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *postgresRepository) List(ctx context.Context) ([]User, error) {
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

	var users []User
	for rows.Next() {
		var u User
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

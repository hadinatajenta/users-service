package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"

	"users-service/internal/config"
)

func NewPostgresConnection(cfg config.DBConfig) (*sql.DB, error) {
	u := url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Name,
		User:   url.UserPassword(cfg.User, cfg.Password),
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	dsn := u.String()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		fmt.Printf("Error getting current database: %v\n", err)
	} else {
		fmt.Printf("CURRENTLY CONNECTED TO DATABASE: '%s'\n", currentDB)
	}

	var currentUser string
	err = db.QueryRow("SELECT current_user").Scan(&currentUser)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	} else {
		fmt.Printf("CURRENT USER: '%s'\n", currentUser)
	}
	return db, nil
}

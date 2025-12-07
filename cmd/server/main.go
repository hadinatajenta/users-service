package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"users-service/internal/config"
	"users-service/internal/db"
	"users-service/internal/users"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer database.Close()

	userRepo := users.NewRepository(database)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	router := gin.Default()
	users.RegisterRoutes(router, userHandler)

	addr := ":" + cfg.AppPort
	log.Printf("starting HTTP server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

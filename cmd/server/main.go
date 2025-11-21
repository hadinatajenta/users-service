package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"users-service/internal/config"
	"users-service/internal/db"
	handler "users-service/internal/delivery/http"
	"users-service/internal/delivery/route"
	"users-service/internal/repository"
	"users-service/internal/service"
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

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	route.Register(router, userHandler)

	addr := ":" + cfg.AppPort
	log.Printf("starting HTTP server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

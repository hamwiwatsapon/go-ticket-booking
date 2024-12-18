package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/database"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/handler"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/repository"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/service"
	"github.com/hamwiwatsapon/go-ticket-booking/pkg/middleware"
)

func main() {
	// Load database configuration
	dbConfig := database.LoadPostgresConfig()

	// Initialize database
	db, err := database.InitPostgres(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := setupRouter(userHandler)

	// Start server
	if err := router.Run(":3333"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh-token", userHandler.RefreshToken)
	}

	protected := router.Group("/api/v1/protected")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// Example of a route only accessible to admin
		protected.GET("/admin", middleware.RoleMiddleware("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome, Admin!",
			})
		})

		// Other protected routes can be added here
	}
	return router
}

package main

import (
	"log"
	"net/http"
	"os"
	"study-quest-backend/internal/handler"
	"study-quest-backend/internal/repository"
	"study-quest-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Repositories (In-Memory for demo)
	taskRepo := repository.NewMemoryTaskRepository()
	userRepo := repository.NewMemoryUserRepository()

	// Initialize Services
	taskService := service.NewTaskService(taskRepo, userRepo)
	
	// Initialize Handlers
	h := handler.NewHandler(taskService)

	// Setup Router
	r := gin.Default()
	r.Use(corsMiddleware())

	// Serve Web Demo
	webDir := "../web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = "./web"
	}
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = "../../web"
	}
	r.Static("/web", webDir)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	api := r.Group("/api/v1")
	{
		api.GET("/config/init", h.GetAppConfig)
		
		// Tasks
		api.GET("/tasks/today", h.GetTodayTasks)
		api.GET("/tasks/pending", h.GetPendingTasks)
		api.POST("/tasks/create", h.CreateTask)
		api.POST("/tasks/submit", h.SubmitTask)
		api.POST("/tasks/approve", h.ApproveTask)

		// Profile
		api.GET("/profile", h.GetProfile)
	}

	// Start Server
	port := "8080"
	log.Printf("Running in DEMO MODE (In-Memory)")
	log.Printf("Server starting on port %s...", port)
	log.Printf("Open http://localhost:%s/web to view the demo", port)
	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}


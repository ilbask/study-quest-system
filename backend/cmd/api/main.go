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
	sessionRepo := repository.NewMemorySessionRepository()

	// Initialize Services
	taskService := service.NewTaskService(taskRepo, userRepo)
	authService := service.NewAuthService(userRepo, sessionRepo)
	
	// Initialize Handlers
	h := handler.NewHandler(taskService, authService)

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
		
		// Auth (public)
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
		api.POST("/auth/logout", h.Logout)
		
		// Ranking (public)
		api.GET("/ranking", h.GetRanking)
	}
	
	// Protected routes (require authentication)
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware(authService))
	{
		// Tasks
		protected.GET("/tasks/today", h.GetTodayTasks)
		protected.GET("/tasks/pending", h.GetPendingTasks)
		protected.POST("/tasks/create", h.CreateTask)
		protected.POST("/tasks/submit", h.SubmitTask)
		protected.POST("/tasks/approve", h.ApproveTask)

		// Profile
		protected.GET("/profile", h.GetProfile)
		
		// Rewards
		protected.POST("/rewards/redeem", h.RedeemReward)
		
		// Students
		protected.GET("/students", h.GetStudentList)
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

func authMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		
		user, err := authService.ValidateSession(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
		// Set user info in context
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}


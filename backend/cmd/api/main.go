package main

import (
	"log"
	"net/http"
	"os"
	"study-quest-backend/internal/config"
	"study-quest-backend/internal/handler"
	"study-quest-backend/internal/repository"
	"study-quest-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Database
	db, err := repository.InitDB(cfg.Database)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		log.Println("Falling back to In-Memory mode...")
		
	// Fallback to memory repositories
	taskRepo := repository.NewMemoryTaskRepository()
	userRepo := repository.NewMemoryUserRepository()
	sessionRepo := repository.NewMemorySessionRepository()
	redemptionRepo := repository.NewMemoryRedemptionRepository()
	rewardRepo := repository.NewMemoryRewardRepository()
	
	taskService := service.NewTaskService(taskRepo, userRepo, redemptionRepo, rewardRepo)
	authService := service.NewAuthService(userRepo, sessionRepo)
	h := handler.NewHandler(taskService, authService)
	
	startServer(h, cfg.Server.Port)
	return
	}

	log.Println("Connected to MySQL successfully!")

	// 3. Auto Migrate Database Schemas
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")

	// 4. Seed Initial Data
	if err := repository.SeedData(db); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	// 5. Initialize Repositories (MySQL)
	taskRepo := repository.NewMySQLTaskRepository(db)
	userRepo := repository.NewMySQLUserRepository(db)
	sessionRepo := repository.NewMySQLSessionRepository(db)
	redemptionRepo := repository.NewMySQLRedemptionRepository(db)
	rewardRepo := repository.NewMySQLRewardRepository(db)

	// 6. Initialize Services
	taskService := service.NewTaskService(taskRepo, userRepo, redemptionRepo, rewardRepo)
	authService := service.NewAuthService(userRepo, sessionRepo)
	
	// 7. Initialize Handlers
	h := handler.NewHandler(taskService, authService)
	
	startServer(h, cfg.Server.Port)
}

func startServer(h *handler.Handler, port string) {

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
	protected.Use(h.AuthMiddleware())
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
	protected.GET("/rewards", h.GetRewards)
	protected.POST("/rewards/redeem", h.RedeemReward)
	protected.GET("/redemptions", h.GetRedemptions)
	
	// Students
	protected.GET("/students", h.GetStudentList)
	}

	// Start Server
	if port == "" {
		port = "8080"
	}
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



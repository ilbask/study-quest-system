package handler

import (
	"log"
	"net/http"
	"study-quest-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	taskService *service.TaskService
	authService *service.AuthService
}

func NewHandler(ts *service.TaskService, as *service.AuthService) *Handler {
	return &Handler{
		taskService: ts,
		authService: as,
	}
}

// Middleware
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		
		user, err := h.authService.ValidateSession(token)
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

func (h *Handler) GetAppConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"review_mode": false,
		"theme":       "default",
	})
}

func (h *Handler) GetTodayTasks(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	tasks, _ := h.taskService.GetTodayTasks(userID.(uint))
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetPendingTasks(c *gin.Context) {
	tasks, _ := h.taskService.GetPendingTasks()
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Points int    `json:"points"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Get current user (parent)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Get user info to find family ID
	user, err := h.taskService.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	
	// Create task and assign to family students
	err = h.taskService.CreateTask(req.Title, req.Points, user.FamilyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func (h *Handler) SubmitTask(c *gin.Context) {
	var req struct {
		TaskID uint `json:"task_id"` // This is actually task_log ID from frontend
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	log.Printf("Submitting task_log ID %d for user %d", req.TaskID, userID.(uint))
	
	err := h.taskService.SubmitTaskByLogID(req.TaskID, userID.(uint))
	if err != nil {
		log.Printf("Error submitting task: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

func (h *Handler) ApproveTask(c *gin.Context) {
	var req struct {
		LogID  uint   `json:"log_id"`
		Action string `json:"action"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	if req.Action == "approve" {
		h.taskService.ApproveTask(req.LogID)
	} else {
		h.taskService.RejectTask(req.LogID)
	}
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	user, _ := h.taskService.GetUserProfile(userID.(uint))
	c.JSON(http.StatusOK, user)
}

func (h *Handler) RedeemReward(c *gin.Context) {
	var req struct {
		RewardID uint `json:"reward_id"`
		Cost     int  `json:"cost"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.taskService.RedeemReward(userID.(uint), req.Cost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points or error occurred"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "redeemed"})
}

// Auth Handlers
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		RealName string `json:"real_name"`
		Grade    int    `json:"grade"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Role, req.RealName, req.Grade)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "message": "Registration successful"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *Handler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	h.authService.Logout(token)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Student List and Ranking
func (h *Handler) GetStudentList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user to find family ID
	user, err := h.taskService.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	students, err := h.taskService.GetStudentsByFamily(user.FamilyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *Handler) GetRanking(c *gin.Context) {
	limit := 10 // Top 10 by default
	students, err := h.taskService.GetTopStudents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ranking"})
		return
	}

	c.JSON(http.StatusOK, students)
}

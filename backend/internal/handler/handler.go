package handler

import (
	"net/http"
	"study-quest-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	taskService *service.TaskService
}

func NewHandler(ts *service.TaskService) *Handler {
	return &Handler{
		taskService: ts,
	}
}

func (h *Handler) GetAppConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"review_mode": false,
		"theme":       "default",
	})
}

func (h *Handler) GetTodayTasks(c *gin.Context) {
	tasks, _ := h.taskService.GetTodayTasks(1)
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
	h.taskService.CreateTask(req.Title, req.Points)
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func (h *Handler) SubmitTask(c *gin.Context) {
	var req struct {
		TaskID uint `json:"task_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	h.taskService.SubmitTask(1, req.TaskID)
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
	user, _ := h.taskService.GetUserProfile(1)
	c.JSON(http.StatusOK, user)
}


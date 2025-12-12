package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"study-quest-backend/internal/model"
	"study-quest-backend/internal/repository"
	"time"
)

type TaskService struct {
	taskRepo repository.ITaskRepository
	userRepo repository.IUserRepository
}

func NewTaskService(taskRepo repository.ITaskRepository, userRepo repository.IUserRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *TaskService) GetTodayTasks(studentID uint) ([]model.TaskLog, error) {
	return s.taskRepo.GetTodayTasks(studentID)
}

func (s *TaskService) GetPendingTasks() ([]model.TaskLog, error) {
	return s.taskRepo.GetPendingTasks()
}

func (s *TaskService) CreateTask(title string, points int) error {
	task := &model.Task{
		Title:  title,
		Points: points,
		Type:   1,
	}
	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) SubmitTask(studentID uint, taskID uint) error {
	return s.taskRepo.SubmitTask(studentID, taskID)
}

func (s *TaskService) ApproveTask(logID uint) error {
	// 1. Get task log to obtain student ID and points
	taskLog, err := s.taskRepo.GetTaskLog(logID)
	if err != nil {
		return err
	}

	// 2. Approve the task
	err = s.taskRepo.ApproveTask(logID)
	if err != nil {
		return err
	}

	// 3. Add points to student
	return s.userRepo.AddPoints(taskLog.StudentID, taskLog.Task.Points)
}

func (s *TaskService) RejectTask(logID uint) error {
	return s.taskRepo.RejectTask(logID)
}

func (s *TaskService) GetUserProfile(userID uint) (*model.User, error) {
	return s.userRepo.GetUser(userID)
}

func (s *TaskService) RedeemReward(studentID uint, rewardCost int) error {
	// 1. Check if user has enough points
	user, err := s.userRepo.GetUser(studentID)
	if err != nil {
		return err
	}

	if user.Points < rewardCost {
		return err
	}

	// 2. Deduct points
	return s.userRepo.AddPoints(studentID, -rewardCost)
}

func (s *TaskService) GetStudentsByFamily(familyID uint) ([]model.User, error) {
	return s.userRepo.GetStudentsByFamily(familyID)
}

func (s *TaskService) GetTopStudents(limit int) ([]model.User, error) {
	return s.userRepo.GetTopStudents(limit)
}

// AuthService
type AuthService struct {
	userRepo    repository.IUserRepository
	sessionRepo repository.ISessionRepository
}

func NewAuthService(userRepo repository.IUserRepository, sessionRepo repository.ISessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(username, password, role, realName string, grade int) (*model.User, error) {
	// Simple validation
	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	
	// Create user
	user := &model.User{
		Username: username,
		Password: password, // In production, should hash the password
		Role:     role,
		RealName: realName,
		Grade:    grade,
		FamilyID: 1, // Default family, in real app should be generated
	}
	
	if role == "student" {
		user.Points = 100 // Initial points for students
	}
	
	err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	// Get user
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}
	
	// Check password (simple comparison, should use bcrypt in production)
	if user.Password != password {
		return nil, "", errors.New("invalid username or password")
	}
	
	// Generate token
	token := generateToken()
	
	// Create session
	session := &model.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	
	err = s.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, "", err
	}
	
	return user, token, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.DeleteSession(token)
}

func (s *AuthService) ValidateSession(token string) (*model.User, error) {
	session, err := s.sessionRepo.GetSession(token)
	if err != nil {
		return nil, err
	}
	
	return s.userRepo.GetUser(session.UserID)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}


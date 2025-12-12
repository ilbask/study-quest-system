package service

import (
	"study-quest-backend/internal/model"
	"study-quest-backend/internal/repository"
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


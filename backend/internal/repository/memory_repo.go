package repository

import (
	"errors"
	"study-quest-backend/internal/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type ITaskRepository interface {
	GetTodayTasks(studentID uint) ([]model.TaskLog, error)
	GetPendingTasks() ([]model.TaskLog, error)
	GetTaskLog(logID uint) (*model.TaskLog, error)
	CreateTask(task *model.Task) error
	SubmitTask(studentID uint, taskID uint) error
	ApproveTask(logID uint) error
	RejectTask(logID uint) error
}

type IUserRepository interface {
	GetUser(id uint) (*model.User, error)
	AddPoints(userID uint, points int) error
}

// Memory Implementation
type MemoryTaskRepository struct {
	tasks    map[uint]*model.Task
	taskLogs map[uint]*model.TaskLog
	idCounter uint
	logCounter uint
	mu       sync.Mutex
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	repo := &MemoryTaskRepository{
		tasks:    make(map[uint]*model.Task),
		taskLogs: make(map[uint]*model.TaskLog),
		idCounter: 1,
		logCounter: 1,
	}
	// Seed Data
	repo.tasks[1] = &model.Task{Model: gorm.Model{ID: 1}, Title: "完成数学作业", Points: 30, Type: 1}
	repo.tasks[2] = &model.Task{Model: gorm.Model{ID: 2}, Title: "整理房间", Points: 20, Type: 2}
	
	// Assign tasks to student (log w/ status 0)
	repo.taskLogs[1] = &model.TaskLog{
		Model: gorm.Model{ID: 1}, StudentID: 1, TaskID: 1, Status: 0, 
		Task: *repo.tasks[1],
	}
	repo.taskLogs[2] = &model.TaskLog{
		Model: gorm.Model{ID: 2}, StudentID: 1, TaskID: 2, Status: 0,
		Task: *repo.tasks[2],
	}
	repo.logCounter = 3
	repo.idCounter = 3

	return repo
}

func (r *MemoryTaskRepository) GetTodayTasks(studentID uint) ([]model.TaskLog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var logs []model.TaskLog
	for _, log := range r.taskLogs {
		if log.StudentID == studentID {
			// Reload task info
			if t, ok := r.tasks[log.TaskID]; ok {
				log.Task = *t
			}
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (r *MemoryTaskRepository) GetPendingTasks() ([]model.TaskLog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var logs []model.TaskLog
	for _, log := range r.taskLogs {
		if log.Status == 1 { // Pending
			if t, ok := r.tasks[log.TaskID]; ok {
				log.Task = *t
			}
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (r *MemoryTaskRepository) GetTaskLog(logID uint) (*model.TaskLog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.taskLogs[logID]; ok {
		// Reload task info
		if t, ok := r.tasks[log.TaskID]; ok {
			log.Task = *t
		}
		return log, nil
	}
	return nil, errors.New("task log not found")
}

func (r *MemoryTaskRepository) CreateTask(task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task.ID = r.idCounter
	r.idCounter++
	r.tasks[task.ID] = task
	
	// Auto assign to student 1 for demo
	log := &model.TaskLog{
		Model: gorm.Model{ID: r.logCounter},
		StudentID: 1,
		TaskID: task.ID,
		Status: 0,
		Task: *task,
	}
	r.taskLogs[r.logCounter] = log
	r.logCounter++
	return nil
}

func (r *MemoryTaskRepository) SubmitTask(studentID uint, taskID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, log := range r.taskLogs {
		if log.StudentID == studentID && log.Model.ID == taskID && log.Status == 0 {
			log.Status = 1 // Pending
			log.SubmittedAt = time.Now()
			return nil
		}
	}
	return errors.New("task not found or not in todo state")
}

func (r *MemoryTaskRepository) ApproveTask(logID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.taskLogs[logID]; ok {
		log.Status = 2
		log.ApprovedAt = time.Now()
		return nil
	}
	return errors.New("log not found")
}

func (r *MemoryTaskRepository) RejectTask(logID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.taskLogs[logID]; ok {
		log.Status = 3
		return nil
	}
	return errors.New("log not found")
}

// Memory User Repo
type MemoryUserRepository struct {
	users map[uint]*model.User
	mu sync.Mutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: map[uint]*model.User{
			1: {Model: gorm.Model{ID: 1}, Username: "Student", Points: 100},
		},
	}
}

func (r *MemoryUserRepository) GetUser(id uint) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (r *MemoryUserRepository) AddPoints(userID uint, points int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.users[userID]; ok {
		u.Points += points
		return nil
	}
	return errors.New("user not found")
}


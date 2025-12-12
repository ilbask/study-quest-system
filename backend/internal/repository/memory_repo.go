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
	AssignTaskToStudent(studentID uint, taskID uint) error
	SubmitTask(studentID uint, taskID uint) error
	ApproveTask(logID uint) error
	RejectTask(logID uint) error
}

type IUserRepository interface {
	GetUser(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
	AddPoints(userID uint, points int) error
	GetStudentsByFamily(familyID uint) ([]model.User, error)
	GetTopStudents(limit int) ([]model.User, error)
}

type ISessionRepository interface {
	CreateSession(session *model.Session) error
	GetSession(token string) (*model.Session, error)
	DeleteSession(token string) error
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

func (r *MemoryTaskRepository) AssignTaskToStudent(studentID uint, taskID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	task, ok := r.tasks[taskID]
	if !ok {
		return errors.New("task not found")
	}
	
	log := &model.TaskLog{
		Model:     gorm.Model{ID: r.logCounter},
		StudentID: studentID,
		TaskID:    taskID,
		Status:    0, // Todo
		Task:      *task,
	}
	r.taskLogs[r.logCounter] = log
	r.logCounter++
	return nil
}

// Memory User Repo
type MemoryUserRepository struct {
	users map[uint]*model.User
	usersByUsername map[string]*model.User
	idCounter uint
	mu sync.Mutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	repo := &MemoryUserRepository{
		users: make(map[uint]*model.User),
		usersByUsername: make(map[string]*model.User),
		idCounter: 1,
	}
	
	// Seed data: Create demo users
	demoStudent := &model.User{
		Model: gorm.Model{ID: 1}, 
		Username: "student1", 
		Password: "123456",
		Role: "student",
		Points: 100,
		FamilyID: 1,
		RealName: "小明",
		Grade: 3,
	}
	demoParent := &model.User{
		Model: gorm.Model{ID: 2}, 
		Username: "parent1", 
		Password: "123456",
		Role: "parent",
		Points: 0,
		FamilyID: 1,
		RealName: "李妈妈",
	}
	
	repo.users[1] = demoStudent
	repo.users[2] = demoParent
	repo.usersByUsername["student1"] = demoStudent
	repo.usersByUsername["parent1"] = demoParent
	repo.idCounter = 3
	
	return repo
}

func (r *MemoryUserRepository) GetUser(id uint) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (r *MemoryUserRepository) GetUserByUsername(username string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.usersByUsername[username]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (r *MemoryUserRepository) CreateUser(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if username exists
	if _, exists := r.usersByUsername[user.Username]; exists {
		return errors.New("username already exists")
	}
	
	user.ID = r.idCounter
	r.idCounter++
	r.users[user.ID] = user
	r.usersByUsername[user.Username] = user
	return nil
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

func (r *MemoryUserRepository) GetStudentsByFamily(familyID uint) ([]model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var students []model.User
	for _, user := range r.users {
		if user.FamilyID == familyID && user.Role == "student" {
			students = append(students, *user)
		}
	}
	return students, nil
}

func (r *MemoryUserRepository) GetTopStudents(limit int) ([]model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	var students []model.User
	for _, user := range r.users {
		if user.Role == "student" {
			students = append(students, *user)
		}
	}
	
	// Sort by points descending
	for i := 0; i < len(students)-1; i++ {
		for j := i + 1; j < len(students); j++ {
			if students[i].Points < students[j].Points {
				students[i], students[j] = students[j], students[i]
			}
		}
	}
	
	if len(students) > limit {
		students = students[:limit]
	}
	
	return students, nil
}

// Memory Session Repo
type MemorySessionRepository struct {
	sessions map[string]*model.Session
	mu sync.Mutex
}

func NewMemorySessionRepository() *MemorySessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[string]*model.Session),
	}
}

func (r *MemorySessionRepository) CreateSession(session *model.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.Token] = session
	return nil
}

func (r *MemorySessionRepository) GetSession(token string) (*model.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[token]; ok {
		// Check if expired
		if time.Now().After(s.ExpiresAt) {
			delete(r.sessions, token)
			return nil, errors.New("session expired")
		}
		return s, nil
	}
	return nil, errors.New("session not found")
}

func (r *MemorySessionRepository) DeleteSession(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, token)
	return nil
}


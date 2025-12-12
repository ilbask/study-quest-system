package repository

import (
	"errors"
	"study-quest-backend/internal/model"
	"sync"
	"time"
)

// Interfaces
type ITaskRepository interface {
	GetTodayTasks(studentID uint) ([]model.TaskLog, error)
	GetPendingTasks() ([]model.TaskLog, error)
	GetTaskLog(logID uint) (*model.TaskLog, error)
	CreateTask(task *model.Task) error
	AssignTaskToStudent(studentID uint, taskID uint) error
	SubmitTask(studentID uint, taskID uint) error
	SubmitTaskByLogID(logID uint) error
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

type IRedemptionRepository interface {
	CreateRedemption(redemption *model.Redemption) error
	GetRedemptionsByFamily(familyID uint) ([]model.Redemption, error)
	GetRedemptionsByStudent(studentID uint) ([]model.Redemption, error)
}

type IRewardRepository interface {
	GetAllRewards() ([]model.Reward, error)
	GetReward(id uint) (*model.Reward, error)
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
	repo.tasks[1] = &model.Task{ID: 1, Title: "完成数学作业", Points: 30, Type: 1, CreatedAt: time.Now()}
	repo.tasks[2] = &model.Task{ID: 2, Title: "整理房间", Points: 20, Type: 2, CreatedAt: time.Now()}
	
	// Assign tasks to student (log w/ status 0)
	repo.taskLogs[1] = &model.TaskLog{
		ID: 1, StudentID: 1, TaskID: 1, Status: 0, 
		Task: *repo.tasks[1], CreatedAt: time.Now(),
	}
	repo.taskLogs[2] = &model.TaskLog{
		ID: 2, StudentID: 1, TaskID: 2, Status: 0,
		Task: *repo.tasks[2], CreatedAt: time.Now(),
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
		if log.StudentID == studentID && log.ID == taskID && log.Status == 0 {
			log.Status = 1 // Pending
			now := time.Now()
			log.SubmittedAt = &now
			return nil
		}
	}
	return errors.New("task not found or not in todo state")
}

func (r *MemoryTaskRepository) SubmitTaskByLogID(logID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.taskLogs[logID]; ok {
		if log.Status != 0 {
			return errors.New("task already submitted or completed")
		}
		log.Status = 1 // Pending
		now := time.Now()
		log.SubmittedAt = &now
		return nil
	}
	return errors.New("task log not found")
}

func (r *MemoryTaskRepository) ApproveTask(logID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.taskLogs[logID]; ok {
		log.Status = 2
		now := time.Now()
		log.ApprovedAt = &now
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
		ID:        r.logCounter,
		StudentID: studentID,
		TaskID:    taskID,
		Status:    0, // Todo
		Task:      *task,
		CreatedAt: time.Now(),
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
		ID: 1,
		Username: "student1", 
		Password: "123456",
		Role: "student",
		Points: 100,
		FamilyID: 1,
		RealName: "小明",
		Grade: 3,
		CreatedAt: time.Now(),
	}
	demoParent := &model.User{
		ID: 2,
		Username: "parent1", 
		Password: "123456",
		Role: "parent",
		Points: 0,
		FamilyID: 1,
		CreatedAt: time.Now(),
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

// MemoryRedemptionRepository
type MemoryRedemptionRepository struct {
	redemptions map[uint]*model.Redemption
	idCounter   uint
	mu          sync.Mutex
}

func NewMemoryRedemptionRepository() *MemoryRedemptionRepository {
	return &MemoryRedemptionRepository{
		redemptions: make(map[uint]*model.Redemption),
		idCounter:   1,
	}
}

func (r *MemoryRedemptionRepository) CreateRedemption(redemption *model.Redemption) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	redemption.ID = r.idCounter
	r.idCounter++
	redemption.CreatedAt = time.Now()
	r.redemptions[redemption.ID] = redemption
	return nil
}

func (r *MemoryRedemptionRepository) GetRedemptionsByFamily(familyID uint) ([]model.Redemption, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.Redemption
	for _, redemption := range r.redemptions {
		// We need to check if the student belongs to the family
		// For now, we'll just return all redemptions and filter in service layer
		result = append(result, *redemption)
	}
	return result, nil
}

func (r *MemoryRedemptionRepository) GetRedemptionsByStudent(studentID uint) ([]model.Redemption, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.Redemption
	for _, redemption := range r.redemptions {
		if redemption.StudentID == studentID {
			result = append(result, *redemption)
		}
	}
	return result, nil
}

// MemoryRewardRepository
type MemoryRewardRepository struct {
	rewards   map[uint]*model.Reward
	idCounter uint
	mu        sync.Mutex
}

func NewMemoryRewardRepository() *MemoryRewardRepository {
	repo := &MemoryRewardRepository{
		rewards:   make(map[uint]*model.Reward),
		idCounter: 1,
	}
	
	// Initialize with default rewards
	defaultRewards := []model.Reward{
		{Title: "看电视 30分钟", Cost: 50, Category: 1, Stock: 999},
		{Title: "玩手机 15分钟", Cost: 30, Category: 1, Stock: 999},
		{Title: "吃冰淇淋", Cost: 40, Category: 2, Stock: 10},
		{Title: "去游乐园", Cost: 200, Category: 2, Stock: 5},
	}
	
	for _, reward := range defaultRewards {
		reward.ID = repo.idCounter
		repo.idCounter++
		reward.CreatedAt = time.Now()
		repo.rewards[reward.ID] = &reward
	}
	
	return repo
}

func (r *MemoryRewardRepository) GetAllRewards() ([]model.Reward, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []model.Reward
	for _, reward := range r.rewards {
		result = append(result, *reward)
	}
	return result, nil
}

func (r *MemoryRewardRepository) GetReward(id uint) (*model.Reward, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reward, ok := r.rewards[id]; ok {
		return reward, nil
	}
	return nil, errors.New("reward not found")
}


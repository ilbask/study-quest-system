package repository

import (
	"study-quest-backend/internal/model"
	"gorm.io/gorm"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) GetUser(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *MySQLUserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *MySQLUserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *MySQLUserRepository) AddPoints(userID uint, points int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", points)).Error
}

func (r *MySQLUserRepository) GetStudentsByFamily(familyID uint) ([]model.User, error) {
	var students []model.User
	err := r.db.Where("family_id = ? AND role = ?", familyID, "student").Find(&students).Error
	return students, err
}

func (r *MySQLUserRepository) GetTopStudents(limit int) ([]model.User, error) {
	var students []model.User
	err := r.db.Where("role = ?", "student").
		Order("points DESC").
		Limit(limit).
		Find(&students).Error
	return students, err
}

type MySQLTaskRepository struct {
	db *gorm.DB
}

func NewMySQLTaskRepository(db *gorm.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

func (r *MySQLTaskRepository) GetTodayTasks(studentID uint) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	err := r.db.Preload("Task").Where("student_id = ?", studentID).Find(&logs).Error
	return logs, err
}

func (r *MySQLTaskRepository) GetPendingTasks() ([]model.TaskLog, error) {
	var logs []model.TaskLog
	err := r.db.Preload("Task").Where("status = ?", 1).Find(&logs).Error
	return logs, err
}

func (r *MySQLTaskRepository) GetTaskLog(logID uint) (*model.TaskLog, error) {
	var log model.TaskLog
	err := r.db.Preload("Task").First(&log, logID).Error
	return &log, err
}

func (r *MySQLTaskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *MySQLTaskRepository) AssignTaskToStudent(studentID uint, taskID uint) error {
	log := &model.TaskLog{
		StudentID: studentID,
		TaskID:    taskID,
		Status:    0, // Todo
	}
	return r.db.Create(log).Error
}

func (r *MySQLTaskRepository) SubmitTask(studentID uint, taskID uint) error {
	return r.db.Model(&model.TaskLog{}).
		Where("student_id = ? AND id = ? AND status = ?", studentID, taskID, 0).
		Updates(map[string]interface{}{
			"status":       1,
			"submitted_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *MySQLTaskRepository) ApproveTask(logID uint) error {
	return r.db.Model(&model.TaskLog{}).
		Where("id = ?", logID).
		Updates(map[string]interface{}{
			"status":      2,
			"approved_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *MySQLTaskRepository) RejectTask(logID uint) error {
	return r.db.Model(&model.TaskLog{}).
		Where("id = ?", logID).
		Update("status", 3).Error
}

type MySQLSessionRepository struct {
	db *gorm.DB
}

func NewMySQLSessionRepository(db *gorm.DB) *MySQLSessionRepository {
	return &MySQLSessionRepository{db: db}
}

func (r *MySQLSessionRepository) CreateSession(session *model.Session) error {
	return r.db.Create(session).Error
}

func (r *MySQLSessionRepository) GetSession(token string) (*model.Session, error) {
	var session model.Session
	err := r.db.Where("token = ? AND expires_at > NOW()", token).First(&session).Error
	return &session, err
}

func (r *MySQLSessionRepository) DeleteSession(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.Session{}).Error
}


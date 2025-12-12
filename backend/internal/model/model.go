package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Role     string `json:"role"` // 'parent', 'student'
	Points   int    `json:"points"`
	Avatar   string `json:"avatar"`
}

type Task struct {
	gorm.Model
	Title      string `json:"title"`
	Points     int    `json:"points"`
	Type       int    `json:"type"` // 1:Study, 2:Chore, 3:Habit
	Recurrence string `json:"recurrence"`
}

type TaskLog struct {
	gorm.Model
	StudentID   uint      `json:"student_id"`
	TaskID      uint      `json:"task_id"`
	Status      int       `json:"status"` // 0:InProgress, 1:Pending, 2:Done, 3:Rejected
	SubmittedAt time.Time `json:"submitted_at"`
	ApprovedAt  time.Time `json:"approved_at"`
	Task        Task      `json:"task"` // Preload
}

type Reward struct {
	gorm.Model
	Title    string `json:"title"`
	Cost     int    `json:"cost"`
	Category int    `json:"category"` // 1:Time, 2:Item
	Stock    int    `json:"stock"`
}

type Redemption struct {
	gorm.Model
	StudentID uint `json:"student_id"`
	RewardID  uint `json:"reward_id"`
	Cost      int  `json:"cost"`
}

type AppConfig struct {
	Key        string `gorm:"primaryKey" json:"key"`
	Value      string `json:"value"`
	Platform   string `json:"platform"`
	MinVersion string `json:"min_version"`
}


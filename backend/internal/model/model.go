package model

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Username  string     `json:"username"`
	Password  string     `json:"-"` // 不在 JSON 中返回
	Role      string     `json:"role"` // 'parent', 'student'
	Points    int        `json:"points"`
	Avatar    string     `json:"avatar"`
	FamilyID  uint       `json:"family_id"` // 家庭组ID
	Grade     int        `json:"grade"` // 年级（学生）
	RealName  string     `json:"real_name"` // 真实姓名
}

type Task struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Title      string     `json:"title"`
	Points     int        `json:"points"`
	Type       int        `json:"type"` // 1:Study, 2:Chore, 3:Habit
	Recurrence string     `json:"recurrence"`
}

type TaskLog struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	StudentID   uint       `json:"student_id"`
	TaskID      uint       `json:"task_id"`
	Status      int        `json:"status"` // 0:InProgress, 1:Pending, 2:Done, 3:Rejected
	SubmittedAt *time.Time `json:"submitted_at"`
	ApprovedAt  *time.Time `json:"approved_at"`
	Task        Task       `json:"task" gorm:"foreignKey:TaskID"` // Preload
}

type Reward struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Title    string `json:"title"`
	Cost     int    `json:"cost"`
	Category int    `json:"category"` // 1:Time, 2:Item
	Stock    int    `json:"stock"`
}

type Redemption struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	StudentID   uint       `json:"student_id"`
	RewardID    uint       `json:"reward_id"`
	RewardTitle string     `json:"reward_title"` // Store title for display
	Cost        int        `json:"cost"`
	Student     User       `json:"student" gorm:"foreignKey:StudentID"`
	Reward      Reward     `json:"reward" gorm:"foreignKey:RewardID"`
}

type AppConfig struct {
	Key        string `gorm:"primaryKey" json:"key"`
	Value      string `json:"value"`
	Platform   string `json:"platform"`
	MinVersion string `json:"min_version"`
}

type Session struct {
	Token     string `gorm:"primaryKey" json:"token"`
	UserID    uint   `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}


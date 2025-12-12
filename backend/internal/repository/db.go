package repository

import (
	"fmt"
	"log"
	"study-quest-backend/internal/config"
	"study-quest-backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Task{},
		&model.TaskLog{},
		&model.Reward{},
		&model.Redemption{},
		&model.AppConfig{},
		&model.Session{},
	)
}

func SeedData(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		log.Println("Data already seeded, skipping...")
		return nil
	}

	log.Println("Seeding initial data...")

	// Create demo users
	users := []model.User{
		{
			Username: "student1",
			Password: "123456",
			Role:     "student",
			Points:   100,
			FamilyID: 1,
			RealName: "小明",
			Grade:    3,
		},
		{
			Username: "parent1",
			Password: "123456",
			Role:     "parent",
			Points:   0,
			FamilyID: 1,
			RealName: "李妈妈",
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Username, err)
		}
	}

	// Create demo tasks
	tasks := []model.Task{
		{Title: "完成数学作业", Points: 30, Type: 1},
		{Title: "整理房间", Points: 20, Type: 2},
	}

	for _, task := range tasks {
		if err := db.Create(&task).Error; err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		// Assign to student1
		log := model.TaskLog{
			StudentID: 1,
			TaskID:    task.ID,
			Status:    0,
		}
		if err := db.Create(&log).Error; err != nil {
			return fmt.Errorf("failed to create task log: %w", err)
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}


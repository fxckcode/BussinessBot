package models

import (
	"github.com/fxckcode/BussinessBot/db"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string         `gorm:"not null;unique_index" json:"title"`
	Description string         `json:"description"`
	Completed   bool           `gorm:"not null;default:false" json:"completed"`
}

func (*Task) FindAllTasks() (*[]Task, error) {
	var tasks []Task
	result := db.DB.Find(&tasks)

	if result.Error != nil {
		return &[]Task{}, result.Error
	}

	return &tasks, nil
}

func (*Task) FindTaskById(id string) (*Task, error) {
	var task Task
	result := db.DB.First(&task, id)

	if result.Error != nil {
		return &Task{}, result.Error
	}

	return &task, nil
}

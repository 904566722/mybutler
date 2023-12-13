package entity

import (
	"database/sql"

	"gorm.io/gorm"

	"changeme/backend/pkg/model/common"
)

// Task 任务模型
type Task struct {
	common.BaseModel
	Name     string          `json:"name" gorm:"size:255"`
	Deadline sql.NullTime    `json:"deadline"`
	Priority common.Priority `json:"priority"`
}

type SubTask struct {
	common.BaseModel
	Name     string             `json:"name"`
	DoneAt   sql.NullTime       `json:"doneAt"`
	Deadline sql.NullTime       `json:"deadline"`
	Type     common.SubTaskType `json:"type"`
	Priority common.Priority    `json:"priority"`

	TaskID string `json:"taskID" gorm:"size:32"` // sub task belongs to task
	Task   Task   `json:"task"`
}

// TaskCheckIn 任务打卡模型
type TaskCheckIn struct {
	common.BaseModel
	TaskID string `json:"taskID" gorm:"size:32"` // sub task belongs to task
}

func (t *Task) GetPrefix() string {
	return "t"
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	return t.BaseModel.BeforeCreate(t.GetPrefix(), tx)
}

func (subt *SubTask) GetPrefix() string {
	return "subt"
}

func (subt *SubTask) BeforeCreate(tx *gorm.DB) (err error) {
	return subt.BaseModel.BeforeCreate(subt.GetPrefix(), tx)
}

func (tci *TaskCheckIn) GetPrefix() string {
	return "tci"
}

func (tci *TaskCheckIn) BeforeCreate(tx *gorm.DB) (err error) {
	return tci.BaseModel.BeforeCreate(tci.GetPrefix(), tx)
}

const (
	ColTaskDeadline = "deadline"

	ColSubTaskTaskId = "task_id"
)

package service

import (
	"context"
	"database/sql"
	"time"

	"changeme/backend/pkg/db"
	"changeme/backend/pkg/log"
	"changeme/backend/pkg/model/common"
	"changeme/backend/pkg/model/entity"
)

type TaskService struct {
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (t *TaskService) Create(ctx context.Context, task *entity.Task) error {
	result := db.Default().Create(task)
	if result.Error != nil {
		return result.Error
	}
	log.Info("create task success", log.String("id", task.ID))
	return nil
}

func (t *TaskService) Get(ctx context.Context, id string) (*entity.Task, error) {
	task := &entity.Task{}
	result := db.Default().Where(common.ColId+" = ?", id).First(task)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}

func (t *TaskService) ModifyDeadline(ctx context.Context, id string, time time.Time) error {
	result := db.Default().Model(&entity.Task{}).Where(common.ColId+" = ?", id).Update(entity.ColTaskDeadline, time)
	if result.Error != nil {
		return result.Error
	}
	log.Info("modify task deadline success", log.String("id", id), log.String("time", time.String()))
	return nil
}

// CalcLeftDays 计算剩余天数
func (t *TaskService) CalcLeftDays(ctx context.Context, deadline sql.NullTime) int {
	if !deadline.Valid {
		return -1
	}
	now := time.Now()
	return int(deadline.Time.Sub(now).Hours() / 24)
}

// CalcProgress 计算进度
func (t *TaskService) CalcProgress(ctx context.Context, subTasks []*entity.SubTask) uint8 {
	if len(subTasks) == 0 {
		return 0
	}
	var doneCount int
	for _, subTask := range subTasks {
		if subTask.DoneAt.Valid {
			doneCount++
		}
	}
	return uint8(doneCount * 100 / len(subTasks))
}

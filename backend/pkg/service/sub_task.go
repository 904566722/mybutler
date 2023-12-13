package service

import (
	"context"

	"changeme/backend/pkg/db"
	"changeme/backend/pkg/log"
	"changeme/backend/pkg/model/common"
	"changeme/backend/pkg/model/entity"
)

type SubTaskService struct {
}

func NewSubTaskService() *SubTaskService {
	return &SubTaskService{}
}

func (t *SubTaskService) Get(ctx context.Context, id string) (*entity.SubTask, error) {
	subTask := &entity.SubTask{}
	result := db.Default().Where(common.ColId+" = ?", id).First(subTask)
	if result.Error != nil {
		return nil, result.Error
	}
	return subTask, nil
}

func (t *SubTaskService) Create(ctx context.Context, subTask *entity.SubTask) error {
	result := db.Default().Create(subTask)
	if result.Error != nil {
		return result.Error
	}
	log.Info("create subTask success", log.String("id", subTask.ID))
	return nil
}

func (t *SubTaskService) ListByConditions(ctx context.Context, conds ListConditions) ([]*entity.SubTask, error) {
	var subTasks []*entity.SubTask
	result := db.Default().Where(conds.WhereExpr()).Find(&subTasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return subTasks, nil
}

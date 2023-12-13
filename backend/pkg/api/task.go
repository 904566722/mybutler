package api

import (
	"context"

	"changeme/backend/pkg/api/converter"
	"changeme/backend/pkg/model/entity"
	"changeme/backend/pkg/model/vo"
	"changeme/backend/pkg/service"
)

type TaskApi struct {
	taskSvc    *service.TaskService
	subTaskSvc *service.SubTaskService
	taskCvt    *converter.TaskConverter
	subTaskCvt *converter.SubTaskConverter
}

func NewTaskApi() *TaskApi {
	return &TaskApi{
		taskSvc:    service.NewTaskService(),
		subTaskSvc: service.NewSubTaskService(),
		taskCvt:    converter.NewTaskConverter(),
		subTaskCvt: converter.NewSubTaskConverter(),
	}
}

func (t *TaskApi) Get(ctx context.Context, id string) (*vo.Task, error) {
	task, err := t.taskSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	subTasks, err := t.subTaskSvc.ListByConditions(ctx, []*service.ListCondition{
		{entity.ColSubTaskTaskId, "=", task.ID},
	})
	if err != nil {
		return nil, err
	}

	taskVo := t.taskCvt.Model2VO(task)
	subTaskVos := t.subTaskCvt.BatchModel2VO(subTasks)
	taskVo.SubTasks = subTaskVos
	// 计算剩余天数、进度
	if task.Deadline.Valid {
		taskVo.Deadline = &task.Deadline.Time
	}
	taskVo.LeftDays = t.taskSvc.CalcLeftDays(ctx, task.Deadline)
	taskVo.Progress = t.taskSvc.CalcProgress(ctx, subTasks)
	if task.Deadline.Valid {
		taskVo.Deadline = &task.Deadline.Time
		taskVo.Done = true
	}
	return taskVo, nil
}

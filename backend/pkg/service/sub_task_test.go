package service

import (
	"context"
	"testing"

	"changeme/backend/pkg/model/common"
	"changeme/backend/pkg/model/entity"
)

func TestSubTaskService_Create(t1 *testing.T) {
	DependInit()
	if err := NewSubTaskService().Create(context.Background(), &entity.SubTask{
		Name:     "编码",
		Type:     common.SubTaskTypeCoding,
		Priority: common.PriorityHigh,
		TaskID:   "t-NPosxgVn",
	}); err != nil {
		t1.Error(err)
		return
	}
}

func TestSubTaskService_ListByConditions(t1 *testing.T) {
	DependInit()
	subTasks, err := NewSubTaskService().ListByConditions(context.Background(), []*ListCondition{
		{
			Field: entity.ColSubTaskTaskId,
			Op:    "=",
			Value: "t-NPosxgVn",
		},
	})
	if err != nil {
		t1.Error(err)
		return
	}
	t1.Log(subTasks)
}

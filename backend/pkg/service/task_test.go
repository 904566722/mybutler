package service

import (
	"context"
	"testing"
	"time"

	"changeme/backend/pkg/configs"
	"changeme/backend/pkg/db"
	"changeme/backend/pkg/log"
	"changeme/backend/pkg/model/common"
	"changeme/backend/pkg/model/entity"
)

func DependInit() {
	if err := configs.InitConfig(); err != nil {
		panic(err)
	}
	log.InitLog()
	db.InitMysql()
}

func TestTaskService_Create(t *testing.T) {
	DependInit()
	ctx := context.Background()
	if err := NewTaskService().Create(ctx, &entity.Task{
		Name:     "安装向导",
		Priority: common.PriorityHigh,
	}); err != nil {
		t.Error(err)
		return
	}
}

func TestTaskService_ModifyDeadline(t1 *testing.T) {
	DependInit()
	if err := NewTaskService().ModifyDeadline(context.Background(), "t-NPosxgVn", time.Now()); err != nil {
		t1.Error(err)
		return
	}
}

func TestTaskService_Get(t1 *testing.T) {
	DependInit()
	taskInf, err := NewTaskService().Get(context.Background(), "t-NPosxgVn")
	if err != nil {
		t1.Error(err)
		return
	}
	t1.Log(taskInf)
}

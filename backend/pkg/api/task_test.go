package api

import (
	"context"
	"testing"

	"changeme/backend/pkg/configs"
	"changeme/backend/pkg/db"
	"changeme/backend/pkg/log"
)

func DependInit() {
	if err := configs.InitConfig(); err != nil {
		panic(err)
	}
	log.InitLog()
	db.InitMysql()
}

func TestTaskApi_Get(t1 *testing.T) {
	DependInit()
	task, err := NewTaskApi().Get(context.Background(), "t-NPosxgVn")
	if err != nil {
		t1.Error(err)
		return
	}
	t1.Log(task)
}

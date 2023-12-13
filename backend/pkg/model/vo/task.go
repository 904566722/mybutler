package vo

import (
	"time"

	"changeme/backend/pkg/model/common"
)

type BaseVO struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
}

type Task struct {
	BaseVO
	Name     string          `json:"name"`
	Done     bool            `json:"done"`     // 是否完成
	DoneAt   *time.Time      `json:"doneAt"`   // 完成时间
	Deadline *time.Time      `json:"deadline"` // 截止日期
	LeftDays int             `json:"leftDays"` // 剩余天数
	Progress uint8           `json:"progress"` // 进度
	Priority common.Priority `json:"priority"` // 优先级
	SubTasks []*SubTask      `json:"subTasks"` // 子任务
	CheckIns []time.Time     `json:"checkIns"` // 打卡记录
}

type SubTask struct {
	BaseVO
	Name     string             `json:"name"`
	Done     bool               `json:"done"`
	DoneAt   *time.Time         `json:"doneAt"`
	Deadline *time.Time         `json:"deadline"`
	Type     common.SubTaskType `json:"type"`
	Priority common.Priority    `json:"priority"`
}

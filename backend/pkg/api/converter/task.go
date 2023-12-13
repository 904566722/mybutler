package converter

import (
	"changeme/backend/pkg/model/entity"
	"changeme/backend/pkg/model/vo"
)

type TaskConverter struct{}

func NewTaskConverter() *TaskConverter {
	return &TaskConverter{}
}

func (t TaskConverter) Model2VO(m *entity.Task) *vo.Task {
	if m == nil {
		return &vo.Task{}
	}

	taskVo := &vo.Task{
		BaseVO: vo.BaseVO{
			ID:      m.ID,
			Comment: m.Comment,
		},
		Name: m.Name,
		//Done:     false,
		//DoneAt:   nil,
		//Deadline: nil,
		//LeftDays: 0,
		//Progress: 0,
		Priority: m.Priority,
		//SubTasks: nil,
		//CheckIns: nil,
	}

	if m.Deadline.Valid {
		taskVo.Deadline = &m.Deadline.Time
	}

	return taskVo
}

type SubTaskConverter struct{}

func NewSubTaskConverter() *SubTaskConverter {
	return &SubTaskConverter{}
}

func (t SubTaskConverter) Model2VO(m *entity.SubTask) *vo.SubTask {
	if m == nil {
		return &vo.SubTask{}
	}

	subTaskVo := &vo.SubTask{
		BaseVO: vo.BaseVO{
			ID:      m.ID,
			Comment: m.Comment,
		},
		Name: m.Name,
		//Done:     false,
		//DoneAt:   time.Time{},
		//Deadline: time.Time{},
		Type:     m.Type,
		Priority: m.Priority,
	}

	if m.DoneAt.Valid {
		subTaskVo.DoneAt = &m.DoneAt.Time
		subTaskVo.Done = true
	}
	if m.Deadline.Valid {
		subTaskVo.Deadline = &m.Deadline.Time
	}
	return subTaskVo
}

func (t SubTaskConverter) BatchModel2VO(mList []*entity.SubTask) []*vo.SubTask {
	if mList == nil {
		return nil
	}
	var voList []*vo.SubTask
	for _, task := range mList {
		voList = append(voList, t.Model2VO(task))
	}
	return voList
}

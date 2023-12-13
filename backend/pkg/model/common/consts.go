package common

type SubTaskType uint8

const (
	SubTaskTypeCoding SubTaskType = iota // 编码任务
	SubTaskTypeTest                      // 测试任务
	SubTaskTypeFix                       // 修复任务
	SubTaskTypePerf                      // 优化任务
	SubTaskTypeReview                    // 评审任务
	SubTaskTypeDocs                      // 文档任务
	SubTaskTypeDeploy                    // 部署任务
)

type Priority uint8

const (
	PriorityLow      Priority = iota // 低
	PriorityMedium                   // 中
	PriorityHigh                     // 高
	PriorityUrgent                   // 紧急
	PriorityCritical                 // 至关重要
)

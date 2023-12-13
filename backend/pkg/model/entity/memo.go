package entity

import (
	"changeme/backend/pkg/model/common"
)

// Memo 备忘录模型
type Memo struct {
	common.BaseModel
	Content string `json:"content"`
}

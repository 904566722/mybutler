package common

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"changeme/backend/pkg/utils"
)

type BaseModel struct {
	ID        string `gorm:"primarykey,size:32"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`

	Comment string `json:"comment"`
}

func (bm *BaseModel) BeforeCreate(idPrefix string, tx *gorm.DB) (err error) {
	if bm.ID == "" {
		bm.ID = utils.GenId(idPrefix)
	}
	if bm.CreatedAt.IsZero() {
		bm.CreatedAt = time.Now()
	}
	return nil
}

const (
	ColId      = "id"
	ColComment = "comment"
)

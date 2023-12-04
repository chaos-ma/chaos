package model

/**
* created by mengqi on 2023/12/4
 */
import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

type BaseSoftDelModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

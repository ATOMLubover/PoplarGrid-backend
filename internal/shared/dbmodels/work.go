package dbmodels

import (
	"gorm.io/gorm"
)

// 作品集（同步尨译信息）
// 考虑到数量较少，前端可以直接加载完整列表，这里不提供除了
// 主键和外键以外的索引
type Workset struct {
	gorm.Model

	Title     string
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`
}

func (Workset) TableName() string {
	return "worksets"
}

// 作品（同步尨译信息）
// 作品数量较大，所以对 title 也启用外键
type Work struct {
	gorm.Model

	Title     string `gorm:"uniqueIndex;type:text;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`

	Description string `gorm:"type:text"`
}

func (Work) TableName() string {
	return "works"
}

// 汉化组
// 数量较少，可以一次加载完全，所以同样没有对 name 使用索引
type Team struct {
	gorm.Model

	Name      string `gorm:"unique;size:256;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`
}

func (Team) TableName() string {
	return "teams"
}

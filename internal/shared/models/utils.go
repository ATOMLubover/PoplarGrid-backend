package models

import "time"

// ---------------------
// utils 指定一些辅助的类型以及函数
// ---------------------

// PKey 定义了一个通用的主键类型
type PKey uint64

// BaseModel 定义了一个通用的基础模型，包含了常用的字段
type BaseModel struct {
	// Id 是主键
	Id PKey `gorm:"primaryKey;autoIncrement"`
	// CreatedAt 是创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	// UpdatedAt 是更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
	// DeletedAt 是删除时间，软删除用
	// 默认是 null，表示未删除
	DeletedAt *time.Time `gorm:"index"`
}

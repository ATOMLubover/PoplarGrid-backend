package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// 为 GORM 自带的 uint 主键类型起一个醒目的别名
type PrimaryKey uint

// 代替 GORM 的 Model
type BaseModel struct {
	// 主键 ID
	Id PrimaryKey `gorm:"primaryKey;autoIncrement;not null"`
	// 创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	// 更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
	// 删除时间
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 职责掩码位移
const (
	LABOR_CREATOR_SHIFT      = iota // 负责人（仅在项目中有效）
	LABOR_PRINCIPAL_SHIFT           // 监制
	LABOR_SRC_PROV_SHIFT            // 图源
	LABOR_CLEANER_SHIFT             // 美工
	LABOR_GRAPHIC_PROC_SHIFT        // 修图
	LABOR_TRANSLATOR_SHIFT          // 翻译
	LABOR_PROOF_SHIFT               // 校对
	LABOR_LETTERER_SHIFT            // 嵌字
	LABOR_REVIEWER_SHIFT            // 审核
	LABOR_PUBLISHER_SHIFT           // 发布
)

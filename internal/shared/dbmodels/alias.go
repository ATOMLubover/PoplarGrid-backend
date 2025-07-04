package dbmodels

import (
	"time"
	"unsafe"

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
	LABOR_CREATOR_SHIFT      = iota // 项目创建者
	LABOR_PRINCIPAL_SHIFT           // 监制
	LABOR_SRC_PROV_SHIFT            // 图源提供者
	LABOR_CLEANER_SHIFT             // 美工
	LABOR_GRAPHIC_PROC_SHIFT        // 修图
	LABOR_TRANSLATOR_SHIFT          // 翻译
	LABOR_PROOF_SHIFT               // 校对
	LABOR_LETTERER_SHIFT            // 嵌字
	LABOR_REVIEWER_SHIFT            // 审核
	LABOR_PUBLISHER_SHIFT           // 发布
)

// 项目状态掩码
// 对应位是1代表正在处于该状态
const (
	// 当前项目 status 如果为 0 是未定义的
	PROJ_STATUS_UNSET_MASK = 1 << iota
	PROJ_STATUS_ON_TRANSLATING_MASK
	PROJ_STATUS_TRANSLATED_MASK
	PROJ_STATUS_ON_PROOF_MASK
	PROJ_STATUS_PROVED_MASK
	PROJ_STATUS_ON_LETTERING_MASK
	PROJ_STATUS_LETTERED_MASK
	PROJ_STATUS_ON_REVIEWING_MASK
	PROJ_STATUS_REVIEWED_MASK
	PROJ_STATUS_PUBLISHED_MASK

	// 这里所取的是最高位，可以直接检查
	PROJ_STATUS_CANCELED = 1 << (unsafe.Sizeof(Project{}.Status)*8 - 1)
)

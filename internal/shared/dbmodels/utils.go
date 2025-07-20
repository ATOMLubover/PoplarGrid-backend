package dbmodels

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

// 职责掩码位移常量
const (
	LABOR_DIRECTOR_SHIFT     = iota // 负责人/组内管理员 - 位移 0
	LABOR_PRINCIPAL_SHIFT           // 监制 - 位移 1
	LABOR_SRC_PROV_SHIFT            // 图源 - 位移 2
	LABOR_CLEANER_SHIFT             // 美工 - 位移 3
	LABOR_GRAPHIC_PROC_SHIFT        // 修图 - 位移 4
	LABOR_TRANSLATOR_SHIFT          // 翻译 - 位移 5
	LABOR_PROOF_SHIFT               // 校对 - 位移 6
	LABOR_LETTERER_SHIFT            // 嵌字 - 位移 7
	LABOR_REVIEWER_SHIFT            // 审核 - 位移 8
	LABOR_PUBLISHER_SHIFT           // 发布 - 位移 9
)

// LaborMask 表示成员的权限
type LaborMask uint

// 各职责对应的掩码常量 (通过左移位移量生成)
const (
	LABOR_CREATOR_MASK      LaborMask = 1 << LABOR_DIRECTOR_SHIFT
	LABOR_PRINCIPAL_MASK    LaborMask = 1 << LABOR_PRINCIPAL_SHIFT
	LABOR_SRC_PROV_MASK     LaborMask = 1 << LABOR_SRC_PROV_SHIFT
	LABOR_CLEANER_MASK      LaborMask = 1 << LABOR_CLEANER_SHIFT
	LABOR_GRAPHIC_PROC_MASK LaborMask = 1 << LABOR_GRAPHIC_PROC_SHIFT
	LABOR_TRANSLATOR_MASK   LaborMask = 1 << LABOR_TRANSLATOR_SHIFT
	LABOR_PROOF_MASK        LaborMask = 1 << LABOR_PROOF_SHIFT
	LABOR_LETTERER_MASK     LaborMask = 1 << LABOR_LETTERER_SHIFT
	LABOR_REVIEWER_MASK     LaborMask = 1 << LABOR_REVIEWER_SHIFT
	LABOR_PUBLISHER_MASK    LaborMask = 1 << LABOR_PUBLISHER_SHIFT
)

// HasRole 检查 LaborMask 是否包含某个职责
func (m LaborMask) HasRole(roleMask LaborMask) bool {
	return (m & roleMask) != 0
}

// AddRole 为 LaborMask 添加一个职责
func (m *LaborMask) AddRole(roleMask LaborMask) {
	*m |= roleMask
}

// RemoveRole 从 LaborMask 中移除一个职责
func (m *LaborMask) RemoveRole(roleMask LaborMask) {
	*m &= ^roleMask
}

// NewLaborMask 根据多个职责掩码创建一个新的复合掩码
func NewLaborMask(roles ...LaborMask) LaborMask {
	var mask LaborMask

	for _, role := range roles {
		mask.AddRole(role)
	}

	return mask
}

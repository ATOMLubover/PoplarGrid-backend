package dbmodels

import (
	"time"

	"gorm.io/gorm"
)

// 成员的基本信息
type Member struct {
	gorm.Model

	// 基本信息
	Nickname     string `gorm:"unique;size:128;not null"`
	Email        string `gorm:"unique;size:128;not null"`
	PasswordHash string `gorm:"size:256;not null"`

	// 接受尨译分配的 ID
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`

	// 在仪表盘中的身份
	PoplarIsAdmin bool
	// 负责角色
	Labors uint16

	// 补充备注
	Remark string `gorm:"type:text"`

	// 上一次活跃时间（可能是通过 ping 来确定）
	LastActive time.Time
}

func (Member) TableName() string {
	return "members"
}

// tag 基础信息
type Tag struct {
	gorm.Model

	Name        string `gorm:"uniqueIndex;size:128;not null"`
	Description string `gorm:"size:256"`
}

func (Tag) TableName() string {
	return "tags"
}

// 成员偏好 tag
// 关联成员与 tag
type MemberPreference struct {
	ID uint `gorm:"primaryKey"`

	MemberId PrimaryKey `gorm:"index"`
	FkMember *Member    `gorm:"foreignKey:MemberId"`

	TagId PrimaryKey `gorm:"index"`
	FkTag *Tag       `gorm:"foreignKey:TagId"`

	// 标志是否是被偏好 tag
	IsResisted bool `gorm:"default:false;index"`
}

func (MemberPreference) TableName() string {
	return "member_preferences"
}

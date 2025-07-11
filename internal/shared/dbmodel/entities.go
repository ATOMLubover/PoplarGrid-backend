package dbmodel

import (
	"time"
)

// 成员的基本信息
type Member struct {
	BaseModel

	// 所属的汉化组
	TeamId PrimaryKey
	FkTeam Team `gorm:"foreignKey:TeamId"`

	// 基本信息
	Nickname     string `gorm:"uniqueIndex;size:128;not null"`
	Email        string `gorm:"unique;size:128;not null"`
	PasswordHash string `gorm:"size:256;not null"`

	// 接受尨译分配的 ID
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`

	// 在仪表盘中的身份
	PoplarIsAdmin bool
	// 负责角色
	Labors uint

	// 补充备注
	Remark string `gorm:"type:text"`
	// QQ 号
	QqNumber string `gorm:"size:64"`

	// 上一次活跃时间（可能是通过 ping 来确定）
	LastActive time.Time
}

func (Member) TableName() string {
	return "members"
}

// tag 基础信息
// 将 tag 设计为各个汉化组共用的
type Tag struct {
	BaseModel

	Name        string `gorm:"uniqueIndex;size:128;not null"`
	Description string `gorm:"size:256"`
}

func (Tag) TableName() string {
	return "tags"
}

// 作品集（同步尨译信息）
// 考虑到数量较少，前端可以直接加载完整列表，这里不提供除了
// 主键和外键以外的索引
type Workset struct {
	BaseModel

	TeamId PrimaryKey
	FkTeam Team `gorm:"foreignKey:TeamId"`

	Name      string `gorm:"unique;type:text;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`
}

func (Workset) TableName() string {
	return "worksets"
}

// 汉化组
// 数量较少，可以一次加载完全，所以同样没有对 name 使用索引
type Team struct {
	BaseModel

	Name      string `gorm:"unique;size:256;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`
}

func (Team) TableName() string {
	return "teams"
}

// 项目进度表
type Project struct {
	BaseModel

	// 内嵌作品的信息
	Title     string `gorm:"index;type:text;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text;not null"`

	// 历史遗留序号（【】中的序号），保留对老作品的兼容
	// 经过观察，有序号重复的地方，如果可以最好重构这部分
	LegacyId uint `gorm:"index"`

	// 所属作品集
	WorksetId PrimaryKey
	FkWorkset Workset `gorm:"foreignKey:WorksetId"`

	// 当前项目的状态
	OnTranslating bool `gorm:"not null;default:false"`
	IsTranslated  bool `gorm:"not null;default:false"`
	OnProoving    bool `gorm:"not null;default:false"`
	IsProoved     bool `gorm:"not null;default:false"`
	OnLettering   bool `gorm:"not null;default:false"`
	IsLettered    bool `gorm:"not null;default:false"`
	OnReviewing   bool `gorm:"not null;default:false"`
	IsReviewed    bool `gorm:"not null;default:false"`
	IsPublished   bool `gorm:"not null;default:false"`
}

func (Project) TableName() string {
	return "projects"
}

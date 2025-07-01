package dbmodels

import (
	"time"
	"unsafe"

	"gorm.io/gorm"
)

// 为 GORM 自带的 uint 主键类型起一个醒目的别名
type PrimaryKey uint

// 职责掩码位移
const (
	LABOR_CREATOR_SHIFT = iota
	LABOR_PRICINPAL_SHIFT
	LABOR_SRC_PROV_SHIFT
	LABOR_CLEANER_SHIFT
	LABOR_GRAPHIC_PROC_SHIFT
	LABOR_TRANSLATOR_SHIFT
	LABOR_PROOF_SHIFT
	LABOR_LETTERER_SHIFT
	LABOR_REVIEWER_SHIFT
)

// 成员的基本信息
type Member struct {
	gorm.Model

	// 基本信息
	Nickname     string `gorm:"unique;size:128;not null"`
	Email        string `gorm:"unique;size:128;not null"`
	PasswordHash string `gorm:"size:256;not null"`

	// 接受尨译分配的 ID
	LongyiId string

	// 在仪表盘中的身份
	IsAdmin bool
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
type MemberPreference struct {
	ID uint `gorm:"primaryKey"`

	MemberId uint
	FkMember Member `gorm:"foreignKey:MemberId"`

	TagId uint
	FkTag Tag `gorm:"foreignKey:TagId"`

	// 标志是否是被偏好 tag
	IsPrefered bool
}

func (MemberPreference) TableName() string {
	return "member_preferences"
}

// 作品集
type Workset struct {
	gorm.Model

	Title string
}

func (Workset) TableName() string {
	return "worksets"
}

// 项目状态掩码
// 对应位是1代表正在处于该状态
const (
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

	PROJ_STATUS_CANCELED = 1 << (unsafe.Sizeof(Project{}.Status)*8 - 1)
)

// 项目进度表
type Project struct {
	gorm.Model

	// 该项目所属于的汉化组（同步尨译）
	TeamAffiliatedTo string
	// 历史遗留序号（【】中的序号）
	LegacyId uint `gorm:"index"`

	// 所属作品集
	WorksetId uint
	FkWorkset Workset `gorm:"foreignKey:WorksetId"`

	Title       string
	Description string

	// 当前项目的状态（partition依据）
	Status uint
	// 用紧急度代替好漫无汉等
	Urgency int16 `gorm:"type:smallint"`
}

func (Project) TableName() string {
	return "projects"
}

// 项目所携带 tag
type ProjectTag struct {
	ProjectId PrimaryKey `gorm:"index"`
	TagId     PrimaryKey `gorm:"index"`

	FkProject Project `gorm:"foreignKey:ProjectId"`
	FkTag     Tag     `gorm:"foreignKey:TagId"`
}

func (ProjectTag) TableName() string {
	return "project_tags"
}

// 项目分工表
type ProjectLaborDivision struct {
	gorm.Model

	ProjectId PrimaryKey
	FkProject Project `gorm:"foreignKey:ProjectId"`

	MemberId PrimaryKey
	FkMember Member `gorm:"foreignKey:MemberId"`

	// 分工，使用掩码计算多重身份
	LaborRole uint
}

func (ProjectLaborDivision) TableName() string {
	return "project_labor_divisions"
}

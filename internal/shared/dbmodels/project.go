package dbmodels

import "gorm.io/gorm"

// 项目进度表
type Project struct {
	gorm.Model

	// 历史遗留序号（【】中的序号），保留对老作品的兼容
	// 经过观察，有序号重复的地方，如果可以最好重构这部分
	LegacyId uint `gorm:"index"`

	// 所属汉化组（同步尨译）
	TeamId PrimaryKey
	FkTeam Team `gorm:"foreignKey:TeamId"`

	// 所属作品集（同步尨译）
	WorksetId PrimaryKey
	FkWorkset Workset `gorm:"foreignKey:WorksetId"`

	// 所属作品（同步龙译）
	WorkId PrimaryKey
	FkWork Work `gorm:"foreignKey:WorkId"`

	// 当前项目的状态
	Status uint
	// 用紧急度代替好漫无汉等
	Urgency int16 `gorm:"type:smallint"`
}

func (Project) TableName() string {
	return "projects"
}

// 项目所携带 tag
// 关联项目与 tag
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

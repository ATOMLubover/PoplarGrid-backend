package dbmodels

// 成员偏好 tag
// 关联成员与 tag
type MemberPreference struct {
	MemberId PrimaryKey `gorm:"index"`
	FkMember Member     `gorm:"foreignKey:MemberId"`

	TagId PrimaryKey `gorm:"index"`
	FkTag Tag        `gorm:"foreignKey:TagId"`

	// 标志是否是被偏好 tag
	IsResisted bool `gorm:"default:false;index"`
}

func (MemberPreference) TableName() string {
	return "member_preferences"
}

// 作品所携带 tag
// 关联作品与 tag
type WorkTag struct {
	WorkId PrimaryKey `gorm:"index"`
	FkWork Work       `gorm:"foreignKey:WorkId"`

	TagId PrimaryKey `gorm:"index"`
	FkTag Tag        `gorm:"foreignKey:TagId"`
}

func (WorkTag) TableName() string {
	return "project_tags"
}

// 项目分工表
type ProjectLaborDivision struct {
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

package dbmodel

// 成员偏好 tag
// 关联成员与 tag
type MemberPreference struct {
	Id PrimaryKey `gorm:"primaryKey;autoIncrement;not null"`

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

// 项目所携带 tag
// 关联项目与 tag
type ProjectTag struct {
	Id PrimaryKey `gorm:"primaryKey;autoIncrement;not null"`

	ProjectId PrimaryKey `gorm:"index"`
	FkProject Project    `gorm:"foreignKey:ProjectId"`

	TagId PrimaryKey `gorm:"index"`
	FkTag Tag        `gorm:"foreignKey:TagId"`
}

func (ProjectTag) TableName() string {
	return "project_tags"
}

// 项目分工表
type ProjectLaborDivision struct {
	Id PrimaryKey `gorm:"primaryKey;autoIncrement;not null"`

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

package dbmodels

// 项目分工表
type ProjectLaborDivision struct {
	BaseModel

	Id PrimaryKey `gorm:"primaryKey;autoIncrement;not null"`

	ProjectId PrimaryKey
	FkProject Project `gorm:"foreignKey:ProjectId"`

	UserId PrimaryKey
	FkUser User `gorm:"foreignKey:UserId"`

	// 分工，使用掩码计算多重身份
	LaborRole LaborMask `gorm:"not null"`
}

func (ProjectLaborDivision) TableName() string {
	return "project_labor_divisions"
}

// 成员在汉化组中的身份
type TeamMember struct {
	BaseModel

	UserId PrimaryKey `gorm:"index"`
	FkUser User       `gorm:"foreignKey:UserId"`

	TeamId PrimaryKey `gorm:"index"`
	FkTeam Team       `gorm:"foreignKey:TeamId"`

	// 职责 ，使用掩码计算多重身份
	Role LaborMask `gorm:"not null;default:0"`
}

func (TeamMember) TableName() string {
	return "team_members"
}

// Status 定义了申请状态的类型
type Status int

// 定义了申请状态的常量
const (
	STATUS_PENDING  Status = 0 // 申请待处理
	STATUS_ACCEPTED Status = 1 // 申请已接受
	STATUS_REJECTED Status = 2 // 申请已拒绝
)

// 项目邀请表
type ProjectInvitation struct {
	BaseModel

	InviterId PrimaryKey `gorm:"index"`
	FkInvitor User       `gorm:"foreignKey:InviterId"`
	InviteeId PrimaryKey `gorm:"index"`
	FkInvitee User       `gorm:"foreignKey:InviteeId"`

	ProjectId PrimaryKey `gorm:"index"`
	FkProject Project    `gorm:"foreignKey:ProjectId"`

	TargetRole LaborMask `gorm:"not null;default:0"` // 邀请的角色，使用位掩码表示

	Status Status `gorm:"not null;default:0"` // 邀请状态，0 pending, 1 accepted, 2 rejected
}

func (ProjectInvitation) TableName() string {
	return "project_invitations"
}

// ProjectApplication 项目申请表
type ProjectApplication struct {
	BaseModel

	ApplicantId PrimaryKey `gorm:"index"`
	FkApplicant User       `gorm:"foreignKey:ApplicantId"`
	ProjectId   PrimaryKey `gorm:"index"`
	FkProject   Project    `gorm:"foreignKey:ProjectId"`

	// 负责人/审核人
	PrincipalId PrimaryKey `gorm:"index;not null"`
	FkPrincipal User       `gorm:"foreignKey:PrincipalId"`

	TargetRole LaborMask `gorm:"not null;default:0"` // 申请的角色，使用位掩码表示

	Status Status `gorm:"not null;default:0"` // 申请状态，0 pending, 1 accepted, 2 rejected
}

func (ProjectApplication) TableName() string {
	return "project_applications"
}

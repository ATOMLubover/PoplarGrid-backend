package models

import "gorm.io/gorm"

// Invitation 定义了邀请的基本信息
type Invitation struct {
	BaseModel

	InvitorMemberId PKey    `gorm:"not null;index"`             // 邀请者成员 ID
	FkInvitor       *Member `gorm:"foreignKey:InvitorMemberId"` // 邀请者成员外键
	InviteeMemberId PKey    `gorm:"not null;index"`             // 接收者成员 ID
	FkInvitee       *Member `gorm:"foreignKey:InviteeMemberId"` // 接收者成员外键

	TargetProjectId PKey     `gorm:"not null;index"`             // 目标项目 ID
	FkProject       *Project `gorm:"foreignKey:TargetProjectId"` // 目标项目外键

	TargetLaborMask uint32 `gorm:"not null;default:0"` // 目标角色掩码

	Status Status `gorm:"not null;default:0"` // 邀请状态，0: 待处理, 1: 已接受, 2: 已拒绝
}

// TableName 返回 Invitation 的表名
func (*Invitation) TableName() string {
	return "invitations"
}

// InvitationSpec 定义了邀请的查询条件
type InvitationSpec struct {
	Id *PKey // 邀请 ID

	InvitorMemberId *PKey // 邀请者成员 ID
	InviteeMemberId *PKey // 接收者成员 ID
	TargetProjectId *PKey // 目标项目 ID

	Status *Status // 邀请状态
}

// Apply 将 InvitationSpec 应用为 WHERE 子句
func (s *InvitationSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.InvitorMemberId != nil {
		query = query.Where("invitor_member_id = ?", *s.InvitorMemberId)
	}

	if s.InviteeMemberId != nil {
		query = query.Where("invitee_member_id = ?", *s.InviteeMemberId)
	}

	if s.TargetProjectId != nil {
		query = query.Where("target_project_id = ?", *s.TargetProjectId)
	}
}

// InvitationFields 定义了 Invitation 的预加载字段
type InvitationFields struct {
	Id   bool // 邀请 ID
	Time bool // 邀请时间

	InvitorMemberId bool          // 邀请者成员 ID
	InviteeMemberId bool          // 接收者成员 ID
	MemberFields    *MemberFields // 邀请者和接收者的预加载字段

	TargetProjectId bool           // 目标项目 ID
	ProjectFields   *ProjectFields // 目标项目的预加载字段

	TargetLaborMask bool // 目标角色掩码
}

// Apply 将 InvitationFields 应用为 SELECT 子句
func (f *InvitationFields) Apply(query *gorm.DB) {
	fields := []string{}

	if f.Id {
		fields = append(fields, "id")
	}
	if f.Time {
		fields = append(fields, "time")
	}

	if f.InvitorMemberId {
		fields = append(fields, "invitor_member_id")
	}
	if f.InviteeMemberId {
		fields = append(fields, "invitee_member_id")
	}

	if f.TargetProjectId {
		fields = append(fields, "target_project_id")
		if f.ProjectFields != nil {
			// 如果指定了 ProjectFields，则添加项目相关字段
			query = query.Preload("FkProject", func(db *gorm.DB) *gorm.DB {
				f.ProjectFields.Apply(db)
				return db
			})
		}
	}

	if f.TargetLaborMask {
		fields = append(fields, "target_role_mask")
	}

	query.Select(fields)
}

// 供外部调用的指针，避免每一次都创建
var emptyInvitation *Invitation = nil

// GetInvitation 获取一个空的 Invitation 实例
func GetInvitation() *Invitation {
	return emptyInvitation
}

// Insert 插入一条新的邀请记录
func (*Invitation) Insert(handle *gorm.DB, invitation *Invitation) error {
	if invitation == nil {
		return &InvalidParameterError{}
	}

	return handle.
		Model(&Invitation{}).
		Create(invitation).Error
}

// Update 更新一条邀请记录
func (*Invitation) Update(handle *gorm.DB, invitation *Invitation) error {
	if invitation == nil {
		return &InvalidParameterError{}
	}

	if invitation.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return handle.
		Model(&Invitation{}).
		Where("id = ?", invitation.BaseModel.Id).
		Updates(invitation).Error
}

// SelectFirst 查询一条邀请记录
func (*Invitation) SelectFirst(handle *gorm.DB, cnd *InvitationSpec, fields *InvitationFields) (*Invitation, error) {
	var invitation Invitation

	query := handle.Model(&Invitation{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.First(&invitation).Error; err != nil {
		return nil, err
	}

	return &invitation, nil
}

// SelectMany 查询多条邀请记录
func (*Invitation) SelectMany(handle *gorm.DB, cnd *InvitationSpec, fields *InvitationFields) ([]*Invitation, error) {
	var invitations []*Invitation

	query := handle.Model(&Invitation{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

// Delete 删除一条邀请记录
func (*Invitation) Delete(handle *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return handle.
		Model(&Invitation{}).
		Where("id = ?", id).
		Delete(&Invitation{}).Error
}

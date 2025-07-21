package repos

import (
	"errors"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// InivitaionRepo 定义了邀请相关的仓库接口
type InvitationRepo interface {
	// GetInvitationSentByUserId 根据用户 ID 获取其发送的邀请信息
	GetInvitationSentByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.ProjectInvitation, error)
	// GetInvitationRecievedByUserId 根据用户 ID 获取收到的邀请信息
	GetInvitationRecievedByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.ProjectInvitation, error)

	// CreateInvitation 创建一个新的邀请
	CreateInvitation(inviter dbmodels.PrimaryKey, invitee dbmodels.PrimaryKey,
		projectId dbmodels.PrimaryKey, targetRole uint) error
}

// invitationRepoImpl 是 InvitationRepo 的实现
type invitationRepoImpl struct {
	handle *gorm.DB
}

// NewInvitationRepo 创建一个新的 InvitationRepo 实例
func NewInvitationRepo(handle *gorm.DB) InvitationRepo {
	return &invitationRepoImpl{
		handle: handle,
	}
}

// GetInvitationSentByUserId 实现 InvitationRepo 接口的 GetInvitationSentByUserId 方法
func (r *invitationRepoImpl) GetInvitationSentByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.ProjectInvitation, error) {
	var invitations []*dbmodels.ProjectInvitation

	if err := r.handle.Model(&dbmodels.ProjectInvitation{}).
		Where("inviter_id = ?", userId).
		// 预加载 FkProject 关联的 Project 信息
		Preload("FkProject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "workset_id", "workset_index", "status")
		}).
		// 预加载 FkInvitee 关联的 User 信息
		Preload("FkInvitee", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nickname") // User 只需要 id、nickname
		}).
		Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

// GetInvitationRecievedByUserId 实现 InvitationRepo 接口的 GetInvitationRecievedByUserId 方法
func (r *invitationRepoImpl) GetInvitationRecievedByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.ProjectInvitation, error) {
	var invitations []*dbmodels.ProjectInvitation

	if err := r.handle.Model(&dbmodels.ProjectInvitation{}).
		Where("invitee_id = ?", userId).
		// 预加载 FkProject 关联的 Project 信息
		Preload("FkProject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "workset_id", "workset_index", "status")
		}).
		// 预加载 FkInvitee 关联的 User 信息
		Preload("FkInvitee", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nickname") // User 只需要 id、nickname
		}).
		Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

// CreateInvitation 实现 InvitationRepo 接口的 CreateInvitation 方法
func (r *invitationRepoImpl) CreateInvitation(
	inviter dbmodels.PrimaryKey, invitee dbmodels.PrimaryKey,
	projectId dbmodels.PrimaryKey, targetRole uint,
) error {
	if targetRole == 0 {
		return errors.New("目标职位不可以为空")
	}

	invitation := &dbmodels.ProjectInvitation{
		InviterId: inviter,
		InviteeId: invitee,
		ProjectId: projectId,

		TargetRole: dbmodels.LaborMask(targetRole),
	}

	if err := r.handle.Create(invitation).Error; err != nil {
		return err
	}

	return nil
}

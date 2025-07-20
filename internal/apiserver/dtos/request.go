package dtos

// CreateProjectRequest 定义了创建项目的请求体
type CreateProjectRequest struct {
	Title       string `json:"title" binding:"required"`      // 项目标题，必填
	WorksetId   uint   `json:"workset_id" binding:"required"` // 作品集 ID，必填
	Description string `json:"description"`                   // 项目描述

	AllowAutoJoin bool `json:"allow_auto_join"` // 是否允许自动加入，默认为 false
	IsHidden      bool `json:"is_hidden"`       // 是否隐藏项目，默认为 false
}

// UpdateProjectRequest 定义了更新项目的请求体
type UpdateProjectRequest struct {
	ProjectId   uint   `json:"project_id" binding:"required"` // 项目 ID，必填
	Title       string `json:"title" binding:"required"`      // 项目标题
	Description string `json:"description"`                   // 项目描述
	Status      uint   `json:"status"`                        // 项目状态，使用位掩码表示
}

// InviteMemberRequest 定义了邀请成员加入项目的请求体
type InviteMemberRequest struct {
	MemberID   uint `json:"member_id" binding:"required"`   // 成员的 ID
	InviteRole uint `json:"invite_role" binding:"required"` // 邀请加入的职位 (使用 uint 存储位掩码)
	ProjectID  uint `json:"project_id" binding:"required"`  // 项目 ID
}

// AcceptInvitationRequest 定义了接受邀请加入项目的请求体
type AcceptInvitationRequest struct {
	InvitationID uint `json:"invitation_id" binding:"required"` // 邀请的 ID
}

// RefuseInvitationRequest 定义了拒绝邀请加入项目的请求体
type RefuseInvitationRequest struct {
	InvitationID uint `json:"invitation_id" binding:"required"` // 邀请的 ID
}

// ApplyProjectRequest 定义了申请加入项目的请求体
type ApplyProjectRequest struct {
	ProjectID uint `json:"project_id" binding:"required"` // 项目 ID
	ApplyRole uint `json:"apply_role" binding:"required"` // 申请加入的职位 (使用 uint 存储位掩码)
}

// AcceptApplicationRequest 定义了接受申请加入项目的请求体
type AcceptApplicationRequest struct {
	ApplicationID uint `json:"application_id" binding:"required"` // 申请的 ID
}

// RefuseApplicationRequest 定义了拒绝申请加入项目的请求体
type RefuseApplicationRequest struct {
	ApplicationID uint `json:"application_id" binding:"required"` // 申请的 ID
}

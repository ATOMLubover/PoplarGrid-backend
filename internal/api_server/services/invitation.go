package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"
	"time"

	"gorm.io/gorm"
)

// InvitationListParams 定义了邀请列表的查询参数
type InvitationListParams struct {
	Offset           int               // 偏移量
	Limit            int               // 限制数量
	CurrentMemberIds map[uint]struct{} // 当前成员 ID 集合
	InvitorMemberId  uint              // 邀请者成员 ID
	InviteeMemberId  uint              // 接收者成员 ID
	TargetProjectId  uint              // 目标项目 ID
}

// InvitationInfo 定义了邀请的基本信息
type InvitationInfo struct {
	Id   uint      // 邀请 ID
	Time time.Time // 邀请时间

	InvitorMemberId uint   // 邀请者成员 ID
	InvitorNickname string // 邀请者昵称
	InviteeMemberId uint   // 接收者成员 ID
	InviteeNickname string // 接收者昵称

	TargetProjectId        uint   // 目标项目 ID
	TargetProjectTitle     string // 目标项目标题
	TargetProjectWorksetId uint   // 目标工作集 ID
	TargetProjectIndex     uint   // 目标工作集索引

	TargetLaborMask LaborMask // 目标角色掩码
}

// CreateInvitationParams 定义了创建邀请的请求参数
type CreateInvitationParams struct {
	CurrentMemberIds map[uint]struct{} // 当前成员 ID 集合
	InvitorMemberId  uint              // 邀请者成员 ID
	InviteeMemberId  uint              // 接收者成员 ID
	TargetProjectId  uint              // 目标项目 ID
	TargetLaborMask  LaborMask         // 目标角色掩码
}

// ProcessInvitationParams 定义了处理邀请的请求参数
type ProcessInvitationParams struct {
	CurrentMemberIds  map[uint]struct{} // 当前成员 ID 集合
	ProcessorMemberId uint              // 处理者成员 ID
	InvitationId      uint              // 邀请 ID
	Accept            bool              // 是否接受邀请
}

// InvitationService 接口定义了邀请服务的基本操作
type InvitationService interface {
	// GetInvitations 获取指定条件下的邀请列表
	GetInvitations(params *InvitationListParams) ([]*InvitationInfo, error)
	// CreateInvitation 创建新的邀请
	CreateInvitation(params *CreateInvitationParams) error
	// ProcessInvitation 处理邀请（接受或拒绝）
	ProcessInvitation(params *ProcessInvitationParams) error
}

// invitationServiceImpl 是 InvitationService 的实现
type invitationServiceImpl struct {
	handle    *gorm.DB
	apiClient apiclient.ApiClient
	logger    *slog.Logger
}

// NewInvitationService 创建一个新的 InvitationService 实例
func NewInvitationService(
	hdl *gorm.DB,
	apiClient apiclient.ApiClient,
	lgr *slog.Logger,
) InvitationService {
	return &invitationServiceImpl{
		handle:    hdl,
		apiClient: apiClient,
		logger:    lgr,
	}
}

// GetInvitations 实现 InvitationService 接口的 GetInvitations 方法
func (s *invitationServiceImpl) GetInvitations(params *InvitationListParams) ([]*InvitationInfo, error) {
	// 检查是否查询的 invitorMeberId、inviteeMemberId 属于当前用户的成员 ID 集合
	if _, ok := params.CurrentMemberIds[params.InvitorMemberId]; params.InvitorMemberId != 0 && !ok {
		// 只有 invitorMeberId 不为默认值的 0 时才检查
		s.logger.Error("GetInvitations 邀请者不在当前成员列表中",
			slog.Uint64("invitor_id", uint64(params.InvitorMemberId)))
		return nil, errors.New("邀请者不在当前成员列表中")
	}
	if _, ok := params.CurrentMemberIds[params.InviteeMemberId]; params.InviteeMemberId != 0 && !ok {
		// 只有 inviteeMemberId 不为默认值的 0 时才检查
		s.logger.Error("GetInvitations 接受者不在当前成员列表中",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)))
		return nil, errors.New("接受者不在当前成员列表中")
	}

	// 查询条件为邀请者或接收者成员 ID
	invitationSpec := buildQueryParams(params)
	// 查询的字段
	invitationFields := &models.InvitationFields{
		Id:              true,
		Time:            true,
		InvitorMemberId: true,
		InviteeMemberId: true,
		MemberFields: &models.MemberFields{
			Id:     true,
			UserId: true,
			UserFields: &models.UserFields{
				Id:       true,
				Nickname: true, // 需要获取被邀请者的昵称
			},
		},
		TargetProjectId: true,
		ProjectFields: &models.ProjectFields{
			Id:           true,
			Title:        true,
			WorksetId:    true,
			WorksetIndex: true,
		},
		TargetLaborMask: true,
	}

	// 执行查询
	invitations, err := models.GetInvitation().SelectMany(s.handle, invitationSpec, invitationFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("GetInvitations 没有找到符合条件的邀请",
				slog.Uint64("invitor_id", uint64(params.InvitorMemberId)),
				slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
				slog.Uint64("target_project_id", uint64(params.TargetProjectId)))
			return nil, errors.New("没有找到符合条件的邀请")
		}
		s.logger.Error("GetInvitations 查询邀请失败", slog.Any("error", err))
		return nil, errors.New("查询邀请列表失败")
	}

	// 将查询结果转换为 InvitationInfo
	invitationInfos := make([]*InvitationInfo, 0, len(invitations))

	for _, invitation := range invitations {
		invitationInfos = append(invitationInfos, &InvitationInfo{
			Id:                     uint(invitation.Id),
			InvitorMemberId:        uint(invitation.InvitorMemberId),
			InviteeMemberId:        uint(invitation.InviteeMemberId),
			TargetProjectId:        uint(invitation.TargetProjectId),
			TargetProjectTitle:     invitation.FkProject.Title,
			TargetProjectWorksetId: uint(invitation.FkProject.WorksetId),
			TargetProjectIndex:     uint(invitation.FkProject.WorksetIndex),
			TargetLaborMask:        LaborMask(invitation.TargetLaborMask),
		})
	}

	return invitationInfos, nil
}

// CreateInvitation 实现 InvitationService 接口的 CreateInvitation 方法
func (s *invitationServiceImpl) CreateInvitation(params *CreateInvitationParams) error {
	// 检查邀请者是否在当前成员 ID 集合中
	if _, exists := params.CurrentMemberIds[params.InvitorMemberId]; !exists {
		s.logger.Error("CreateInvitation 邀请者不在当前成员列表中",
			slog.Uint64("invitor_id", uint64(params.InvitorMemberId)))
		return errors.New("邀请者不在当前成员列表中")
	}

	// 检查受邀者是否能够担当指定的分工
	memberPKey := models.PKey(params.InviteeMemberId)
	memberSpec := &models.MemberSpec{
		Id: &memberPKey,
	}
	memberFields := &models.MemberFields{
		Id:    true,
		Roles: true, // 需要获取成员的角色掩码
	}

	member, err := models.GetMember().SelectFirst(s.handle, memberSpec, memberFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("CreateInvitation 没有找到指定成员",
				slog.Uint64("member_id", uint64(params.InviteeMemberId)))
			return errors.New("没有找到指定成员")
		}
		s.logger.Error("CreateInvitation 查询指定成员失败",
			slog.Uint64("member_id", uint64(params.InviteeMemberId)),
			slog.Any("error", err))
		return errors.New("查询指定成员失败")
	}

	// 检查受邀者是否有足够的权限
	if err := checkValidLabor(member, params.TargetLaborMask); err != nil {
		s.logger.Error("CreateInvitation 受邀者没有足够的权限",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
			slog.Any("target_labor_mask", params.TargetLaborMask),
			slog.Any("error", err))
		return errors.New("受邀者没有足够的权限")
	}

	// 检查邀请者是否是指定项目的负责人
	projectPKey := models.PKey(params.TargetProjectId)
	projectSpec := &models.ProjectSpec{
		Id: &projectPKey,
	}
	projectFields := &models.ProjectFields{
		Id:          true,
		PrincipalId: true,
	}

	project, err := models.GetProject().SelectFirst(s.handle, projectSpec, projectFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("CreateInvitation 没有找到指定项目",
				slog.Uint64("project_id", uint64(params.TargetProjectId)))
			return errors.New("没有找到指定项目")
		}
		s.logger.Error("CreateInvitation 查询项目信息失败",
			slog.Uint64("project_id", uint64(params.TargetProjectId)),
			slog.Any("error", err))
		return errors.New("查询项目信息失败")
	}

	if project.PrincipalId != models.PKey(params.InvitorMemberId) {
		// 如果不是邀请者是项目负责人，则返回错误
		s.logger.Error("CreateInvitation 邀请者不是项目负责人",
			slog.Uint64("invitor_id", uint64(params.InvitorMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)))
		return errors.New("邀请者不是项目负责人")
	}

	// 检查接受者是否已经在当前项目
	inviteePKey := models.PKey(params.InviteeMemberId)
	laborSpec := &models.LaborSpec{
		MemberId:  &inviteePKey,
		ProjectId: &projectPKey,
	}
	laborFields := &models.LaborFields{
		Id: true, // 其余字段可以不需要，只要检查存在性即可
	}

	labor, err := models.GetLabor().SelectFirst(s.handle, laborSpec, laborFields)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("CreateInvitation 查询接受者在项目中的分工记录失败",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)),
			slog.Any("error", err))
		return errors.New("查询接受者在项目中的分工记录失败")
	}

	if labor != nil {
		// 如果接受者已经在当前项目，则返回错误
		s.logger.Error("CreateInvitation 接受者已经在当前项目中",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)))
		return errors.New("接受者已经在当前项目中")
	}

	// 再检查是否有重复的未处理邀请记录
	invitationStatus := models.Status(0)
	invitationSpec := &models.InvitationSpec{
		InviteeMemberId: &inviteePKey,
		TargetProjectId: &projectPKey,
		Status:          &invitationStatus,
	}
	invitationFields := &models.InvitationFields{
		Id: true, // 其余字段可以不需要，只要检查存在性即可
	}

	invitations, err := models.GetInvitation().SelectMany(s.handle, invitationSpec, invitationFields)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("CreateInvitation 查询未处理的邀请记录失败",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)),
			slog.Any("error", err))
		return errors.New("查询未处理的邀请记录失败")
	}

	if len(invitations) > 0 {
		// 如果存在未处理的邀请记录，则返回错误
		s.logger.Error("CreateInvitation 存在未处理的邀请记录",
			slog.Uint64("invitee_id", uint64(params.InviteeMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)))
		return errors.New("存在未处理的邀请记录")
	}

	// 创建新的邀请记录
	invitation := &models.Invitation{
		InvitorMemberId: models.PKey(params.InvitorMemberId),
		InviteeMemberId: models.PKey(params.InviteeMemberId),
		TargetProjectId: models.PKey(params.TargetProjectId),
		TargetLaborMask: uint32(params.TargetLaborMask),
	}

	if err := models.GetInvitation().Insert(s.handle, invitation); err != nil {
		s.logger.Error("CreateInvitation 插入邀请记录失败", slog.Any("error", err))
		return errors.New("创建邀请失败")
	}

	return nil
}

// ProcessInvitation 实现 InvitationService 接口的 ProcessInvitation 方法
func (s *invitationServiceImpl) ProcessInvitation(params *ProcessInvitationParams) error {
	// 检查邀请 ID 是否在当前成员 ID 集合中
	if _, exists := params.CurrentMemberIds[params.ProcessorMemberId]; !exists {
		s.logger.Error("ProcessInvitation 处理成员 ID 不在当前成员列表中",
			slog.Uint64("processor_member_id", uint64(params.ProcessorMemberId)))
		return errors.New("处理成员 ID 不在当前成员列表中")
	}

	// 查询邀请记录
	invitationPKey := models.PKey(params.InvitationId)
	newStatus := models.Status(0)

	if err := models.GetInvitation().Update(s.handle, &models.Invitation{
		BaseModel: models.BaseModel{
			Id: invitationPKey,
		},
		Status: newStatus,
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("ProcessInvitation 没有找到指定邀请记录",
				slog.Uint64("invitation_id", uint64(params.InvitationId)))
			return errors.New("没有找到指定邀请记录")
		}
		s.logger.Error("ProcessInvitation 更新邀请记录失败",
			slog.Uint64("invitation_id", uint64(params.InvitationId)),
			slog.Any("error", err))
		return errors.New("更新邀请记录失败")
	}

	return nil
}

// ========= 辅助函数 =========

// buildQueryParams 构建查询参数
func buildQueryParams(params *InvitationListParams) *models.InvitationSpec {
	invitationSpec := &models.InvitationSpec{}

	if params.InvitorMemberId != 0 {
		invitorPKey := models.PKey(params.InvitorMemberId)
		invitationSpec.InvitorMemberId = &invitorPKey
	}
	if params.InviteeMemberId != 0 {
		inviteePKey := models.PKey(params.InviteeMemberId)
		invitationSpec.InviteeMemberId = &inviteePKey
	}
	if params.TargetProjectId != 0 {
		targetProjectPKey := models.PKey(params.TargetProjectId)
		invitationSpec.TargetProjectId = &targetProjectPKey
	}

	return invitationSpec
}

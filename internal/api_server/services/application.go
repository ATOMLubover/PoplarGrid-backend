package services

import (
	"context"
	"errors"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"
	"poplargrid/internal/shared/txutils"
	"time"

	"gorm.io/gorm"
)

// ApplicationListParams 定义了申请列表的查询参数
type ApplicationListParams struct {
	Offset            int               // 偏移量
	Limit             int               // 限制数量
	CurrentMemberIds  map[uint]struct{} // 当前成员 ID 集合
	ApplicantMemberId uint              // 申请者成员 ID
	ProcessorMemberId uint              // 处理者成员 ID
	TargetProjectId   uint              // 目标项目 ID
}

// ApplicationInfo 定义了申请的基本信息
type ApplicationInfo struct {
	Id   uint      // 申请 ID
	Time time.Time // 申请时间

	ApplicantMemberId uint   // 申请者成员 ID
	ApplicantNickname string // 申请者昵称
	ProcessorMemberId uint   // 处理者成员 ID
	ProcessorNickname string // 处理者昵称

	TargetProjectId        uint   // 目标项目 ID
	TargetProjectTitle     string // 目标项目标题
	TargetProjectWorksetId uint   // 目标工作集 ID
	TargetProjectIndex     uint   // 目标工作集索引

	TargetLaborMask LaborMask // 目标角色掩码
}

// CreateApplicationParams 定义了创建申请的请求参数
type CreateApplicationParams struct {
	CurrentMemberIds  map[uint]struct{} // 当前成员 ID 集合
	ApplicantMemberId uint              // 申请者成员 ID
	TargetProjectId   uint              // 目标项目 ID
	TargetLaborMask   LaborMask         // 目标角色掩码
}

// ProcessApplicationParams 定义了处理申请的请求参数
type ProcessApplicationParams struct {
	CurrentMemberIds  map[uint]struct{} // 当前成员 ID 集合
	ProcessorMemberId uint              // 处理者成员 ID
	ApplicationId     uint              // 申请 ID
	Accept            bool              // 是否接受申请
}

// ApplicationService 接口定义了申请服务的基本操作
type ApplicationService interface {
	// GetApplications 获取指定条件下的申请列表
	GetApplications(params *ApplicationListParams) ([]*ApplicationInfo, error)
	// CreateApplication 创建新的申请
	CreateApplication(params *CreateApplicationParams) error
	// ProcessApplication 处理申请（接受或拒绝）
	ProcessApplication(params *ProcessApplicationParams) error
}

// applicationServiceImpl 是 ApplicationService 的实现
type applicationServiceImpl struct {
	handle    *gorm.DB
	apiClient apiclient.ApiClient
	logger    *slog.Logger
}

// NewApplicationService 创建一个新的 ApplicationService 实例
func NewApplicationService(
	hdl *gorm.DB,
	apiClient apiclient.ApiClient,
	lgr *slog.Logger,
) ApplicationService {
	return &applicationServiceImpl{
		handle:    hdl,
		apiClient: apiClient,
		logger:    lgr,
	}
}

// GetApplications 实现 ApplicationService 接口的 GetApplications 方法
func (s *applicationServiceImpl) GetApplications(params *ApplicationListParams) ([]*ApplicationInfo, error) {
	// 检查是否查询的 applicantMemberId、processorMemberId 属于当前用户的成员 ID 集合
	if _, ok := params.CurrentMemberIds[params.ApplicantMemberId]; params.ApplicantMemberId != 0 && !ok {
		// 只有 applicantMemberId 不为默认值的 0 时才检查
		s.logger.Error("GetApplications 申请者不在当前成员列表中",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)))
		return nil, errors.New("申请者不在当前成员列表中")
	}
	if _, ok := params.CurrentMemberIds[params.ProcessorMemberId]; params.ProcessorMemberId != 0 && !ok {
		// 只有 processorMemberId 不为默认值的 0 时才检查
		s.logger.Error("GetApplications 处理者不在当前成员列表中",
			slog.Uint64("processor_id", uint64(params.ProcessorMemberId)))
		return nil, errors.New("处理者不在当前成员列表中")
	}

	// 查询条件为申请者或处理者成员 ID
	applicationSpec := buildApplicationQueryParams(params)
	// 查询的字段
	applicationFields := &models.ApplicationFields{
		Id:                true,
		Time:              true,
		ApplicantMemberId: true,
		MemberFields: &models.MemberFields{
			Id:     true,
			UserId: true,
			UserFields: &models.UserFields{
				Id:       true,
				Nickname: true, // 需要获取处理者的昵称
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
	applications, err := models.GetApplication().SelectMany(s.handle, applicationSpec, applicationFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("GetApplications 没有找到符合条件的申请",
				slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
				slog.Uint64("processor_id", uint64(params.ProcessorMemberId)),
				slog.Uint64("target_project_id", uint64(params.TargetProjectId)))
			return nil, errors.New("没有找到符合条件的申请")
		}
		s.logger.Error("GetApplications 查询申请失败", slog.Any("error", err))
		return nil, errors.New("查询申请列表失败")
	}

	// 将查询结果转换为 ApplicationInfo
	applicationInfos := make([]*ApplicationInfo, 0, len(applications))

	for _, application := range applications {
		applicationInfos = append(applicationInfos, &ApplicationInfo{
			Id:                     uint(application.Id),
			ApplicantMemberId:      uint(application.ApplicantMemberId),
			ProcessorMemberId:      uint(application.ProcessorMemberId),
			TargetProjectId:        uint(application.TargetProjectId),
			TargetProjectTitle:     application.FkProject.Title,
			TargetProjectWorksetId: uint(application.FkProject.WorksetId),
			TargetProjectIndex:     uint(application.FkProject.WorksetIndex),
			TargetLaborMask:        LaborMask(application.TargetLaborMask),
		})
	}

	return applicationInfos, nil
}

// CreateApplication 实现 ApplicationService 接口的 CreateApplication 方法
func (s *applicationServiceImpl) CreateApplication(params *CreateApplicationParams) error {
	// 检查申请者是否在当前成员 ID 集合中
	if _, exists := params.CurrentMemberIds[params.ApplicantMemberId]; !exists {
		s.logger.Error("CreateApplication 申请者不在当前成员列表中",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)))
		return errors.New("申请者不在当前成员列表中")
	}

	// 检查当前申请者是否能够担任申请的角色
	memberPKey := models.PKey(params.ApplicantMemberId)
	memberSpec := &models.MemberSpec{
		Id: &memberPKey,
	}
	memberFields := &models.MemberFields{
		Id:    true,
		Roles: true, // 需要检查申请者的角色
	}

	member, err := models.GetMember().SelectFirst(s.handle, memberSpec, memberFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("CreateApplication 没有找到指定申请者成员",
				slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
				slog.Any("error", err))
			return errors.New("没有找到指定申请者成员")
		}
		s.logger.Error("CreateApplication 查询申请者成员信息失败",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Any("error", err))
		return errors.New("查询申请者成员信息失败")
	}

	if err := checkValidLabor(member, params.TargetLaborMask); err != nil {
		s.logger.Error("CreateApplication 申请者不能担任申请的分工",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Any("error", err))
		return errors.New("申请者不能担任申请的分工")
	}

	// 检查申请者是否已经在当前项目
	applicantPKey := models.PKey(params.ApplicantMemberId)
	projectPKey := models.PKey(params.TargetProjectId)
	laborSpec := &models.LaborSpec{
		MemberId:  &applicantPKey,
		ProjectId: &projectPKey,
	}
	laborFields := &models.LaborFields{
		Id: true, // 其余字段可以不需要，只要检查存在性即可
	}

	labor, err := models.GetLabor().SelectFirst(s.handle, laborSpec, laborFields)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("CreateApplication 查询申请者在项目中的分工记录失败",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)),
			slog.Any("error", err))
		return errors.New("查询申请者在项目中的分工记录失败")
	}

	if labor != nil {
		// 如果申请者已经在当前项目，则返回错误
		s.logger.Error("CreateApplication 申请者已经在当前项目中",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)))
		return errors.New("申请者已经在当前项目中")
	}

	// 再检查是否有重复的未处理申请记录
	applicationStatus := models.Status(0)
	applicationSpec := &models.ApplicationSpec{
		ApplicantMemberId: &applicantPKey,
		TargetProjectId:   &projectPKey,
		Status:            &applicationStatus,
	}
	applicationFields := &models.ApplicationFields{
		Id: true, // 其余字段可以不需要，只要检查存在性即可
	}

	applications, err := models.GetApplication().SelectMany(s.handle, applicationSpec, applicationFields)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("CreateApplication 查询未处理的申请记录失败",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)),
			slog.Any("error", err))
		return errors.New("查询未处理的申请记录失败")
	}

	if len(applications) > 0 {
		// 如果存在未处理的申请记录，则返回错误
		s.logger.Error("CreateApplication 存在未处理的申请记录",
			slog.Uint64("applicant_id", uint64(params.ApplicantMemberId)),
			slog.Uint64("project_id", uint64(params.TargetProjectId)))
		return errors.New("存在未处理的申请记录")
	}

	// 创建新的申请记录
	application := &models.Application{
		ApplicantMemberId: models.PKey(params.ApplicantMemberId),
		TargetProjectId:   models.PKey(params.TargetProjectId),
		TargetLaborMask:   uint32(params.TargetLaborMask),
	}

	if err := models.GetApplication().Insert(s.handle, application); err != nil {
		s.logger.Error("CreateApplication 插入申请记录失败", slog.Any("error", err))
		return errors.New("创建申请失败")
	}

	return nil
}

// ProcessApplication 实现 ApplicationService 接口的 ProcessApplication 方法
func (s *applicationServiceImpl) ProcessApplication(params *ProcessApplicationParams) error {
	// 检查处理成员 ID 是否在当前成员 ID 集合中
	if _, exists := params.CurrentMemberIds[params.ProcessorMemberId]; !exists {
		s.logger.Error("ProcessApplication 处理成员 ID 不在当前成员列表中",
			slog.Uint64("processor_member_id", uint64(params.ProcessorMemberId)))
		return errors.New("处理成员 ID 不在当前成员列表中")
	}

	// 创建事务协调器
	c := txutils.NewTransactionCoordinator(s.handle)

	// 在一个事务中处理更改申请状态和添加分工记录
	if err := c.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 查询申请记录
		applicationPKey := models.PKey(params.ApplicationId)
		applicationSpec := &models.ApplicationSpec{
			Id: &applicationPKey,
		}
		applicationFields := &models.ApplicationFields{
			Id:              true,
			TargetProjectId: true,
			ProjectFields: &models.ProjectFields{
				Id:        true,
				MoetranId: true, // 需要获取龙译 ID
			},
			Status: true,
		}

		application, err := models.GetApplication().SelectFirst(s.handle, applicationSpec, applicationFields)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("ProcessApplication 没有找到指定申请记录",
					slog.Uint64("application_id", uint64(params.ApplicationId)))
				return errors.New("没有找到指定申请记录"), nil
			}
			s.logger.Error("ProcessApplication 查询申请记录失败",
				slog.Uint64("application_id", uint64(params.ApplicationId)),
				slog.Any("error", err))
			return errors.New("查询申请记录失败"), nil
		}

		// 检查处理者是否是申请的处理者
		if application.ProcessorMemberId != models.PKey(params.ProcessorMemberId) {
			s.logger.Error("ProcessApplication 处理者不是申请的处理者",
				slog.Uint64("application_id", uint64(params.ApplicationId)),
				slog.Uint64("processor_member_id", uint64(params.ProcessorMemberId)))
			return errors.New("处理者不是申请的处理者"), nil
		}

		// 更新申请状态
		var newStatus models.Status

		switch params.Accept {
		case true:
			newStatus = models.STATUS_ACCEPTED
		case false:
			newStatus = models.STATUS_REJECTED
		}

		if err := models.GetApplication().Update(s.handle, &models.Application{
			BaseModel: models.BaseModel{
				Id: application.Id,
			},
			Status: newStatus,
		}); err != nil {
			s.logger.Error("ProcessApplication 更新申请记录失败",
				slog.Uint64("application_id", uint64(params.ApplicationId)),
				slog.Any("error", err))
			return errors.New("更新申请记录失败"), nil
		}

		// 如果不接受申请，则略过创建新的分工记录
		if !params.Accept {
			return nil, nil
		}

		labor := &models.Labor{
			MemberId:  models.PKey(application.ApplicantMemberId),
			ProjectId: models.PKey(application.TargetProjectId),
			LaborMask: application.TargetLaborMask,
		}
		if err := models.GetLabor().Insert(s.handle, labor); err != nil {
			s.logger.Error("ProcessApplication 创建分工记录失败", slog.Any("error", err))
			return errors.New("创建分工记录失败"), nil
		}

		// 利用龙译发出邀请
		// 获取两者的龙译 ID 和 JWT 等
		memberFields := &models.MemberFields{
			Id: true,
			UserFields: &models.UserFields{
				Id:         true,
				MoetranId:  true, // 需要获取龙译 ID
				MoetranJwt: true, // 需要获取用户的 JWT
			},
		}
		applicantSpec := &models.MemberSpec{
			Id: &application.ApplicantMemberId,
		}
		processorSpec := &models.MemberSpec{
			Id: &application.ProcessorMemberId,
		}

		applicant, err := models.GetMember().SelectFirst(s.handle, applicantSpec, memberFields)
		if err != nil {
			s.logger.Error("ProcessApplication 查询申请者成员信息失败",
				slog.Uint64("applicant_member_id", uint64(application.ApplicantMemberId)),
				slog.Any("error", err))
			return errors.New("查询申请者成员信息失败"), nil
		}
		processor, err := models.GetMember().SelectFirst(s.handle, processorSpec, memberFields)
		if err != nil {
			s.logger.Error("ProcessApplication 查询处理者成员信息失败",
				slog.Uint64("processor_member_id", uint64(application.ProcessorMemberId)),
				slog.Any("error", err))
			return errors.New("查询处理者成员信息失败"), nil
		}

		inviteInfo := &apiclient.InviteMemberParams{
			MoetranAuth:      processor.FkUser.MoetranJwt,
			MoetranProjectId: application.FkProject.MoetranId,
			MoetranInviteeID: applicant.FkUser.MoetranId,
			InviteRole:       roleMaskToMoetranRole(LaborMask(application.TargetLaborMask)),
		}

		if _, err := s.apiClient.InviteMemberToProject(inviteInfo); err != nil {
			s.logger.Error("ProcessApplication 邀请成员到龙译项目失败",
				slog.Uint64("applicant_member_id", uint64(application.ApplicantMemberId)),
				slog.Uint64("processor_member_id", uint64(application.ProcessorMemberId)),
				slog.Any("error", err))
			return errors.New("邀请成员到龙译项目失败"), nil
		}

		return nil, nil

	}); err != nil {
		s.logger.Error("ProcessApplication 处理申请失败", slog.Any("error", err))
		return errors.New("处理申请失败")
	}

	return nil
}

// ========= 辅助函数 =========

// buildApplicationQueryParams 构建查询参数
func buildApplicationQueryParams(params *ApplicationListParams) *models.ApplicationSpec {
	applicationSpec := &models.ApplicationSpec{}

	if params.ApplicantMemberId != 0 {
		applicantPKey := models.PKey(params.ApplicantMemberId)
		applicationSpec.ApplicantMemberId = &applicantPKey
	}
	if params.ProcessorMemberId != 0 {
		processorPKey := models.PKey(params.ProcessorMemberId)
		applicationSpec.ProcessorMemberId = &processorPKey
	}
	if params.TargetProjectId != 0 {
		targetProjectPKey := models.PKey(params.TargetProjectId)
		applicationSpec.TargetProjectId = &targetProjectPKey
	}

	return applicationSpec
}

// checkValidLabor 检查申请者能否担任其申请的分工
func checkValidLabor(member *models.Member, laborMask LaborMask) error {
	if member == nil {
		return errors.New("申请者成员信息不能为空")
	}

	// 组装 member 的职责掩码
	memberLabor := buildLaborMask(member)

	if !memberLabor.HasRole(laborMask) {
		return errors.New("申请者不能担任申请的分工")
	}

	return nil
}

// roleMaskToMoetranRole 将角色掩码转换为龙译角色 ID
func roleMaskToMoetranRole(mask LaborMask) apiclient.MoetranRole {
	switch {
	case mask.HasRole(LABOR_PRINCIPAL_MASK):
		// 如果是负责人，再任命也只基于监理的职位
		return apiclient.ROLE_SUPERVISOR
	case mask.HasRole(LABOR_TRANSLATOR_MASK):
		return apiclient.ROLE_TRANSLATOR
	case mask.HasRole(LABOR_PROOFREADER_MASK):
		return apiclient.ROLE_PROOFREADER
	case mask.HasRole(LABOR_LETTERER_MASK):
		return apiclient.ROLE_EMBEDDER
	case mask.HasRole(LABOR_REVIEWER_MASK):
		return apiclient.ROLE_SUPERVISOR
	default:
		// 如果没有匹配的角色，返回实习翻译
		return apiclient.ROLE_INTERN
	}
}

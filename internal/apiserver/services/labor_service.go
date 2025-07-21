package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
)

// LaborService 接口定义了邀请申请服务的基本操作
// 懒得写两个 service 了，直接整合 Ov<
type LaborService interface {
	// GetInvitationSentByUserId 获取指定用户发送的邀请列表
	GetInvitationSentByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.InvitationBasic, error)
	// GetInvitationRecievedByUserId 获取指定用户收到的邀请列表
	GetInvitationRecievedByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.InvitationBasic, error)

	// CreateInvitation 创建一个新的邀请
	CreateInvitation(inviterId, inviteeId, projectId uint, targetRole uint) error

	// GetApplisRecievedByUserId 获取指定用户收到的申请列表，根据项目数量分页
	GetApplisRecievedByUserId(projectId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error)
	// GetApplisSentByUserId 获取指定用户发出的申请列表
	GetApplisSentByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error)

	// CreateApplication 创建一个新的申请
	CreateApplication(applicantId, projectId uint, targetRole uint) error
}

// laborServiceImpl 是 LaborService 接口的实现
type laborServiceImpl struct {
	invitationRepo repos.InvitationRepo
	appliRepo      repos.AppliRepo
	logger         *slog.Logger
}

// NewLaborService 创建一个新的 LaborService 实例
func NewLaborService(
	invitationRepo repos.InvitationRepo,
	appliRepo repos.AppliRepo,
	logger *slog.Logger,
) LaborService {
	return &laborServiceImpl{
		invitationRepo: invitationRepo,
		appliRepo:      appliRepo,
		logger:         logger,
	}
}

// GetInvitationSentByUserId 实现 LaborService 接口的 GetInvitationSentByUserId 方法
func (s *laborServiceImpl) GetInvitationSentByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.InvitationBasic, error) {
	invitations, err := s.invitationRepo.GetInvitationSentByUserId(dbmodels.PrimaryKey(userId))
	if err != nil {
		s.logger.Error("GetInvitationSentByUserId 调用 GetInvitationSentByUserId 中出现错误",
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err))
		return nil, err
	}

	var dtosInvitations []*dtos.InvitationBasic

	for _, invitation := range invitations {
		// 封装 project 为 overall status
		overallStatus := dtos.ProjectOverallStatus(0)

		overallStatus.SetTranslatingStatus(uint(invitation.FkProject.TranslateStatus))
		overallStatus.SetReviewingStatus(uint(invitation.FkProject.ReviewStatus))
		overallStatus.SetLetteringStatus(uint(invitation.FkProject.LetterStatus))
		overallStatus.SetProofreadingStatus(uint(invitation.FkProject.ProofStatus))
		overallStatus.SetReviewingStatus(uint(invitation.FkProject.ReviewStatus))
		switch invitation.FkProject.IsPublished {
		case true:
			overallStatus.SetPublishedStatus(dtos.PROJECT_STATUS_COMPLETED)
		case false:
			overallStatus.SetPublishedStatus(dtos.PROJECT_STATUS_UNSET)
		}

		dtosInvitations = append(dtosInvitations, &dtos.InvitationBasic{
			Id:              uint(invitation.Id),
			InviterId:       uint(invitation.InviterId),
			InviteeId:       uint(invitation.InviteeId),
			InviteeNickname: invitation.FkInvitee.Nickname,

			InviteRole: uint(invitation.TargetRole),
			Status:     invitation.Status,
			Project: dtos.InnerProject{
				Id:           uint(invitation.FkProject.Id),
				Title:        invitation.FkProject.Title,
				WorksetId:    uint(invitation.FkProject.WorksetId),
				WorksetIndex: invitation.FkProject.WorksetIndex,
				Status:       overallStatus,
			},
		})
	}

	return dtosInvitations, nil
}

// GetInvitationRecievedByUserId 实现 LaborService 接口的 GetInvitationRecievedByUserId 方法
func (s *laborServiceImpl) GetInvitationRecievedByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.InvitationBasic, error) {
	invitations, err := s.invitationRepo.GetInvitationRecievedByUserId(dbmodels.PrimaryKey(userId))
	if err != nil {
		s.logger.Error("GetInvitationRecievedByUserId 调用 GetInvitationRecievedByUserId 中出现错误",
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err))
		return nil, err
	}

	var dtosInvitations []*dtos.InvitationBasic

	for _, invitation := range invitations {
		// 封装 project 为 overall status
		overallStatus := dtos.ProjectOverallStatus(0)

		overallStatus.SetTranslatingStatus(uint(invitation.FkProject.TranslateStatus))
		overallStatus.SetReviewingStatus(uint(invitation.FkProject.ReviewStatus))
		overallStatus.SetLetteringStatus(uint(invitation.FkProject.LetterStatus))
		overallStatus.SetProofreadingStatus(uint(invitation.FkProject.ProofStatus))
		overallStatus.SetReviewingStatus(uint(invitation.FkProject.ReviewStatus))
		switch invitation.FkProject.IsPublished {
		case true:
			overallStatus.SetPublishedStatus(dtos.PROJECT_STATUS_COMPLETED)
		case false:
			overallStatus.SetPublishedStatus(dtos.PROJECT_STATUS_UNSET)
		}

		dtosInvitations = append(dtosInvitations, &dtos.InvitationBasic{
			Id:              uint(invitation.Id),
			InviterId:       uint(invitation.InviterId),
			InviteeId:       uint(invitation.InviteeId),
			InviteeNickname: invitation.FkInvitee.Nickname,

			InviteRole: uint(invitation.TargetRole),
			Status:     invitation.Status,
			Project: dtos.InnerProject{
				Id:           uint(invitation.FkProject.Id),
				Title:        invitation.FkProject.Title,
				WorksetId:    uint(invitation.FkProject.WorksetId),
				WorksetIndex: invitation.FkProject.WorksetIndex,
				Status:       overallStatus,
			},
		})
	}

	return dtosInvitations, nil
}

// CreateInvitation 实现 LaborService 接口的 CreateInvitation 方法
func (s *laborServiceImpl) CreateInvitation(inviterId, inviteeId, projectId uint, targetRole uint) error {
	if targetRole == 0 {
		return errors.New("目标职位不可以为空")
	}

	if err := s.invitationRepo.CreateInvitation(
		dbmodels.PrimaryKey(inviterId),
		dbmodels.PrimaryKey(inviteeId),
		dbmodels.PrimaryKey(projectId),
		targetRole,
	); err != nil {
		s.logger.Error("CreateInvitation 调用 CreateInvitation 中出现错误",
			slog.Uint64("inviter_id", uint64(inviterId)),
			slog.Uint64("invitee_id", uint64(inviteeId)),
			slog.Uint64("project_id", uint64(projectId)),
			slog.Any("error", err))
		return err
	}

	return nil
}

// GetApplisRecievedByUserId 实现 LaborService 接口的 GetApplisRecievedByUserId 方法
func (s *laborServiceImpl) GetApplisRecievedByUserId(projectId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error) {
	applications, err := s.appliRepo.GetApplicationRecievedByUserId(dbmodels.PrimaryKey(projectId), (pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetApplisRecievedByUserId 调用 GetApplicationRecievedByUserId 中出现错误",
			slog.Uint64("project_id", uint64(projectId)),
			slog.Any("error", err))
		return nil, err
	}

	var dtosApplications []*dtos.ApplicationBasic

	for _, application := range applications {
		dtosApplications = append(dtosApplications, &dtos.ApplicationBasic{
			Id:          uint(application.Id),
			ApplicantId: uint(application.ApplicantId),
			Nickname:    application.FkApplicant.Nickname,
			Role:        uint(application.TargetRole),
			Status:      application.Status,
			Project: dtos.InnerProject{
				Id:           uint(application.FkProject.Id),
				Title:        application.FkProject.Title,
				WorksetId:    uint(application.FkProject.WorksetId),
				WorksetIndex: application.FkProject.WorksetIndex,
				Status:       dtos.ProjectOverallStatus(0), // 这里可以根据需要设置状态
			},
		})
	}

	return dtosApplications, nil
}

// GetApplisSentByUserId 实现 LaborService 接口的 GetApplisSentByUserId 方法
func (s *laborServiceImpl) GetApplisSentByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error) {
	applications, err := s.appliRepo.GetApplicationSentByUserId(dbmodels.PrimaryKey(userId), (pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetApplisSentByUserId 调用 GetApplicationSentByUserId 中出现错误",
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err))
		return nil, err
	}

	var dtosApplications []*dtos.ApplicationBasic

	for _, application := range applications {
		dtosApplications = append(dtosApplications, &dtos.ApplicationBasic{
			Id:          uint(application.Id),
			ApplicantId: uint(application.ApplicantId),
			Nickname:    application.FkApplicant.Nickname,
			Role:        uint(application.TargetRole),
			Status:      application.Status,
			Project: dtos.InnerProject{
				Id:           uint(application.FkProject.Id),
				Title:        application.FkProject.Title,
				WorksetId:    uint(application.FkProject.WorksetId),
				WorksetIndex: application.FkProject.WorksetIndex,
				Status:       dtos.ProjectOverallStatus(0), // 这里可以根据需要设置状态
			},
		})
	}

	return dtosApplications, nil
}

// CreateApplication 实现 LaborService 接口的 CreateApplication 方法
func (s *laborServiceImpl) CreateApplication(applicantId, projectId uint, targetRole uint) error {
	if targetRole == 0 {
		return errors.New("目标职位不可以为空")
	}

	if err := s.appliRepo.CreateApplication(
		dbmodels.PrimaryKey(applicantId),
		dbmodels.PrimaryKey(projectId),
		targetRole,
	); err != nil {
		s.logger.Error("CreateApplication 调用 CreateApplication 中出现错误",
			slog.Uint64("applicant_id", uint64(applicantId)),
			slog.Uint64("project_id", uint64(projectId)),
			slog.Any("error", err))
		return err
	}

	return nil
}

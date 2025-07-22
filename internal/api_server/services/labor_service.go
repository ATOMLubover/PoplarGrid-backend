package services

import (
	"context"
	"errors"
	"log/slog"
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/repos"
	"poplargrid/internal/shared/dbmodels"
	"poplargrid/internal/shared/txutils"

	"gorm.io/gorm"
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
	// ProcessInvitation 处理申请的状态更新
	ProcessInvitation(userId, invitationId uint, status int) error

	// GetApplisRecievedByUserId 获取指定用户收到的申请列表，根据项目数量分页
	GetApplisRecievedByUserId(projectId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error)
	// GetApplisSentByUserId 获取指定用户发出的申请列表
	GetApplisSentByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.ApplicationBasic, error)

	// CreateApplication 创建一个新的申请
	CreateApplication(applicantId, projectId uint, targetRole uint) error
	// ProcessApplication 处理申请的状态更新
	ProcessApplication(userId, applicationId uint, status int) error
}

// laborServiceImpl 是 LaborService 接口的实现
type laborServiceImpl struct {
	invitationRepo repos.InvitationRepo
	appliRepo      repos.AppliRepo
	laborRepo      repos.LaborRepo
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
			Status:     int(invitation.Status),
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
			Status:     int(invitation.Status),
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

	// 如果被邀请者已经在对应项目中，则不允许重复邀请
	labors, err := s.laborRepo.SelectByUserId(
		dbmodels.PrimaryKey(inviteeId), []dbmodels.PrimaryKey{dbmodels.PrimaryKey(projectId)},
	)
	if err != nil || len(labors) != 0 {
		return errors.New("受邀请成员无法再次受邀于同一项目")
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

// ProcessInvitation 实现 LaborService 接口的 ProcessInvitation 方法
func (s *laborServiceImpl) ProcessInvitation(userId, invitationId uint, status int) error {
	if status != int(dbmodels.STATUS_PENDING) && status != int(dbmodels.STATUS_REJECTED) {
		return errors.New("无效的状态")
	}

	// 在事务中执行该操作
	coordinater := txutils.NewTransactionCoordinator(s.invitationRepo.GetHandle())

	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 创建基于事务的 repo
		invitationRepo := repos.NewInvitationRepo(tx)
		laborRepo := repos.NewLaborRepo(tx)

		// 首先检查是否是当前用户可以处理的邀请
		invitation, err := invitationRepo.SelectById(dbmodels.PrimaryKey(invitationId))
		if err != nil {
			s.logger.Error("ProcessInvitation 调用 SelectById 中出现错误",
				slog.Uint64("invitation_id", uint64(invitationId)),
				slog.Any("error", err))
			return err, nil
		}

		if invitation.InviteeId != dbmodels.PrimaryKey(userId) {
			return errors.New("无权处理该邀请"), nil
		}

		// 首先更新邀请的状态
		if err := invitationRepo.UpdateInvitationStatus(dbmodels.PrimaryKey(invitationId), dbmodels.Status(status)); err != nil {
			s.logger.Error("ProcessInvitation 调用 UpdateInvitationStatus 中出现错误",
				slog.Uint64("invitation_id", uint64(invitationId)),
				slog.Any("error", err))
			return err, nil
		}

		// 如果状态是接受，则创建 labor 记录
		if status == int(dbmodels.STATUS_ACCEPTED) {
			laborDivision := &dbmodels.ProjectLaborDivision{
				UserId:    invitation.InviteeId,
				ProjectId: invitation.ProjectId,
				LaborRole: dbmodels.LaborMask(invitation.TargetRole),
			}

			if err := laborRepo.CreateLaborDivision(laborDivision); err != nil {
				s.logger.Error("ProcessInvitation 调用 CreateLaborDivision 中出现错误",
					slog.Uint64("invitation_id", uint64(invitationId)),
					slog.Any("error", err))
				return err, nil
			}
		}

		return nil, nil

	}); err != nil {
		s.logger.Error("ProcessInvitation 调用 RunInTransaction 中出现错误",
			slog.Uint64("invitation_id", uint64(invitationId)),
			slog.Any("error", err))
		return errors.New("处理邀请失败")
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
			Status:      int(application.Status),
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
			Status:      int(application.Status),
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

	// 如果已经加入该项目，则不允许重复申请
	labors, err := s.laborRepo.SelectByUserId(
		dbmodels.PrimaryKey(applicantId), []dbmodels.PrimaryKey{dbmodels.PrimaryKey(projectId)},
	)
	if err != nil || len(labors) != 0 {
		return errors.New("申请成员无法再次申请同一项目")
	}

	// 在事务中执行下述逻辑，因为要提取 project 的 principal_id
	coordinater := txutils.NewTransactionCoordinator(s.appliRepo.GetHandle())

	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 创建基于事务的 repo
		projectRepo := repos.NewProjectRepo(tx)
		appliRepo := repos.NewAppliRepo(tx)

		// 先从 project 中获取 principal_id
		// TODO：这里是非并发安全的（在 project 被删除时）
		proj, err := projectRepo.SelectById(dbmodels.PrimaryKey(projectId))
		if err != nil {
			s.logger.Error("CreateApplication 调用 SelectById 中出现错误",
				slog.Uint64("project_id", uint64(projectId)),
				slog.Any("error", err))
			return err, nil
		}

		if err := appliRepo.CreateApplication(
			dbmodels.PrimaryKey(applicantId),
			dbmodels.PrimaryKey(proj.PrincipalId), // 使用 project 的 principal_id
			dbmodels.PrimaryKey(projectId),
			targetRole,
		); err != nil {
			return err, nil
		}

		return nil, nil

	}); err != nil {
		s.logger.Error("CreateApplication 调用 RunInTransaction 中出现错误",
			slog.Uint64("applicant_id", uint64(applicantId)),
			slog.Uint64("project_id", uint64(projectId)),
			slog.Any("error", err))
		return errors.New("创建申请失败")
	}

	return nil
}

// ProcessApplication 实现 laborServiceImpl 的 ProcessApplication 方法
func (s *laborServiceImpl) ProcessApplication(userId, applicationId uint, status int) error {
	if status != int(dbmodels.STATUS_ACCEPTED) && status != int(dbmodels.STATUS_REJECTED) {
		return errors.New("无效的状态")
	}

	// 在事务中执行该操作
	coordinater := txutils.NewTransactionCoordinator(s.appliRepo.GetHandle())

	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 创建基于事务的 repo
		appliRepo := repos.NewAppliRepo(tx)
		laborRepo := repos.NewLaborRepo(tx)

		// 首先检查是否是当前用户可以处理的申请
		application, err := appliRepo.SelectById(dbmodels.PrimaryKey(applicationId))
		if err != nil {
			s.logger.Error("ProcessApplication 调用 SelectById 中出现错误",
				slog.Uint64("application_id", uint64(applicationId)),
				slog.Any("error", err))
			return err, nil
		}

		if application.FkProject.PrincipalId != dbmodels.PrimaryKey(userId) {
			return errors.New("无权处理该申请"), nil
		}

		// 首先更新申请的状态
		if err := appliRepo.UpdateApplicationStatus(dbmodels.PrimaryKey(applicationId), dbmodels.Status(status)); err != nil {
			s.logger.Error("ProcessApplication 调用 UpdateApplicationStatus 中出现错误",
				slog.Uint64("application_id", uint64(applicationId)),
				slog.Any("error", err))
			return err, nil
		}

		// 如果状态是接受，则创建 labor 记录
		if status == int(dbmodels.STATUS_ACCEPTED) {
			laborDivision := &dbmodels.ProjectLaborDivision{
				UserId:    application.ApplicantId,
				ProjectId: application.ProjectId,
				LaborRole: dbmodels.LaborMask(application.TargetRole),
			}

			if err := laborRepo.CreateLaborDivision(laborDivision); err != nil {
				s.logger.Error("ProcessApplication 调用 CreateLaborDivision 中出现错误",
					slog.Uint64("application_id", uint64(applicationId)),
					slog.Any("error", err))
				return err, nil
			}
		}

		return nil, nil

	}); err != nil {
		s.logger.Error("ProcessApplication 调用 RunInTransaction 中出现错误",
			slog.Uint64("application_id", uint64(applicationId)),
			slog.Any("error", err))
		return errors.New("处理申请失败")
	}

	return nil
}

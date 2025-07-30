package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// LaborMask 定义了分工的掩码类型
type LaborMask uint32

const (
	LABOR_PRINCIPAL_MASK   LaborMask = 1 << iota // 创建者 + 负责人
	LABOR_SRC_PROV_MASK                          // 图源
	LABOR_PERFECTOR_MASK                         // 美工
	LABOR_TRANSLATOR_MASK                        // 翻译
	LABOR_PROOFREADER_MASK                       // 校对
	LABOR_LETTERER_MASK                          // 嵌字
	LABOR_REVIEWER_MASK                          // 嵌字审核
	LABOR_PUBLISHER_MASK                         // 发布者
)

// HasRole 检查 LaborMask 是否包含某个职责
func (m LaborMask) HasRole(roleMask LaborMask) bool {
	return (m & roleMask) != 0
}

// AddRole 为 LaborMask 添加一个职责
func (m *LaborMask) AddRole(roleMask LaborMask) {
	*m |= roleMask
}

// RemoveRole 从 LaborMask 中移除一个职责
func (m *LaborMask) RemoveRole(roleMask LaborMask) {
	*m &= ^roleMask
}

// MemberListParams 定义了成员列表的查询参数
type MemberListParams struct {
	Offset int  // 偏移量
	Limit  int  // 限制数量
	TeamId uint // 汉化组 ID
}

// MemberInfo 定义了成员的基本信息
type MemberInfo struct {
	Id   uint      // 成员 ID
	User *UserInfo // 对应的用户信息
	Team *TeamInfo // 对应的团队信息
	Role LaborMask // 在组内的职责（掩码格式）
}

// MemberService 接口定义了成员服务的基本操作
type MemberService interface {
	// GetMembers 获取指定条件下的成员列表
	GetMembers(params *MemberListParams) ([]*MemberInfo, error)
}

// memberServiceImpl 是 MemberService 接口的实现
type memberServiceImpl struct {
	handle *gorm.DB
	logger *slog.Logger
}

// NewMemberService 创建一个新的 MemberService 实例
func NewMemberService(hdl *gorm.DB, lgr *slog.Logger) MemberService {
	return &memberServiceImpl{
		handle: hdl,
		logger: lgr,
	}
}

// GetTeams 实现 MemberService 接口的方法
func (s *memberServiceImpl) GetMembers(params *MemberListParams) ([]*MemberInfo, error) {
	// 查询条件为汉化组 ID
	teamPKey := models.PKey(params.TeamId)
	memberSpec := &models.MemberSpec{
		TeamId: &teamPKey,
	}
	// 查询的字段
	memberFields := &models.MemberFields{
		Id:     true,
		UserId: true,
		UserFields: &models.UserFields{
			Id:       true,
			Nickname: true,
		},
		TeamId:     true,
		TeamFields: true, // 需要 Preload 团队信息
		Roles:      true,
	}

	// 执行查询
	members, err := models.GetMember().SelectMany(
		s.handle, memberSpec, memberFields,
		&params.Offset, &params.Limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("GetMembers 没有找到指定汉化组的成员信息",
				slog.Any("team_id", params.TeamId),
				slog.Any("error", err))
			return nil, errors.New("没有找到指定汉化组的成员信息")
		}

		s.logger.Error("GetMembers 查询指定汉化组的成员信息失败",
			slog.Any("team_id", params.TeamId),
			slog.Any("error", err))
		return nil, errors.New("查询指定汉化组的成员信息失败")
	}

	// 将查询结果转换为 MemberInfo
	var memberInfos []*MemberInfo

	for _, m := range members {
		// 转换成 MemberInfo
		memberInfos = append(memberInfos, memberModelToInfo(m))
	}

	return memberInfos, nil
}

// buildLaborMask 将数据库的角色字段转换为 LaborMask
func buildLaborMask(member *models.Member) LaborMask {
	var mask LaborMask

	if member.IsAdmin {
		mask.AddRole(LABOR_PRINCIPAL_MASK)
	}
	if member.IsTranslator {
		mask.AddRole(LABOR_TRANSLATOR_MASK)
	}
	if member.IsProofreader {
		mask.AddRole(LABOR_PROOFREADER_MASK)
	}
	if member.IsLetterer {
		mask.AddRole(LABOR_LETTERER_MASK)
	}
	if member.IsReviewer {
		mask.AddRole(LABOR_REVIEWER_MASK)
	}
	if member.IsPublisher {
		mask.AddRole(LABOR_PUBLISHER_MASK)
	}
	if member.IsSourceProvider {
		mask.AddRole(LABOR_SRC_PROV_MASK)
	}
	if member.IsPerfector {
		mask.AddRole(LABOR_PERFECTOR_MASK)
	}

	return mask
}

// memberModelToInfo 将 Member 模型转换为 MemberInfo
func memberModelToInfo(member *models.Member) *MemberInfo {
	m := &MemberInfo{
		Id: uint(member.BaseModel.Id),
		User: &UserInfo{
			Id: uint(member.UserId),
		},
		Team: &TeamInfo{
			Id: uint(member.TeamId),
		},
		Role: buildLaborMask(member),
	}

	if member.FkUser != nil {
		m.User.Nickname = member.FkUser.Nickname
		m.Team.MoetranId = member.FkTeam.MoetranId
	}

	if member.FkTeam != nil {
		m.Team.Name = member.FkTeam.Name
		m.Team.Description = member.FkTeam.Description
		m.Team.MoetranId = member.FkTeam.MoetranId
	}

	return m
}

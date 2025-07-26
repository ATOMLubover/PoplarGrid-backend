package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// RoleMask 定义了成员角色的掩码
type RoleMask uint32

// RoleMask 的掩码常量定义
const (
	ROLE_PRINCIPAL_MASK   RoleMask = 1 << iota // 管理员
	ROLE_SRC_PROV_MASK                         // 图源
	ROLE_PERFECTOR_MASK                        // 美工
	ROLE_TRANSLATOR_MASK                       // 翻译
	ROLE_PROOFREADER_MASK                      // 校对
	ROLE_LETTERER_MASK                         // 嵌字
	ROLE_REVIEWER_MASK                         // 嵌字审核
	ROLE_PUBLISHER_MASK                        // 发布者
)

// NewRoleMask 根据多个职责掩码创建一个新的复合掩码
func NewRoleMask(roles ...RoleMask) RoleMask {
	var mask RoleMask
	for _, role := range roles {
		mask |= role
	}
	return mask
}

// AddRole 为 RoleMask 添加一个职责
func (m *RoleMask) AddRole(roleMask RoleMask) {
	*m |= roleMask
}

// MemberListParams 定义了成员列表的查询参数
type MemberListParams struct {
	Offset int  // 偏移量
	Limit  int  // 限制数量
	UserId uint // 用户 ID
}

// MemberInfo 定义了成员的基本信息
type MemberInfo struct {
	Id   uint     // 成员 ID
	User UserInfo // 对应的用户信息
	Team TeamInfo // 对应的团队信息
	Role RoleMask // 在组内的职责（掩码格式）
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
	// 查询条件为用户 ID
	userPKey := models.PKey(params.UserId)
	memberSpec := &models.MemberSpec{
		UserId: &userPKey,
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
			s.logger.Warn("没有找到指定用户的成员信息",
				slog.Any("user_id", params.UserId),
				slog.Any("error", err))
			return nil, errors.New("没有找到指定用户的成员信息")
		}

		s.logger.Error("获取指定用户",
			slog.Any("user_id", params.UserId),
			slog.Any("error", err))
		return nil, errors.New("查询用户的成员信息失败")
	}

	// 将查询结果转换为 MemberInfo
	var memberInfos []*MemberInfo

	for _, m := range members {
		// 先构建 RoleMask
		roleMask := buildRoleMask(m)

		// 转换成 MemberInfo
		memberInfos = append(memberInfos, &MemberInfo{
			Id: uint(m.BaseModel.Id),
			User: UserInfo{
				Id:       uint(m.UserId),
				Nickname: m.FkUser.Nickname,
			},
			Team: TeamInfo{
				Id:          uint(m.TeamId),
				Name:        m.FkTeam.Name,
				Description: m.FkTeam.Description,
				MoetranId:   m.FkTeam.MoetranId,
			},
			Role: roleMask,
		})
	}

	return memberInfos, nil
}

// buildRoleMask 将数据库的角色字段转换为 RoleMask
func buildRoleMask(member *models.Member) RoleMask {
	var mask RoleMask

	if member.IsAdmin {
		mask.AddRole(ROLE_PRINCIPAL_MASK)
	}
	if member.IsPerfector {
		mask.AddRole(ROLE_PERFECTOR_MASK)
	}
	if member.IsSourceProvider {
		mask.AddRole(ROLE_SRC_PROV_MASK)
	}
	if member.IsTranslator {
		mask.AddRole(ROLE_TRANSLATOR_MASK)
	}
	if member.IsProofreader {
		mask.AddRole(ROLE_PROOFREADER_MASK)
	}
	if member.IsLetterer {
		mask.AddRole(ROLE_LETTERER_MASK)
	}
	if member.IsReviewer {
		mask.AddRole(ROLE_REVIEWER_MASK)
	}
	if member.IsPublisher {
		mask.AddRole(ROLE_PUBLISHER_MASK)
	}

	return mask
}

package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// UserInfo 定义了用户的基本信息
type UserInfo struct {
	ID         uint         // 用户 ID
	Nickname   string       // 昵称
	Email      string       // 邮箱
	QQNumber   int          // QQ 号
	IsAdmin    bool         // 是否是管理员
	Remark     string       // 补充备注
	MoetranId  string       // 龙译 ID
	MoetranJwt string       // 龙译 JWT
	Members    []MemberInfo // 在各个汉化组中的成员信息
}

// UserService 接口定义了用户服务的基本操作
type UserService interface {
	// GetUserDetail 获取指定用户的详细信息
	GetUserDetail(userId uint) (*UserInfo, Err)
}

// userServiceImpl 是 UserService 接口的实现
type userServiceImpl struct {
	handle *gorm.DB
	logger *slog.Logger
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(hdl *gorm.DB, lgr *slog.Logger) UserService {
	return &userServiceImpl{
		handle: hdl,
		logger: lgr,
	}
}

// GetUserDetail 实现 UserService 接口的方法，获取指定用户的详细信息
func (s *userServiceImpl) GetUserDetail(userId uint) (*UserInfo, Err) {
	// 查询条件为用户 ID
	userPKey := models.PKey(userId)
	userSpec := &models.UserSpec{
		Id: &userPKey,
	}
	// 查询的字段
	userFields := &models.UserFields{
		Id:         true,
		Nickname:   true,
		Email:      true,
		QQNumber:   true,
		IsAdmin:    true,
		Remark:     true,
		MoetranId:  true,
		MoetranJwt: true,
	}

	// 执行查询
	user, err := models.GetUser().SelectFirst(s.handle, userSpec, userFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("GetUserDetail 未获取到指定用户详情",
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("error", err))
			return nil, ErrNoSatifiedResults
		}
		s.logger.Error("GetUserDetail 获取指定用户详情失败", "error", err)
		return nil, ErrDatabaseFailure
	}

	// 将查询结果转换为 UserInfo
	userInfo := &UserInfo{
		ID:         uint(user.Id),
		Nickname:   user.Nickname,
		Email:      user.Email,
		IsAdmin:    user.IsAdmin,
		MoetranId:  user.MoetranId,
		MoetranJwt: user.MoetranJwt,
	}
	if user.QQNumber != nil {
		userInfo.QQNumber = int(*user.QQNumber)
	}
	if user.Remark != nil {
		userInfo.Remark = *user.Remark
	}

	// 获取用户的成员信息
	memberSpec := &models.MemberSpec{
		UserId: &userPKey,
	}
	memberFields := &models.MemberFields{
		Id:         true,
		Roles:      true,
		TeamId:     true,
		TeamFields: true,
	}

	members, err := models.GetMember().SelectMany(
		s.handle, memberSpec, memberFields,
		nil, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("GetUserDetail 查询 members 没有结果",
				slog.Uint64("userID", uint64(userId)))
			return nil, ErrNoSatifiedResults
		}
		s.logger.Error("GetUserDetail 查询用户成员失败", slog.Any("error", err))
		return nil, ErrDatabaseFailure
	}

	// 将成员信息转换为 MemberInfo
	for _, member := range members {
		userInfo.Members = append(userInfo.Members, *memberModelToInfo(member))
	}

	return userInfo, nil
}

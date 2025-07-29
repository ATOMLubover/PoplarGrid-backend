package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// UserInfo 定义了用户的基本信息
type UserInfo struct {
	Id       uint   // 用户 ID
	Nickname string // 昵称
	Email    string // 邮箱
	QQNumber int    // QQ 号
	IsAdmin  bool   // 是否是管理员
	Remark   string // 补充备注
}

// UserService 接口定义了用户服务的基本操作
type UserService interface {
	// GetUserDetail 获取指定用户的详细信息
	GetUserDetail(userId uint) (*UserInfo, error)
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
func (s *userServiceImpl) GetUserDetail(userId uint) (*UserInfo, error) {
	// 查询条件为用户 ID
	userPKey := models.PKey(userId)
	userSpec := &models.UserSpec{
		Id: &userPKey,
	}
	// 查询的字段
	userFields := &models.UserFields{
		Id:       true,
		Nickname: true,
		Email:    true,
		QQNumber: true,
		IsAdmin:  true,
		Remark:   true,
	}

	// 执行查询
	user, err := models.GetUser().SelectFirst(s.handle, userSpec, userFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("GetUserDetail 未获取到指定用户详情",
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("error", err))
			return nil, errors.New("用户不存在")
		}
		s.logger.Error("GetUserDetail 获取指定用户详情失败", "error", err)
		return nil, errors.New("获取指定用户详情失败")
	}

	// 将查询结果转换为 UserInfo
	userInfo := &UserInfo{
		Id:       uint(user.Id),
		Nickname: user.Nickname,
		Email:    user.Email,
		QQNumber: *user.QQNumber,
		IsAdmin:  user.IsAdmin,
		Remark:   *user.Remark,
	}

	return userInfo, nil
}

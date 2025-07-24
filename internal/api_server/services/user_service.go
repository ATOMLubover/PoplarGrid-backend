package services

import (
	"errors"
	"fmt"
	"log/slog"
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/repos"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// UserService 接口定义了用户服务的基本操作
type UserService interface {
	// GetUserDetail 获取指定用户的详细信息
	GetUserDetail(userId uint) (*dtos.UserDetail, error)
}

// userServiceImpl 是 UserService 接口的实现
type userServiceImpl struct {
	userRepo repos.UserRepo
	logger   *slog.Logger
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(userRepo repos.UserRepo) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// GetUserDetail 实现 UserService 接口的方法，获取指定用户的详细信息
func (s *userServiceImpl) GetUserDetail(userId uint) (*dtos.UserDetail, error) {
	user, err := s.userRepo.SelectByUserId(dbmodels.PrimaryKey(userId))
	if err != nil {
		// 此处单独拦截 not found 错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("GetUserDetail 调用 SelectByUserId 中未找到用户",
				slog.Uint64("user_id", uint64(userId)))
			return nil, fmt.Errorf("未找到指定 ID 的用户")
		}

		s.logger.Error("GetUserDetail 调用 SelectByUserId 中出现错误",
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err))
		return nil, errors.New("无法读取到指定 user 的详细信息")
	}

	// 将 dbmodels.User 转换为 dtos.UserDetail
	userDetail := &dtos.UserDetail{
		UserBasic: dtos.UserBasic{
			Id:       uint(user.Id),
			Nickname: user.Nickname,
		},
		Email:         user.Email,
		PoplarIsAdmin: user.PoplarIsAdmin,
		Remark:        user.Remark,
		QqNumber:      fmt.Sprintf("%d", user.QqNumber),
		LastActive:    user.LastActive.Format(dtos.DTO_TIME_FORMAT),
	}

	return userDetail, nil
}

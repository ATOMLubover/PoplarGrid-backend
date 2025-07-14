package repos

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// user 所需基础字段
var kUserBasicFields = []string{
	"id",
	"nickname",
	"created_at",
	"last_active",
}

// UserRepo 接口定义了用户仓库的基本操作
type UserRepo interface {
	// SelectByUserId 获取指定用户 ID 的详细信息
	SelectByUserId(userId dbmodels.PrimaryKey) (*dbmodels.User, error)
}

// userRepoImpl 是 UserRepo 接口的实现
type userRepoImpl struct {
	handle *gorm.DB
}

// NewUserRepo 创建一个新的 UserRepo 实例
func NewUserRepo(handle *gorm.DB) UserRepo {
	return &userRepoImpl{
		handle: handle,
	}
}

// Table 限定当前操作的表名
func (r *userRepoImpl) Table() *gorm.DB {
	return r.handle.Table(dbmodels.User{}.TableName())
}

// SelectByUserId 实现 UserRepo 接口的方法，获取指定用户 ID 的详细信息
func (r *userRepoImpl) SelectByUserId(userId dbmodels.PrimaryKey) (*dbmodels.User, error) {
	var user dbmodels.User

	if err := r.Table().
		Select(kUserBasicFields).
		Where("id = ?", userId).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}

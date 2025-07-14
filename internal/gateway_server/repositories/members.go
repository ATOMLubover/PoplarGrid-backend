package repositories

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// members 表对应的 repo
type MembersRepo struct {
	DbCtx *gorm.DB
}

// NewMembersRepo 构造一个新的 MembersRepo 实例
func NewMembersRepo(db *gorm.DB) *MembersRepo {
	return &MembersRepo{
		DbCtx: db,
	}
}

// GetTable 获取 members 表的上下文引用
func (r *MembersRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodels.User{})
}

// SelectFullModelById 根据 memver.id 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullModelById(memberId uint) (*dbmodels.User, error) {
	var member dbmodels.User

	if err := r.GetTable().
		Where("id = ?", memberId).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectFullModelByNickname 根据 memeber.nickname 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullModelByNickname(nickname string) (*dbmodels.User, error) {
	var member dbmodels.User

	if err := r.GetTable().
		Where("nickname = ?", nickname).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectAllToSlice 获取所有成员的完整信息
func (r *MembersRepo) SelectAllToSlice() ([]dbmodels.User, error) {
	var members []dbmodels.User

	if err := r.GetTable().
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

package repositories

import (
	"poplargrid/internal/shared/dbmodel"

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
	return r.DbCtx.Model(&dbmodel.Member{})
}

// SelectFullModelById 根据 memver.id 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullModelById(memberId uint) (*dbmodel.Member, error) {
	var member dbmodel.Member

	if err := r.GetTable().
		Where("id = ?", memberId).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectFullModelByNickname 根据 memeber.nickname 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullModelByNickname(nickname string) (*dbmodel.Member, error) {
	var member dbmodel.Member

	if err := r.GetTable().
		Where("nickname = ?", nickname).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectAllToSlice 获取所有成员的完整信息
func (r *MembersRepo) SelectAllToSlice() ([]dbmodel.Member, error) {
	var members []dbmodel.Member

	if err := r.GetTable().
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

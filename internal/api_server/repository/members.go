package repository

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

// SelectFullById 根据 memver.id 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullById(memberId uint) (*dbmodel.Member, error) {
	var member dbmodel.Member

	if err := r.GetTable().
		Where("id = ?", memberId).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectFullByNickname 根据 memeber.nickname 获取 member 完整信息
// 如果不存在则返回 nil 和错误
func (r *MembersRepo) SelectFullByNickname(nickname string) (*dbmodel.Member, error) {
	var member dbmodel.Member

	if err := r.GetTable().
		Where("nickname = ?", nickname).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectAllToSlice 获取所有成员的完整信息
func (r *MembersRepo) SelectAllToSlice(offset, num int) ([]*dbmodel.Member, error) {
	var members []*dbmodel.Member

	if err := r.GetTable().
		Select("id", "nickname", "email", "moetran_id",
			"poplar_is_admin", "labors", "remark", "last_active"). // 此处不返回 password_hash
		Offset(offset).
		Limit(num).
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}
